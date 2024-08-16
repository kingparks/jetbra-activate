package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"
)

var jetPath = ""
var agtPath = ""

func init() {
	switch runtime.GOOS {
	case "windows":
		jetPath = os.Getenv("APPDATA") + "/JetBrains"
	case "darwin":
		jetPath = "/Users/" + os.Getenv("USER") + "/Library/Application Support/JetBrains"
	case "linux":
		jetPath = "/home/" + os.Getenv("USER") + "/.config/JetBrains"
	}
	agtPath = jetPath + "/active-agt.jar"
	if runtime.GOOS == "windows" {
		agtPath = strings.ReplaceAll(agtPath, "/", "\\")
	}
}

func Active(software string) {
	software = strings.ToLower(software)
	currCrackPath := "script"
	var softVM []string
	switch runtime.GOOS {
	case "windows":
		softVM = append(softVM, software+".exe.vmoptions")
		softVM = append(softVM, software+"64.exe.vmoptions")
	case "darwin":
		softVM = append(softVM, software+".vmoptions")
	case "linux":
		softVM = append(softVM, software+".vmoptions")
		softVM = append(softVM, software+"64.vmoptions")
	}
	// 检查systemUser是否有中文字符，如果有中文，提示可能会因此而不生效，需要搜索如何把用户名目录转换为英文，然后再执行
	for _, runeValue := range jetPath {
		if unicode.Is(unicode.Han, runeValue) {
			switch lang {
			case "zh":
				fmt.Printf(red, "用户名目录中有中文字符，可能会因此而不生效，需要搜索如何把用户名目录转换为英文，转为英文后再执行!")
			default:
				fmt.Printf(red, "The username directory contains Chinese characters, which may cause it to fail to take effect. You need to search for how to convert the username directory to English and then execute it!")
			}
		}
	}

	if _, err := os.Stat(jetPath); os.IsNotExist(err) {
		_ = os.MkdirAll(jetPath, os.ModePerm)
	}

	jarFile := currCrackPath + "/active-agt.jar"
	plugins := currCrackPath + "/plugins"
	config := currCrackPath + "/config"

	jarFileData, err := scriptFS.ReadFile(jarFile)
	if err == nil {
		os.WriteFile(jetPath+"/active-agt.jar", jarFileData, 0644)
		copyDir(scriptFS, plugins, jetPath+"/plugins")
		copyDir(scriptFS, config, jetPath+"/config")
	} else {
		fmt.Printf(red, "active-agt.jar is missing, "+software+" crack failed!")
		os.Exit(1)
	}

	softwareInstall := false
	files, _ := os.ReadDir(jetPath)
	for _, file := range files {
		if file.IsDir() {
			if strings.Contains(strings.ToLower(file.Name()), software) {
				softwareInstall = true
				for _, vm := range softVM {
					vmPath := jetPath + "/" + file.Name() + "/" + vm
					err := os.WriteFile(vmPath, []byte("-javaagent:"+agtPath+"\n--add-opens=java.base/jdk.internal.org.objectweb.asm=ALL-UNNAMED\n--add-opens=java.base/jdk.internal.org.objectweb.asm.tree=ALL-UNNAMED\n"), 0644)
					if err != nil {
						fmt.Printf(red, err.Error())
					}
				}
			}
		}
	}

	// toolbox 1.20
	var toolBoxDir []string
	switch runtime.GOOS {
	case "windows":
		toolBoxDir = append(toolBoxDir, os.Getenv("USERPROFILE")+"/AppData/Local/JetBrains/Toolbox/apps")
		toolBoxDir = append(toolBoxDir, "C:/Program Files/Toolbox/apps")
		toolBoxDir = append(toolBoxDir, "D:/Program Files/Toolbox/apps")
		toolBoxDir = append(toolBoxDir, "E:/Program Files/Toolbox/apps")
		toolBoxDir = append(toolBoxDir, "F:/Program Files/Toolbox/apps")
		toolBoxDir = append(toolBoxDir, "C:/Program Files (x86)/Toolbox/apps")
		toolBoxDir = append(toolBoxDir, "D:/Program Files (x86)/Toolbox/apps")
		toolBoxDir = append(toolBoxDir, "E:/Program Files (x86)/Toolbox/apps")
		toolBoxDir = append(toolBoxDir, "F:/Program Files (x86)/Toolbox/apps")
	case "darwin":
		toolBoxDir = append(toolBoxDir, jetPath+"/Toolbox/apps")
	case "linux":
		toolBoxDir = append(toolBoxDir, os.Getenv("HOME")+"/.local/share/JetBrains/Toolbox/apps")
	}
	for _, dir := range toolBoxDir {
		// 如果目录不存在则跳过
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}
		_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if info == nil {
				return nil
			}
			if info.IsDir() {
				return nil
			}
			if strings.HasSuffix(info.Name(), ".vmoptions") {
				// /Users/xx/Library/Application Support/JetBrains/Toolbox/apps/Goland/ch-0/241.14494.238/GoLand.app.vmoptions
				// C:\Users\XX\AppData\Local\JetBrains\Toolbox\apps\Goland\ch-0\241.14494.238.vmoptions
				// /home/xxx/.local/share/JetBrains/Toolbox/apps/Goland/ch-0/241.14494.238.vmoptions
				if strings.Contains(strings.ToLower(path), software) {
					softwareInstall = true
					err := os.WriteFile(path, []byte("-javaagent:"+agtPath+"\n--add-opens=java.base/jdk.internal.org.objectweb.asm=ALL-UNNAMED\n--add-opens=java.base/jdk.internal.org.objectweb.asm.tree=ALL-UNNAMED\n"), 0644)
					if err != nil {
						fmt.Printf(red, err.Error())
					}
				}
			}
			return nil
		})
	}

	// 2019年及之前的版本
	var dir2019 string
	switch runtime.GOOS {
	case "windows":
		dir2019 = os.Getenv("USERPROFILE")
	case "darwin":
		dir2019 = os.Getenv("HOME") + "/Library/Preferences"
	case "linux":
		dir2019 = os.Getenv("HOME")
	}
	files, _ = os.ReadDir(dir2019)
	for _, file := range files {
		if file.IsDir() {
			if runtime.GOOS == "windows" || runtime.GOOS == "linux" {
				if !strings.HasPrefix(file.Name(), ".") {
					continue
				}
			}
			if strings.Contains(strings.ToLower(file.Name()), software) {
				softwareInstall = true
				for _, vm := range softVM {
					vmPath := dir2019 + "/" + file.Name()
					if runtime.GOOS == "windows" || runtime.GOOS == "linux" {
						vmPath += "/config"
					}
					vmPath += "/" + vm
					_ = os.MkdirAll(filepath.Dir(vmPath), 0755)
					err := os.WriteFile(vmPath, []byte("-javaagent:"+agtPath), 0644)
					if err != nil {
						fmt.Printf(red, err.Error())
					}
				}
			}
		}
	}

	if !softwareInstall {
		switch lang {
		case "zh":
			fmt.Printf(red, "\n请先运行过 "+software+" !")
		default:
			fmt.Printf(red, "\nPlease run "+software+" first!")
		}
		os.Exit(1)
	}
	return
}

func copyDir(srcFS embed.FS, src string, dst string) error {
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		os.MkdirAll(dst, os.ModePerm)
	}
	entries, err := srcFS.ReadDir(src)
	if err != nil {
		fmt.Println(err)
		return err
	}
	for _, entry := range entries {
		srcPath := src + "/" + entry.Name()
		dstPath := dst + "/" + entry.Name()
		if entry.IsDir() {
			err := os.MkdirAll(dstPath, os.ModePerm)
			if err != nil {
				fmt.Println(err)
				return err
			}
			err = copyDir(srcFS, srcPath, dstPath)
			if err != nil {
				fmt.Println(err)
				return err
			}
		} else {
			data, err := srcFS.ReadFile(srcPath)
			if err != nil {
				fmt.Println(err)
				return err
			}
			err = os.WriteFile(dstPath, data, 0644)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	}

	return nil
}
