package main

import (
	"crypto/md5"
	"embed"
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/denisbrodbeck/machineid"
	"howett.net/plist"
	"net"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/unknwon/i18n"
)

var version = 226

var hosts = []string{"https://idea.jeter.eu.org", "http://129.154.205.7:7191", "http://jetbra.serv00.net:7191", "http://ba.serv00.net:7191"}
var host = hosts[0]
var githubPath = "https://ghp.ci/https://github.com/kingparks/jetbra-activate/releases/download/latest/"
var err error

var green = "\033[32m%s\033[0m\n"
var yellow = "\033[33m%s\033[0m\n"
var hGreen = "\033[1;32m%s\033[0m"
var dGreen = "\033[4;32m%s\033[0m\n"
var red = "\033[31m%s\033[0m\n"
var defaultColor = "%s"
var lang, _ = getLocale()
var deviceID = getMacMD5()
var client = Client{Hosts: hosts}

//go:embed all:script
var scriptFS embed.FS

//go:embed all:locales
var localeFS embed.FS

type Tr struct {
	i18n.Locale
}

var tr *Tr

func main() {
	language := flag.String("l", lang, "set language, eg: zh, en, nl, ru, hu, tr")
	flag.Parse()

	localeFileEn, _ := localeFS.ReadFile("locales/en.ini")
	_ = i18n.SetMessage("en", localeFileEn)
	localeFileNl, _ := localeFS.ReadFile("locales/nl.ini")
	_ = i18n.SetMessage("nl", localeFileNl)
	localeFileRu, _ := localeFS.ReadFile("locales/ru.ini")
	_ = i18n.SetMessage("ru", localeFileRu)
	localeFileHu, _ := localeFS.ReadFile("locales/hu.ini")
	_ = i18n.SetMessage("hu", localeFileHu)
	localeFileTr, _ := localeFS.ReadFile("locales/tr.ini")
	_ = i18n.SetMessage("tr", localeFileTr)
	localeFileEs, _ := localeFS.ReadFile("locales/es.ini")
	_ = i18n.SetMessage("es", localeFileEs)
	lang = *language
	switch lang {
	case "zh":
		tr = &Tr{Locale: i18n.Locale{Lang: "zh"}}
	case "nl":
		tr = &Tr{Locale: i18n.Locale{Lang: "nl"}}
	case "ru":
		tr = &Tr{Locale: i18n.Locale{Lang: "ru"}}
	case "hu":
		tr = &Tr{Locale: i18n.Locale{Lang: "hu"}}
	case "tr":
		tr = &Tr{Locale: i18n.Locale{Lang: "tr"}}
	case "es":
		tr = &Tr{Locale: i18n.Locale{Lang: "es"}}
	default:
		tr = &Tr{Locale: i18n.Locale{Lang: "en"}}
	}

	fmt.Printf(green, tr.Tr("IntelliJ 授权")+` v`+strings.Join(strings.Split(fmt.Sprint(version), ""), "."))
	client.SetProxy(lang)
	sCount, sPayCount, _, _, exp := client.GetMyInfo(deviceID)
	fmt.Printf(green, tr.Tr("设备码")+":"+deviceID)
	expTime, _ := time.ParseInLocation("2006-01-02 15:04:05", exp, time.Local)
	fmt.Printf(green, tr.Tr("付费到期时间")+":"+exp)
	fmt.Printf("\033[32m%s\033[0m\u001B[1;32m %s \u001B[0m\033[32m%s\033[0m\u001B[1;32m %s \u001B[0m\u001B[32m%s\u001B[0m\n",
		tr.Tr("推广命令：(已推广"), sCount, tr.Tr("人,推广已付费"), sPayCount, tr.Tr("人；每推广10人或推广付费2人可获得一年授权)"))
	fmt.Printf(hGreen, "bash <(curl -Lk "+githubPath+"install.sh) "+deviceID+"\n")

	printAD()
	checkUpdate(version)
	fmt.Println()
	fmt.Printf(defaultColor, tr.Tr("选择要授权的产品："))
	jbProduct := []string{"IntelliJ IDEA", "CLion", "PhpStorm", "Goland", "PyCharm", "WebStorm", "Rider", "DataGrip", "DataSpell"}
	jbProductChoice := []string{"idea", "clion", "phpstorm", "goland", "pycharm", "webstorm", "rider", "datagrip", "dataspell"}
	for i, v := range jbProduct {
		fmt.Printf(hGreen, fmt.Sprintf("%d. %s\t", i+1, v))
	}
	fmt.Println()
	fmt.Print(tr.Tr("请输入产品编号（直接回车默认为1）："))
	productIndex := 1
	_, _ = fmt.Scanln(&productIndex)
	if productIndex < 1 || productIndex > len(jbProduct) {
		fmt.Println(tr.Tr("输入有误"))
		return
	}
	fmt.Println(tr.Tr("选择的产品为：") + jbProduct[productIndex-1])
	fmt.Println()

	// 到期了
	periodIndex := 1
	_ = []time.Duration{367 * 24 * time.Hour, 24 * time.Hour}
	if expTime.Before(time.Now()) {
		fmt.Printf(defaultColor, tr.Tr("选择有效期："))
		jbPeriod := []string{"1" + tr.Tr("年(购买)"), "24" + tr.Tr("小时(免费)")}
		for i, v := range jbPeriod {
			fmt.Printf(hGreen, fmt.Sprintf("%d. %s\t", i+1, v))
		}
		fmt.Println()
		fmt.Printf("%s", tr.Tr("请输入有效期编号（直接回车默认为1）："))
		_, _ = fmt.Scanln(&periodIndex)
		if periodIndex < 1 || periodIndex > len(jbPeriod) {
			fmt.Println(tr.Tr("输入有误"))
			return
		}
		fmt.Println(tr.Tr("选择的有效期为：") + jbPeriod[periodIndex-1])
		fmt.Println()
	}

	lic := ""
	for i := 0; i < 50; i++ {
		if i == 20 {
			Clean()
			Active(jbProductChoice[productIndex-1])
		}
		time.Sleep(20 * time.Millisecond)
		h := strings.Repeat("=", i) + strings.Repeat(" ", 49-i)
		fmt.Printf("\r%.0f%%[%s]", float64(i)/49*100, h)
	}
	fmt.Println()
	fmt.Println()

	switch periodIndex {
	case 2:
		isOk, result := client.GetLic(jbProductChoice[productIndex-1], periodIndex-1)
		if !isOk {
			fmt.Printf(red, result)
			return
		}
		lic = result
	case 1:
		// 没到期
		if expTime.After(time.Now()) {
			isOk, result := client.GetLic(jbProductChoice[productIndex-1], periodIndex-1)
			if !isOk {
				fmt.Printf(red, result)
				return
			}
			lic = result
			fmt.Println()
			goto Process
		}
		// 到期了
		payUrl, orderID := client.GetPayUrl()
		isCopyText := ""
		errClip := clipboard.WriteAll(payUrl)
		if errClip == nil {
			isCopyText = tr.Tr("（已复制到剪贴板）")
		}
		fmt.Println(tr.Tr("付费已到期,捐赠以获取一年期授权") + isCopyText)
		fmt.Printf(dGreen, payUrl)
		fmt.Println(tr.Tr("捐赠完成后请回车"))
		//检测控制台回车
	checkPay:
		_, _ = fmt.Scanln()
		isPay := client.PayCheck(orderID, deviceID)
		if !isPay {
			fmt.Println(tr.Tr("未捐赠,请捐赠完成后回车"))
			goto checkPay
		}
		isOk, result := client.GetLic(jbProductChoice[productIndex-1], periodIndex-1)
		if !isOk {
			fmt.Printf(red, result)
			return
		}
		fmt.Println()
	}

Process:
	isCopyText := ""
	err = clipboard.WriteAll(lic)
	if err == nil {
		isCopyText = tr.Tr("（已复制到剪贴板）")
	}
	fmt.Printf(yellow, tr.Tr("首次执行请重启IDE，然后填入下面授权码；非首次执行直接填入下面授权码即可")+isCopyText)
	switch runtime.GOOS {
	case "windows":
		_ = exec.Command("taskkill", "/IM", jbProductChoice[productIndex-1]+".exe", "/F").Run()
		_ = exec.Command("taskkill", "/IM", jbProductChoice[productIndex-1]+"64.exe", "/F").Run()
	case "darwin":
		_ = exec.Command("killall", "-9", jbProductChoice[productIndex-1]).Run()
	case "linux":
		_ = exec.Command("killall", "-9", "java").Run()
	}
	fmt.Println()
	fmt.Printf(hGreen, lic)
	fmt.Println()
	for i := 0; i < 4; i++ {
		_, _ = fmt.Scanln()
	}
}

