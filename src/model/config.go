package model

import (
	"fmt"
	"os"
	"syscall"
)

var expectedPermission os.FileMode = 0o644

var IsOwnedByRoot = func(fi os.FileInfo) bool {
	uid := fi.Sys().(*syscall.Stat_t).Uid
	return uid == 0 // Check if the file is owned by root
}

func LoadConfigFile(confPath string) (*Config, error) {
	fileInfo, err := os.Stat(confPath)
	if err != nil {
		return nil, fmt.Errorf("failed to find file: %v", err)
	}

	if fileInfo.Mode() != expectedPermission {
		return nil, fmt.Errorf(
			"file %s has invalid permissions: %v, expected permissions %v",
			confPath, fileInfo.Mode(), expectedPermission)
	}

	if !IsOwnedByRoot(fileInfo) {
		return nil, fmt.Errorf("file %s is not owned by root", confPath)
	}

	content, err := ReadYAML(confPath)
	if err != nil {
		return nil, err
	}

	/*
		TODO: Needs to implement proper validation of the parameters
		and parameters format

		validations to be configured:
			- key=value
			- flag
	*/
	return content, nil
}
