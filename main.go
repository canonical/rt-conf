package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	BOOT_CFG_GRUBCFG = "/boot/cfg/grub.cfg"
	ETC_DEFAULT_GRUB = "/etc/default/grub"

	cmdlineParams    = "isolcpus=8-9" // Text to append
	regexGrubDefault = `^(GRUB_CMDLINE_LINUX_DEFAULT=")([^"]*)(")$`

	// Old regex: linux\s*\/*\w*\/vmlinuz-\d.\d.\d
	regexGrubcfg = `linux\s*\/*\w*\/vmlinuz-\d.\d.\d`
)

var defaultConfig string = ""

func init() {
	snapCommon := os.Getenv("SNAP_COMMON")
	if snapCommon == "" {
		defaultConfig = snapCommon + "/config.yaml"
	}
}

// For /boot/cfg/grub.cfg and /etc/default/grub
func (c *InternalConfig) InjectToFile() error {
	err := readConfigFile(c)
	if err != nil {
		return err
	}

	fmt.Println("KernelCmdline: ", c.data.KernelCmdline)

	cfgFile, err := os.Open(c.grubCfg.file)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		return err
	}
	defer cfgFile.Close()

	// Create a temporary file to write the modified content
	tmpFileCfg, err := os.CreateTemp("", "config-modified-")
	if err != nil {
		fmt.Printf("Failed to create temp file: %v\n", err)
		return err
	}
	defer os.Remove(tmpFileCfg.Name()) // Clean up after execution if necessary

	scannerCfg := bufio.NewScanner(cfgFile)
	for scannerCfg.Scan() {
		line := scannerCfg.Text()

		if c.grubCfg.pattern.MatchString(line) {
			for _, param := range c.data.KernelCmdline {
				line += " " + param
			}
		}
		// Write the line to the temporary file
		_, err := tmpFileCfg.WriteString(line + "\n")
		if err != nil {
			fmt.Printf("Failed to write to temp file: %v\n", err)
			return err
		}
	}

	if err := scannerCfg.Err(); err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		return err
	}

	// Replace the original file with the modified one
	tmpFileCfg.Close()
	err = os.Rename(tmpFileCfg.Name(), c.grubCfg.file)
	if err != nil {
		fmt.Printf("Failed to replace original file: %v\n", err)
	}
	fmt.Println("File /boot/grub/grub.cfg updated successfully.")

	// Second part ----------------------------------------------
	// Modifying the /etc/default/grub file

	defaultFile, err := os.Open(c.grubDefault.file)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		return err
	}
	defer cfgFile.Close()

	// Create a temporary file to write the modified content
	tmpFileDefault, err := os.CreateTemp("", "config-modified-")
	if err != nil {
		fmt.Printf("Failed to create temp file: %v\n", err)
		return err
	}
	defer os.Remove(tmpFileDefault.Name()) // Clean up after execution if necessary

	scannerDefault := bufio.NewScanner(defaultFile)
	for scannerDefault.Scan() {
		line := scannerDefault.Text()

		// BEGIN: THIS IS GOING TO BE DIFFERENT FOR /etc/default/grub

		if c.grubDefault.pattern.MatchString(line) {

			// Extract existing parameters
			matches := c.grubDefault.pattern.FindStringSubmatch(line)
			existing := matches[2]
			// Append new parameters
			updatedParams := strings.TrimSpace(existing + " " + strings.Join(c.data.KernelCmdline, " "))
			// Reconstruct the line with updated parameters
			line = fmt.Sprintf(`%s%s%s`, matches[1], updatedParams, matches[3])

			// Extract existing parameters
			// matches := c.grubDefault.pattern.FindStringSubmatch(line)
			// existing := matches[len(matches)-1]
			// // Append new parameters
			// updatedParams := existing + " " + strings.Join(c.data.KernelCmdline, " ")
			// // Reconstruct the line with updated parameters
			// line = fmt.Sprintf(`%s%s%s`, matches[0], updatedParams, matches[len(matches)-1])
		}

		// END: THIS IS GOING TO BE DIFFERENT FOR /etc/default/grub

		// Write the line to the temporary file
		_, err := tmpFileDefault.WriteString(line + "\n")
		if err != nil {
			fmt.Printf("Failed to write to temp file: %v\n", err)
			return err
		}
	}

	if err := scannerDefault.Err(); err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		return err
	}

	// Replace the original file with the modified one
	tmpFileDefault.Close()
	err = os.Rename(tmpFileDefault.Name(), c.grubDefault.file)
	if err != nil {
		fmt.Printf("Failed to replace original file: %v\n", err)
	}
	fmt.Println("File /etc/default/grub updated successfully.")

	return nil
}

func main() {

	configPath := flag.String("config", defaultConfig, "Path to the configuration file")

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

	err := iCfg.InjectToFile()
	if err != nil {
		fmt.Printf("Failed to inject to file: %v\n", err)
		os.Exit(1)
	}

	// // Also modifying the /etc/default/grub file
	// regexDefaultGrub := `GRUB_CMDLINE_LINUX_DEFAULT="(.*)"` // regex to match vmlinuz line
	// iCfg.pattern = regexp.MustCompile(regexDefaultGrub)
	// iCfg.grubFile = "/etc/default/grub"
	// fmt.Println("Config path: ", iCfg.configFile)
	// fmt.Println("Grub path: ", iCfg.grubFile)
	// injectToFile(&iCfg)

	// data, err := os.ReadFile(*configPath)
	// if err != nil {
	// 	fmt.Printf("Failed to read file: %v\n", err)
	// 	return
	// }

	// config := Config{}
	// {
	// 	err := yaml.Unmarshal([]byte(data), &config)
	// 	if err != nil {
	// 		fmt.Printf("Failed to unmarshal data: %v\n", err)
	// 		return
	// 	}
	// }

	// fmt.Println("KernelCmdline: ", config.KernelCmdline)

	// // Open the configuration file
	// file, err := os.Open(*grubCfgPath)
	// if err != nil {
	// 	fmt.Printf("Failed to open file: %v\n", err)
	// 	return
	// }
	// defer file.Close()

	// // Create a temporary file to write the modified content
	// tmpFile, err := os.CreateTemp("", "config-modified-")
	// if err != nil {
	// 	fmt.Printf("Failed to create temp file: %v\n", err)
	// 	return
	// }
	// defer os.Remove(tmpFile.Name()) // Clean up after execution if necessary

	// // Compile the regex pattern
	// re := regexp.MustCompile(regexPattern)

	// scanner := bufio.NewScanner(file)
	// for scanner.Scan() {
	// 	line := scanner.Text()

	// 	if re.MatchString(line) {
	// 		for _, param := range config.KernelCmdline {
	// 			line += " " + param
	// 		}
	// 	}
	// 	// Write the line to the temporary file
	// 	_, err := tmpFile.WriteString(line + "\n")
	// 	if err != nil {
	// 		fmt.Printf("Failed to write to temp file: %v\n", err)
	// 		return
	// 	}
	// }

	// if err := scanner.Err(); err != nil {
	// 	fmt.Printf("Failed to read file: %v\n", err)
	// 	return
	// }

	// // Replace the original file with the modified one
	// tmpFile.Close()
	// err = os.Rename(tmpFile.Name(), *grubCfgPath)
	// if err != nil {
	// 	fmt.Printf("Failed to replace original file: %v\n", err)
	// }
	// fmt.Println("File updated successfully.")
}