func getMacMD5() string {
	// 获取本机的MAC地址
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("err:", err)
		return ""
	}
	var macAddress []string
	var wifiAddress []string
	var bluetoothAddress []string
	var macErrorStr string
	for _, inter := range interfaces {
		// 排除虚拟网卡
		hardwareAddr := inter.HardwareAddr.String()
		if hardwareAddr == "" {
			//fmt.Println(fmt.Sprintf("log: have not hardwareAddr :%+v",inter))
			continue
		}
		macErrorStr += inter.Name + ":" + hardwareAddr + "\n"
		virtualMacPrefixes := []string{
			"00:05:69", "00:0C:29", "00:1C:14", "00:50:56", // VMware
			"00:15:5D",             // Hyper-V
			"08:00:27", "0A:00:27", // VirtualBox
		}
		isVirtual := false
		for _, prefix := range virtualMacPrefixes {
			if strings.HasPrefix(hardwareAddr, strings.ToLower(prefix)) {
				isVirtual = true
				break
			}
		}
		if isVirtual {
			//fmt.Println(fmt.Sprintf("log: isVirtual :%+v",inter))
			continue
		}
		// 大于en6的排除
		if strings.HasPrefix(inter.Name, "en") {
			numStr := inter.Name[2:]
			num, _ := strconv.Atoi(numStr)
			if num > 6 {
				//fmt.Println(fmt.Sprintf("log: is num>6 :%+v",inter))
				continue
			}
		}
		if strings.HasPrefix(inter.Name, "en") || strings.HasPrefix(inter.Name, "Ethernet") || strings.HasPrefix(inter.Name, "以太网") || strings.HasPrefix(inter.Name, "WLAN") {
			//fmt.Println(fmt.Sprintf("log: add :%+v",inter))
			macAddress = append(macAddress, hardwareAddr)
		} else if strings.HasPrefix(inter.Name, "Wi-Fi") || strings.HasPrefix(inter.Name, "无线网络") {
			wifiAddress = append(wifiAddress, hardwareAddr)
		} else if strings.HasPrefix(inter.Name, "Bluetooth") || strings.HasPrefix(inter.Name, "蓝牙网络连接") {
			bluetoothAddress = append(bluetoothAddress, hardwareAddr)
		} else {
			//fmt.Println(fmt.Sprintf("log: not add :%+v",inter))
		}
	}
	if len(macAddress) == 0 {
		macAddress = append(macAddress, wifiAddress...)
		if len(macAddress) == 0 {
			macAddress = append(macAddress, bluetoothAddress...)
		}
		if len(macAddress) == 0 {
			fmt.Printf(red, "no mac address found,Please contact customer service")
			_, _ = fmt.Scanln()
			return macErrorStr
		}
	}
	sort.Strings(macAddress)
	return fmt.Sprintf("%x", md5.Sum([]byte(strings.Join(macAddress, ","))))
}

