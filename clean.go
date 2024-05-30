package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var jbProducts = []string{"idea", "clion", "phpstorm", "goland", "pycharm", "webstorm", "webide", "rider", "datagrip", "rubymine", "appcode", "dataspell", "gateway", "jetbrains_client", "jetbrainsclient", "studio", "devecostudio"}

func Clean() {
	switch runtime.GOOS {
	case "darwin":
		homeDir, _ := os.UserHomeDir()
		vmSHFile := ".jetbrains.vmoptions.sh"
		myVmOptionsShellFile := homeDir + "/" + vmSHFile

		_ = os.Remove(myVmOptionsShellFile)

		for _, prd := range jbProducts {
			envName := strings.ToUpper(prd) + "_VM_OPTIONS"
			_ = os.Unsetenv(envName)
			cmd := fmt.Sprintf("launchctl unsetenv %s", envName)
			_ = exec.Command("sh", "-c", cmd).Run()
		}

		plistPath := homeDir + "/Library/LaunchAgents/jetbrains.vmoptions.plist"
		_ = os.Remove(plistPath)

		removeLineFromFile(homeDir+"/.profile", vmSHFile)
		removeLineFromFile(homeDir+"/.bash_profile", vmSHFile)
		removeLineFromFile(homeDir+"/.zshrc", vmSHFile)
	case "windows":
		isClean := true
		for _, prd := range jbProducts {
			envKey := strings.ToUpper(prd) + "_VM_OPTIONS"
			if os.Getenv(envKey) != "" {
				isClean = false
			}
			_ = os.Unsetenv(envKey)
		}
		if !isClean {
			var cleanVBSPath = jetPath + "/win_clean.vbs"
			cleanFileData, err := scriptFS.ReadFile("script/win_clean.vbs")
			if err == nil {
				os.WriteFile(cleanVBSPath, cleanFileData, 0644)
			}
			err = exec.Command("cmd.exe", "/c", cleanVBSPath).Run()
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
			}
		}
	case "linux":
		homeDir, _ := os.UserHomeDir()
		vmSHFile := ".jetbrains.vmoptions.sh"
		myVmOptionsShellFile := homeDir + "/" + vmSHFile

		_ = os.Remove(myVmOptionsShellFile)

		for _, prd := range jbProducts {
			envName := strings.ToUpper(prd) + "_VM_OPTIONS"
			_ = os.Unsetenv(envName)
		}

		removeLineFromFile(homeDir+"/.profile", vmSHFile)
		removeLineFromFile(homeDir+"/.bashrc", vmSHFile)
		removeLineFromFile(homeDir+"/.zshrc", vmSHFile)

		kdeEnvDir := homeDir + "/.config/plasma-workspace/env"
		_ = os.Remove(kdeEnvDir + "/jetbrains.vmoptions.sh")
	}
}

func removeLineFromFile(filePath string, lineToRemove string) {
	input, err := os.ReadFile(filePath)
	if err != nil {
		//fmt.Println(err)
		return
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, lineToRemove) {
			lines = append(lines[:i], lines[i+1:]...)
		}
	}
	output := strings.Join(lines, "\n")
	err = os.WriteFile(filePath, []byte(output), 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}
