package data

import "testing"

func TestIRQTunning_Validate(t *testing.T) {
	tests := []struct {
		name    string
		c       IRQTunning
		wantErr bool
	}{
		{
			name: "valid test",
			c: IRQTunning{
				CPUs: "0,1",
				Filter: IRQFilter{
					Actions:  `action`,
					ChipName: `chip_name`,
					Name:     `name`,
					Type:     `type`,
				},
			},
			wantErr: false,
		},
		{
			name: "valid regex",
			c: IRQTunning{
				CPUs: "0-n",
				Filter: IRQFilter{
					Actions:  `nvme`,
					ChipName: `-PCI-`,
					Name:     `name`,
					Type:     `type`,
				},
			},
			wantErr: false,
		},
		{
			name: "valid regex",
			c: IRQTunning{
				CPUs: "0,n",
				Filter: IRQFilter{
					Actions:  `nvme`,
					ChipName: `-PCI-`,
					Name:     `\d`,
					Type:     `type`,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid regex",
			c: IRQTunning{
				CPUs: "0,1",
				Filter: IRQFilter{
					Actions:  `(?!abc)def`, // Negative lookahead
					ChipName: `chip_name`,
					Name:     `name`,
					Type:     `type`,
				},
			},
			wantErr: true,
		},
	}

	t.Log("Running IRQTunning.Validate() tests")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Validate(); (err != nil) != tt.wantErr {
				t.Logf("Testing: %v\nFilters:\n\t%v", tt.name, tt.c.Filter)
				t.Errorf("IRQTunning.Validate() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}