func getMacMD5_241018() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("err:", err)
		return ""
	}

	var macAddress, bluetoothAddress, wifiAddress []string
	var macError []string

	//virtualMacPrefixes := []string{
	//	"00:05:69", "00:0C:29", "00:1C:14", "00:50:56", // VMware
	//	"00:15:5D",             // Hyper-V
	//	"08:00:27", "0A:00:27", // VirtualBox
	//}

	for _, inter := range interfaces {
		hardwareAddr := inter.HardwareAddr.String()
		if hardwareAddr == "" {
			continue
		}
		macError = append(macError, inter.Name+": "+hardwareAddr)

		//if isVirtualMac(hardwareAddr, virtualMacPrefixes) {
		//	continue
		//}

		switch {
		case inter.Name == "en0", inter.Name == "Ethernet0", inter.Name == "以太网":
			macAddress = append(macAddress, hardwareAddr)
		case inter.Name == "Bluetooth", inter.Name == "蓝牙网络连接":
			bluetoothAddress = append(bluetoothAddress, hardwareAddr)
		case inter.Name == "Wi-Fi", inter.Name == "WLAN", inter.Name == "无线网络":
			wifiAddress = append(wifiAddress, hardwareAddr)
		}
	}

	if len(macAddress) == 0 {
		macAddress = append(macAddress, bluetoothAddress...)
		if len(macAddress) == 0 {
			macAddress = append(macAddress, wifiAddress...)
		}
		if len(macAddress) == 0 {
			//fmt.Printf(red, "no mac address found,Please contact customer service")
			//_, _ = fmt.Scanln()
			//return macErrorStr
		}
	}
	sort.Strings(macError)
	return strings.Join(macError, "\n")
	//sort.Strings(macAddress)
	//return fmt.Sprintf("%x", md5.Sum([]byte(strings.Join(macAddress, ","))))
}
func getMacMD5_241019() string {
	id, err := machineid.ID()
	if err != nil {
		return err.Error()
	}
	id = strings.ToLower(id)
	id = strings.ReplaceAll(id, "-", "")
	return id
}

