package model_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/canonical/rt-conf/src/kcmd"
	"github.com/canonical/rt-conf/src/model"
)

const (
	grubSample = `
GRUB_DEFAULT=0
GRUB_TIMEOUT_STYLE=hidden
GRUB_TIMEOUT=0
GRUB_CMDLINE_LINUX_DEFAULT="quiet splash"
GRUB_CMDLINE_LINUX=""
`
)

type TestCase struct {
	Name        string
	Yaml        string
	Validations []struct {
		param string
		value string
	}
}

func assertError(t *testing.T, err error, expectErr bool) {
	t.Helper()
	if expectErr && err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !expectErr && err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func mainLogic(t *testing.T, c TestCase, i int) (string, error) {
	dir := t.TempDir()
	// Set up temporary config-file
	tempConfigPath := filepath.Join(dir, fmt.Sprintf("config-%d.yaml", i))
	if err := os.WriteFile(tempConfigPath, []byte(c.Yaml), 0o644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	// Setup temporary custom grub config file
	tempCustomCfgPath := filepath.Join(dir, fmt.Sprintf("rt-conf-%d.cfg", i))
	if err := os.WriteFile(tempCustomCfgPath, []byte(grubSample), 0o644); err != nil {
		t.Fatalf("failed to write grub: %v", err)
	}

	t.Cleanup(func() {
		if err := os.Remove(tempConfigPath); err != nil {
			t.Fatal(err)
		}
		if err := os.Remove(tempCustomCfgPath); err != nil {
			t.Fatal(err)
		}
	})

	t.Logf("tempConfigPath: %s\n", tempConfigPath)
	t.Logf("tempCustomCfgPath: %s\n", tempCustomCfgPath)

	var conf model.InternalConfig
	if err := conf.Data.LoadFromFile(tempConfigPath); err != nil {
		return "", fmt.Errorf("failed to load config file: %v", err)
	}

	conf.GrubCfg = model.Grub{
		GrubDropInFile: tempCustomCfgPath,
	}

	t.Logf("Config: %+v\n", conf)

	// Run the InjectToFile method
	_, err := kcmd.ProcessKcmdArgs(&conf)
	if err != nil {
		return "", fmt.Errorf("ProcessKcmdArgs failed: %v", err)
	}

	// Verify default-grub updates
	updatedGrub, err := os.ReadFile(tempCustomCfgPath)
	if err != nil {
		return "", fmt.Errorf("failed to read modified grub file: %v", err)
	}

	t.Log("\nGrub file: \n", string(updatedGrub))

	return string(updatedGrub), nil
}

func TestHappyYamlKcmd(t *testing.T) {
	happyCases := []TestCase{
		{
			Name: "Using all cpus 0-N",
			Yaml: `
kernel-cmdline:
  parameters:
    - isolcpus=0-N
    - nohz=on
    - nohz_full=0-N
    - kthread_cpus=0-N
    - irqaffinity=0-N
`,
			Validations: []struct {
				param string
				value string
			}{
				{"isolcpus", "0-N"},
				{"nohz", "on"},
				{"nohz_full", "0-N"},
				{"kthread_cpus", "0-N"},
				{"irqaffinity", "0-N"},
			},
		},
		{
			Name: "Isolating cpu 0",
			Yaml: `
kernel-cmdline:
  parameters:
    - isolcpus=0
    - nohz=off
    - nohz_full=1-N
    - kthread_cpus=1-N
    - irqaffinity=1-N
`,
			Validations: []struct {
				param string
				value string
			}{
				{"isolcpus", "0"},
				{"nohz", "off"},
				{"nohz_full", "1-N"},
				{"kthread_cpus", "1-N"},
				{"irqaffinity", "1-N"},
			},
		},
	}

	for i, c := range happyCases {
		s, err := mainLogic(t, c, i)
		if err != nil {
			t.Fatal(err)
		}
		t.Run("HappyCases", func(t *testing.T) {
			for j, tc := range c.Validations {
				t.Log("Test case: ", j)
				if !strings.Contains(s,
					fmt.Sprintf("%s=%s", tc.param, tc.value)) {
					t.Errorf("\nExpected %s=%s in grub file, but not found",
						tc.param, tc.value)
				}
			}
		})
	}
}

func TestUnhappyYamlKcmd(t *testing.T) {
	UnhappyCases := []TestCase{
		{
			Name: "Invalid parameter name",
			Yaml: `
kernel-cmdline:
  parameters:
		- 34foo=a
`,
			Validations: nil,
		},
		{
			Name: "isolcpus a is invalid",
			Yaml: `
kernel-cmdline:
  parameters:
		- isolcpus=a
		- nohz=on
		- nohz_full=0-N
		- kthread_cpus=0
		- irqaffinity=0-N
`,
			Validations: nil,
		},
		{
			Name: "kthread_cpus z is invalid",
			Yaml: `
kernel-cmdline:
  parameters:
		- isolcpus=0
		- nohz=on
		- nohz_full=0-N
		- kthread_cpus=z
		- irqaffinity=0-N
`,
			Validations: nil,
		},
		{
			Name: "nohz true is invalid",
			Yaml: `
kernel-cmdline:
  parameters:
		- isolcpus=0
		- nohz=true
		- nohz_full=0-N
		- kthread_cpus=0-N
		- irqaffinity=0-N
`,
			Validations: nil,
		},
		{
			Name: "isolcpus 100000000 is invalid",
			Yaml: `
kernel-cmdline:
  parameters:
		- isolcpus=100000000
		- nohz=off
		- nohz_full=0-N
		- kthread_cpus=0-N
		- irqaffinity=0-N
`,
			Validations: nil,
		},
	}

	for i, c := range UnhappyCases {
		t.Run(c.Name, func(t *testing.T) {
			_, err := mainLogic(t, c, i)
			if err == nil {
				t.Fatalf("Expected error, but got nil on YAML: \n%v",
					c.Yaml)
			}
		})
	}
}

func TestKcmdDuplicates(t *testing.T) {
	testCases := []struct {
		Name      string
		Cfg       model.KernelCmdline
		ExpectErr bool
	}{
		{
			Name: "No duplicates",
			Cfg: model.KernelCmdline{
				Parameters: []string{
					"nohz=on",
					"isolcpus=0",
				},
			},
			ExpectErr: false,
		},
		{
			Name: "Duplicated parameter with diferent value",
			Cfg: model.KernelCmdline{
				Parameters: []string{
					"nohz=on",
					"nohz=off",
				},
			},
			ExpectErr: true,
		},
		{
			Name: "Tag parameter duplicated",
			Cfg: model.KernelCmdline{
				Parameters: []string{
					"nohz",
					"nohz",
				},
			},
			ExpectErr: false,
		},
		{
			Name: "Empty value in list",
			Cfg: model.KernelCmdline{
				Parameters: []string{
					"nohz=on",
					"",
					"quiet",
				},
			},
			ExpectErr: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := tc.Cfg.HasDuplicates()
			assertError(t, err, tc.ExpectErr)
		},
		)
	}
}

func TestKcmdValidation(t *testing.T) {
	// Building long kernel parameters list for testing
	// case that violates COMMAND_LINE_SIZE limit
	parameters := make([]string, 0, 65)
	// If each key=value parameter is 31 characters long
	// So being x the maximum number of parameters:
	// 31 * x + x <= 2048 (+x because we need to count the spaces between parameters)
	// The maximum x is 64 so with 65 parameters we exceed the limit
	for i := range 65 {
		parameters = append(parameters, fmt.Sprintf("some_kernel_cmd_option%03d=value", i))
	}

	testCases := []struct {
		Name      string
		Cfg       model.KernelCmdline
		ExpectErr bool
	}{
		{
			Name: "invalid Long command line",
			Cfg: model.KernelCmdline{
				Parameters: parameters,
			},
			ExpectErr: true,
		},
		{
			Name: "invalid parameter name",
			Cfg: model.KernelCmdline{
				Parameters: []string{
					"42foo=bar",
				},
			},
			ExpectErr: true,
		},
		{
			Name: "invalid cpulist",
			Cfg: model.KernelCmdline{
				Parameters: []string{
					"nohz_full=foo",
				},
			},
			ExpectErr: true,
		},
		{
			Name: "invalid isolcpus",
			Cfg: model.KernelCmdline{
				Parameters: []string{
					"isolcpus=foo,0-3",
				},
			},
			ExpectErr: true,
		},
		{
			Name: "invalid empty parameters",
			Cfg: model.KernelCmdline{
				Parameters: []string{
					"",
				},
			},
			ExpectErr: true,
		},
		{
			Name: "valid tag paramter",
			Cfg: model.KernelCmdline{
				Parameters: []string{
					"nohz",
				},
			},
			ExpectErr: false,
		},
		{
			Name: "valid not validated parameter",
			Cfg: model.KernelCmdline{
				Parameters: []string{
					"amd_iommu=pgtbl_v2",
				},
			},
			ExpectErr: false,
		},
	}
	for _, tc := range testCases {
		err := tc.Cfg.Validate()
		assertError(t, err, tc.ExpectErr)
	}
}
