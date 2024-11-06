package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
)

const (
	GRUB_FILE     = "/boot/cfg/grub.cfg"
	cmdlineParams = "isolcpus=8-9" // Text to append
)

var defaultConfig string = ""

func init() {
	snapCommon := os.Getenv("SNAP_COMMON")
	if snapCommon == "" {
		defaultConfig = snapCommon + "/config.yaml"
	}
}

// For /boot/cfg/grub.cfg and /etc/default/grub
func (c *InternalConfig) InjectToFile(pattern *regexp.Regexp) error {
	err := readConfigFile(c)
	if err != nil {
		return err
	}

	fmt.Println("KernelCmdline: ", c.data.KernelCmdline)

	file, err := os.Open(c.grubFile)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		return err
	}
	defer file.Close()

	// Create a temporary file to write the modified content
	tmpFile, err := os.CreateTemp("", "config-modified-")
	if err != nil {
		fmt.Printf("Failed to create temp file: %v\n", err)
		return err
	}
	defer os.Remove(tmpFile.Name()) // Clean up after execution if necessary

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if c.patternCfgGrub.MatchString(line) {
			for _, param := range c.data.KernelCmdline {
				line += " " + param
			}
		}
		// Write the line to the temporary file
		_, err := tmpFile.WriteString(line + "\n")
		if err != nil {
			fmt.Printf("Failed to write to temp file: %v\n", err)
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		return err
	}

	// Replace the original file with the modified one
	tmpFile.Close()
	err = os.Rename(tmpFile.Name(), c.grubFile)
	if err != nil {
		fmt.Printf("Failed to replace original file: %v\n", err)
	}
	fmt.Println("File updated successfully.")

	return nil
}

func scanFile(b bufio.Scanner) error {

}

func main() {

	configPath := flag.String("config", defaultConfig, "Path to the configuration file")
	grubPath := flag.String("grub", GRUB_FILE, "Path to the grub file")
	flag.Parse()

	regexGrubcfg := `linux\s\/boot\/vmlinuz-\d.\d.\d` // regex to match vmlinuz line
	iCfg := InternalConfig{
		grubFile:   *grubPath,
		configFile: *configPath,
		pattern:    regexp.MustCompile(regexGrubcfg),
	}
	fmt.Println("Config path: ", iCfg.configFile)
	fmt.Println("Grub path: ", iCfg.grubFile)
	injectToFile(&iCfg)

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
	// file, err := os.Open(*grubPath)
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
	// err = os.Rename(tmpFile.Name(), *grubPath)
	// if err != nil {
	// 	fmt.Printf("Failed to replace original file: %v\n", err)
	// }
	// fmt.Println("File updated successfully.")
}
