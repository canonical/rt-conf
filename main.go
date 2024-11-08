package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	// Grub files paths
	BOOT_CFG_GRUBCFG = "/boot/cfg/grub.cfg"
	ETC_DEFAULT_GRUB = "/etc/default/grub"

	// Default configuration file path
	DEFAULT_CONFIG_PATH = "/var/snap/rt-conf/common/config.yaml"

	regexGrubDefault = `^(GRUB_CMDLINE_LINUX=")([^"]*)(")$`

	regexGrubcfg = `linux\s*\/*\w*\/vmlinuz-\d.\d.\d`
)

// FileTransformer interface with a TransformLine method.
// This method is used to transform a line of a file.
//
// NOTE: This interface can be implemented also for RPi on classic
type FileTransformer interface {
	TransformLine(string) string
	GetFilePath() string
	GetPattern() *regexp.Regexp
}

// grubCfgTransformer handles transformations for /boot/grub/grub.cfg
type grubCfgTransformer struct {
	filePath string
	pattern  *regexp.Regexp
	params   []string
}

func (g *grubCfgTransformer) TransformLine(line string) string {
	// Append each kernel command line parameter to the matched line
	for _, param := range g.params {
		line += " " + param
	}
	return line
}

func (g *grubCfgTransformer) GetFilePath() string {
	return g.filePath
}

func (g *grubCfgTransformer) GetPattern() *regexp.Regexp {
	return g.pattern
}

// grubDefaultTransformer handles transformations for /etc/default/grub
type grubDefaultTransformer struct {
	filePath string
	pattern  *regexp.Regexp
	params   []string
}

func (g *grubDefaultTransformer) TransformLine(line string) string {
	// TODO: Add functionality to avoid duplications of parameters

	// Extract existing parameters
	matches := g.pattern.FindStringSubmatch(line)
	// Append new parameters
	updatedParams := strings.TrimSpace(matches[2] + " " + strings.Join(g.params, " "))
	// Reconstruct the line with updated parameters
	return fmt.Sprintf(`%s%s%s`, matches[1], updatedParams, matches[3])
}

func (g *grubDefaultTransformer) GetFilePath() string {
	return g.filePath
}

func (g *grubDefaultTransformer) GetPattern() *regexp.Regexp {
	return g.pattern
}

// InjectToGrubFiles inject the kernel command line parameters to the grub files.
// /boot/grub/grub.cfg and /etc/default/grub
func (c *InternalConfig) InjectToGrubFiles() error {
	err := readConfigFile(c)
	if err != nil {
		return err
	}

	fmt.Println("KernelCmdline: ", c.data.KernelCmdline)

	// grubCfg := &grubCfgTransformer{
	// 	filePath: c.grubCfg.file,
	// 	pattern:  c.grubCfg.pattern,
	// 	params:   c.data.KernelCmdline,
	// }

	grubDefault := &grubDefaultTransformer{
		filePath: c.grubDefault.file,
		pattern:  c.grubDefault.pattern,
		params:   c.data.KernelCmdline,
	}

	// // Process each file with its specific transformer
	// if err := processFile(grubCfg); err != nil {
	// 	return fmt.Errorf("error updating %s: %v", grubCfg.filePath, err)
	// }
	// fmt.Println("File /boot/grub/grub.cfg updated successfully.")

	if err := processFile(grubDefault); err != nil {
		return fmt.Errorf("error updating %s: %v", grubDefault.filePath, err)
	}
	fmt.Println("File /etc/default/grub updated successfully.")

	return nil
}

func main() {

	configPath := flag.String("config", DEFAULT_CONFIG_PATH, "Path to the configuration file")

	// Define the paths to the grub files as flags
	// To be used for testing purposes
	grubCfgPath := flag.String("grub-cfg", BOOT_CFG_GRUBCFG, "Path to the processed grub file")
	grubDefaultPath := flag.String("grub-default", ETC_DEFAULT_GRUB, "Path to the default grub file")

	flag.Parse()

	iCfg := InternalConfig{
		configFile: *configPath,
		grubCfg: grub{
			file:    *grubCfgPath,
			pattern: regexp.MustCompile(regexGrubcfg),
		},
		grubDefault: grub{
			file:    *grubDefaultPath,
			pattern: regexp.MustCompile(regexGrubDefault),
		},
	}

	fmt.Println("Config path: ", iCfg.configFile)
	fmt.Println("Grub path: ", iCfg.grubCfg.file)

	err := iCfg.InjectToGrubFiles()
	if err != nil {
		fmt.Printf("Failed to inject to file: %v\n", err)
		os.Exit(1)
	}

	// Run update-grub command
	cmd := exec.Command("update-grub")
	var out []byte
	out, err = cmd.CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		fmt.Printf("Failed to update grub: %v\n", err)
		os.Exit(1)
	}

}
