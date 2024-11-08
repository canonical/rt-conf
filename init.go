package main

import (
	"os"
	"sync"
)

var lock = &sync.Mutex{}

type singletonDefaultConfig struct {
	path string
}

var instance *singletonDefaultConfig

func getDefaultConfig() string {
	lock.Lock()
	defer lock.Unlock()

	return instance.path
}

func init() {
	defaultConfig := ""
	snapCommon := os.Getenv("SNAP_COMMON")
	if snapCommon != "" {
		defaultConfig = snapCommon + "/config.yaml"
	}

	instance = &singletonDefaultConfig{path: defaultConfig}
}