func printAD() {
	ad := client.GetAD()
	if len(ad) == 0 {
		return
	}
	fmt.Printf(yellow, ad)
}

func checkUpdate(version int) {
	upUrl := client.CheckVersion(fmt.Sprint(version))
	if upUrl == "" {
		return
	}
	isCopyText := ""
	installCmd := `bash -c "$(curl -fsSLk ` + githubPath + `install.sh)"`
	errClip := clipboard.WriteAll(installCmd)
	if errClip == nil {
		isCopyText = tr.Tr("（已复制到剪贴板）")
	}
	switch runtime.GOOS {
	case "windows":
		fmt.Printf(red, tr.Tr("有新版本，请关闭本窗口，将下面命令粘贴到GitBash窗口执行")+isCopyText+`：`)
	default:
		fmt.Printf(red, tr.Tr("有新版本，请关闭本窗口，将下面命令粘贴到新终端窗口执行")+isCopyText+`：`)
	}
	fmt.Printf(hGreen, installCmd)
	_, _ = fmt.Scanln()
	os.Exit(0)
	return
}

// 获取推广人
func getPromotion() (promotion string) {
	b, _ := os.ReadFile(os.Getenv("HOME") + "/.jetbrarc")
	promotion = strings.TrimSpace(string(b))
	if len(promotion) == 0 {
		if len(os.Args) > 1 {
			promotion = os.Args[1]
		}
	}
	return
}

func getLocale() (langRes, locRes string) {
	osHost := runtime.GOOS
	langRes = "en"
	locRes = "US"
	switch osHost {
	case "windows":
		// Exec powershell Get-Culture on Windows.
		cmd := exec.Command("powershell", "Get-Culture | select -exp Name")
		output, err := cmd.Output()
		if err == nil {
			langLocRaw := strings.TrimSpace(string(output))
			langLoc := strings.Split(langLocRaw, "-")
			langRes = langLoc[0]
			langRes = strings.Split(langRes, "-")[0]
			locRes = langLoc[1]
			return
		}
	case "darwin":
		// Exec shell Get-Culture on MacOS.
		cmd := exec.Command("sh", "osascript -e 'user locale of (get system info)'")
		output, err := cmd.Output()
		if err == nil {
			langLocRaw := strings.TrimSpace(string(output))
			langLoc := strings.Split(langLocRaw, "_")
			langRes = langLoc[0]
			langRes = strings.Split(langRes, "-")[0]
			if len(langLoc) == 1 {
				return
			}
			locRes = langLoc[1]
			return
		}
		plistB, err := os.ReadFile(os.Getenv("HOME") + "/Library/Preferences/.GlobalPreferences.plist")
		if err != nil {
			panic(err)
		}
		var a map[string]interface{}
		_, err = plist.Unmarshal(plistB, &a)
		if err != nil {
			panic(err)
		}
		langLocRaw := a["AppleLocale"].(string)
		langLoc := strings.Split(langLocRaw, "_")
		langRes = langLoc[0]
		langRes = strings.Split(langRes, "-")[0]
		if len(langLoc) == 1 {
			return
		}
		locRes = langLoc[1]
		return
	case "linux":
		envlang, ok := os.LookupEnv("LANG")
		if ok {
			langLocRaw := strings.TrimSpace(envlang)
			langLocRaw = strings.Split(envlang, ".")[0]
			langLoc := strings.Split(langLocRaw, "_")
			langRes = langLoc[0]
			langRes = strings.Split(langRes, "-")[0]
			if len(langLoc) == 1 {
				return
			}
			locRes = langLoc[1]
			return
		}
	}
	if langRes == "" {
		langLocRaw := os.Getenv("LC_CTYPE")
		langLocRaw = strings.Split(langLocRaw, ".")[0]
		langLoc := strings.Split(langLocRaw, "_")
		langRes = langLoc[0]
		langRes = strings.Split(langRes, "-")[0]
		if len(langLoc) == 1 {
			return
		}
		locRes = langLoc[1]
		return
	}
	return
}
