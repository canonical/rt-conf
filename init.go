package main

import "os"

var defaultConfig string = ""

func init() {
	snapCommon := os.Getenv("SNAP_COMMON")
	if snapCommon == "" {
		defaultConfig = snapCommon + "/config.yaml"
	}
}
