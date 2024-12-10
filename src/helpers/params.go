package helpers

import (
	"fmt"

	"github.com/canonical/rt-conf/src/data"
	"github.com/canonical/rt-conf/src/validator"
)

var Parameters = []data.Param{
	{
		YAMLName:    "isolcpus",
		CmdlineName: "isolcpus",
		TransformFn: func(value interface{}) string {
			return fmt.Sprintf("isolcpus=%s", value)
		},
	},
	{
		YAMLName:    "dyntick-idle",
		CmdlineName: "nohz",
		TransformFn: func(value interface{}) string {
			validator.ValidateType(validator.TypeEnum["bool"], "dyntick-idle",
				value)
			if v, ok := value.(bool); ok && v {
				return "nohz=on"
			}
			return "nohz=off"
		},
	},
	{
		YAMLName:    "adaptive-ticks",
		CmdlineName: "nohz_full",
		TransformFn: func(value interface{}) string {
			return fmt.Sprintf("nohz_full=%s", value)
		},
	},
}
