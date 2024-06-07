package main

import (
	"crypto/md5"
	"embed"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/atotto/clipboard"
	"github.com/tidwall/gjson"
	"howett.net/plist"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"time"
)

var version = 212

var hosts = []string{"http://string.jeter.eu.org", "http://string.eiyou.fun", "http://jetbra.serv00.net:7191", "http://ba.serv00.net:7191"}
var host = hosts[0]
var githubPath = "https://mirror.ghproxy.com/https://github.com/kingparks/jetbra-activate/releases/download/latest/"
var err error

var green = "\033[32m%s\033[0m\n"
var yellow = "\033[33m%s\033[0m\n"
var hGreen = "\033[1;32m%s\033[0m"
var dGreen = "\033[4;32m%s\033[0m\n"
var red = "\033[31m%s\033[0m\n"
var defaultColor = "%s"
var lang, _ = getLocale()

//go:embed all:script
var scriptFS embed.FS

func main() {
	setProxy()
	switch lang {
	case "zh":
		fmt.Printf(green, `IntelliJ 授权 v`+strings.Join(strings.Split(fmt.Sprint(version), ""), "."))
	default:
		fmt.Printf(green, `IntelliJ Activate v`+strings.Join(strings.Split(fmt.Sprint(version), ""), "."))
	}
	checkHost()
	deviceID := getMacMD5()
	switch lang {
	case "zh":
		fmt.Printf(green, "设备码: "+deviceID)
	default:
		fmt.Printf(green, "Device ID: "+deviceID)
	}
	printAD()
	checkUpdate(version)
	fmt.Println()
	switch lang {
	case "zh":
		fmt.Printf(defaultColor, "选择要授权的产品：")
	default:
		fmt.Printf(defaultColor, "Choose the product to authorize: ")
	}
	jbProduct := []string{"IntelliJ IDEA", "CLion", "PhpStorm", "Goland", "PyCharm", "WebStorm", "Rider", "DataGrip", "DataSpell"}
	jbProductChoice := []string{"idea", "clion", "phpstorm", "goland", "pycharm", "webstorm", "rider", "datagrip", "dataspell"}
	for i, v := range jbProduct {
		fmt.Printf(hGreen, fmt.Sprintf("%d. %s\t", i+1, v))
	}
	fmt.Println()
	switch lang {
	case "zh":
		fmt.Print("请输入产品编号（直接回车默认为1）：")
	default:
		fmt.Print("Please enter the product number (press Enter directly to default to 1): ")
	}
	productIndex := 1
	_, _ = fmt.Scanln(&productIndex)
	if productIndex < 1 || productIndex > len(jbProduct) {
		switch lang {
		case "zh":
			fmt.Println("输入有误")
		default:
			fmt.Println("Input error")
		}
		return
	}
	switch lang {
	case "zh":
		fmt.Println("选择的产品为：" + jbProduct[productIndex-1])
	default:
		fmt.Println("Selected product：" + jbProduct[productIndex-1])
	}
	fmt.Println()
	switch lang {
	case "zh":
		fmt.Printf(defaultColor, "选择有效期：")
	default:
		fmt.Printf(defaultColor, "Choose the validity period: ")
	}
	jbPeriod := []string{"两天(免费)", "一年(赠送)"}
	switch lang {
	case "zh":
	default:
		jbPeriod = []string{"Two days (free)", "One year (gift)"}
	}
	for i, v := range jbPeriod {
		fmt.Printf(hGreen, fmt.Sprintf("%d. %s\t", i+1, v))
	}
	fmt.Println()
	switch lang {
	case "zh":
		fmt.Printf("%s", "请输入有效期编号（直接回车默认为1）：")
	default:
		fmt.Printf("%s", "Please enter the validity period number (press Enter directly to default to 1): ")
	}
	periodIndex := 1
	_, _ = fmt.Scanln(&periodIndex)
	if periodIndex < 1 || periodIndex > len(jbPeriod) {
		switch lang {
		case "zh":
			fmt.Println("输入有误")
		default:
			fmt.Println("Input error")
		}
		return
	}
	switch lang {
	case "zh":
		fmt.Println("选择的有效期为：" + jbPeriod[periodIndex-1])
	default:
		fmt.Println("Selected validity period：" + jbPeriod[periodIndex-1])
	}

	fmt.Println()
	lic := ""
	for i := 0; i < 50; i++ {
		if i == 20 {
			Clean()
			Active(jbProductChoice[productIndex-1])
			lic = ReadLic(jbProductChoice[productIndex-1])
		}
		time.Sleep(20 * time.Millisecond)
		h := strings.Repeat("=", i) + strings.Repeat(" ", 49-i)
		fmt.Printf("\r%.0f%%[%s]", float64(i)/49*100, h)
	}
	fmt.Println()
	fmt.Println()

	switch periodIndex {
	case 1:
		isOk, result := getLic(jbProductChoice[productIndex-1], periodIndex-1)
		if !isOk {
			fmt.Printf(red, result)
			return
		}
		lic = result
	case 2:
		// 获取之前的 lic
		if lic != "" {
			goto Process
		}
		switch lang {
		case "zh":
			fmt.Println("该项目由捐赠字符串藏品后赠送，捐赠获取藏品地址：")
			fmt.Printf(dGreen, host)
			fmt.Println("请输入字符串藏品元数据：")
		default:
			fmt.Println("This project is donated by the string collection and then given away. The donation address of the collection is: ")
			fmt.Printf(dGreen, host)
			fmt.Println("Please enter the string collection metadata: ")
		}
		var ticket string
		_, _ = fmt.Scanln(&ticket)
		if ticket == "" {
			switch lang {
			case "zh":
				fmt.Printf(red, "输入有误")
			default:
				fmt.Printf(red, "Input error")
			}
			return
		}
		ticket = strings.TrimSpace(ticket)
		ticket = strings.ReplaceAll(ticket, "https://idea.eiyou.fun/", "")
		ticket = strings.ReplaceAll(ticket, "http://idea.eiyou.fun/", "")
		ticket = strings.ReplaceAll(ticket, "https://idea.jeter.eu.org/", "")
		ticket = strings.ReplaceAll(ticket, "http://idea.jeter.eu.org/", "")
		ticket = strings.ReplaceAll(ticket, "/", "")
		// 检查是否捐赠
		res, err := httplib.Get(host + "/check?ticket=" + ticket + "&device=" + deviceID).String()
		if err != nil {
			fmt.Printf(red, err.Error())
			return
		}
		if gjson.Get(res, "code").Int() != 0 {
			fmt.Printf(red, gjson.Get(res, "msg").String())
			return
		}
		isOk, result := getLic(jbProductChoice[productIndex-1], periodIndex-1)
		if !isOk {
			fmt.Printf(red, result)
			return
		}
		lic = result
		// 保存到本地
		WriteLic(jbProductChoice[productIndex-1], lic)
		fmt.Println()
	}

Process:
	isCopyText := ""
	err = clipboard.WriteAll(lic)
	if err == nil {
		switch lang {
		case "zh":
			isCopyText = "（已复制到剪贴板）"
		default:
			isCopyText = "（Copied to clipboard）"
		}
	}
	switch lang {
	case "zh":
		fmt.Printf(yellow, "首次执行请重启IDE，然后填入下面授权码；非首次执行直接填入下面授权码即可"+isCopyText)
	default:
		fmt.Printf(yellow, "For the first execution, please restart the IDE, then fill in the authorization code below; for non-first time execution, fill in the authorization code directly"+isCopyText)
	}
	fmt.Println()
	fmt.Printf(hGreen, lic)
	fmt.Println()
}

func getMacMD5() string {
	// 获取本机的MAC地址
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("err:", err)
		return ""
	}
	var macAddress []string
	for _, inter := range interfaces {
		if strings.HasPrefix(inter.Name, "en") || strings.HasPrefix(inter.Name, "Ethernet") || strings.HasPrefix(inter.Name, "以太网") || strings.HasPrefix(inter.Name, "本地连接") || strings.HasPrefix(inter.Name, "WLAN") {
			macAddress = append(macAddress, inter.HardwareAddr.String())
		}
	}
	sort.Strings(macAddress)
	return fmt.Sprintf("%x", md5.Sum([]byte(strings.Join(macAddress, ","))))
}

func printAD() {
	res, err := httplib.Get(host + "/ad").String()
	if err != nil {
		return
	}
	if len(res) == 0 {
		return
	}
	fmt.Printf(yellow, res)
}

func checkUpdate(version int) {
	res, err := httplib.Get(host + "/version?version=" + fmt.Sprint(version) + "&plat=" + runtime.GOOS + "_" + runtime.GOARCH).String()
	if err != nil {
		switch lang {
		case "zh":
			fmt.Printf(red, "网络错误！\n"+err.Error())
			fmt.Printf(red, "排除网络问题后，仍存在问题的话可尝试下载新版本:\n")
		default:
			fmt.Printf(red, "Network error！\n"+err.Error())
			fmt.Printf(red, "If there is still a problem after excluding the network problem, you can try to download the new version:\n")
		}
		fmt.Printf(dGreen, githubPath+"jetbra_"+runtime.GOOS+"_"+runtime.GOARCH)
		os.Exit(0)
		return
	}
	upUrl := gjson.Get(res, "url").String()
	if upUrl != "" {
		switch lang {
		case "zh":
			fmt.Printf(red, "有新版本可用，尝试自动更新中，若失败，请输入下面命令并回车手动更新程序:\n")
		default:
			fmt.Printf(red, "There is a new version available, trying to update automatically. If it fails, please enter the following command and press Enter to update the program manually:\n")
		}
		fmt.Println(`bash -c "$(curl -fsSL ` + host + `/app/install.sh)"`)
		cmd := exec.Command("bash", "-c", fmt.Sprintf(`bash -c "$(curl -fsSL %s/app/install.sh)"`, host))
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		switch lang {
		case "zh":
			fmt.Println("更新完成，重新运行程序即可")
		default:
			fmt.Println("Update completed, just run the program again")
		}
		os.Exit(0)
		return
	}
}

func setProxy() {
	lang经由 := "经由"
	lang代理访问 := "代理访问 "
	switch lang {
	case "zh":
	default:
		lang经由 = "via"
		lang代理访问 = "proxy access "
	}
	httplib.SetDefaultSetting(httplib.BeegoHTTPSettings{
		ReadWriteTimeout: 30 * time.Second,
		ConnectTimeout:   30 * time.Second,
		Gzip:             true,
		DumpBody:         true,
		UserAgent:        lang,
	})
	if os.Getenv("http_proxy") != "" {
		httplib.SetDefaultSetting(httplib.BeegoHTTPSettings{
			Proxy: func(request *http.Request) (*url.URL, error) {
				return url.Parse(os.Getenv("http_proxy"))
			},
			ReadWriteTimeout: 30 * time.Second,
			ConnectTimeout:   30 * time.Second,
			Gzip:             true,
			DumpBody:         true,
			UserAgent:        lang,
		})
		fmt.Printf(yellow, lang经由+" http_proxy "+lang代理访问+os.Getenv("http_proxy"))
		return
	}
	if os.Getenv("https_proxy") != "" {
		httplib.SetDefaultSetting(httplib.BeegoHTTPSettings{
			Proxy: func(request *http.Request) (*url.URL, error) {
				return url.Parse(os.Getenv("https_proxy"))
			},
			ReadWriteTimeout: 30 * time.Second,
			ConnectTimeout:   30 * time.Second,
			Gzip:             true,
			DumpBody:         true,
			UserAgent:        lang,
		})
		fmt.Printf(yellow, lang经由+" https_proxy "+lang代理访问+os.Getenv("https_proxy"))
		return
	}
	if os.Getenv("all_proxy") != "" {
		httplib.SetDefaultSetting(httplib.BeegoHTTPSettings{
			Proxy: func(request *http.Request) (*url.URL, error) {
				return url.Parse(os.Getenv("all_proxy"))
			},
			ReadWriteTimeout: 30 * time.Second,
			ConnectTimeout:   30 * time.Second,
			Gzip:             true,
			DumpBody:         true,
			UserAgent:        lang,
		})
		fmt.Printf(yellow, lang经由+" all_proxy "+lang代理访问+os.Getenv("all_proxy"))
		return
	}
}

func getLocale() (lang, loc string) {
	osHost := runtime.GOOS
	defaultLang := "en"
	defaultLoc := "US"
	switch osHost {
	case "windows":
		// Exec powershell Get-Culture on Windows.
		cmd := exec.Command("powershell", "Get-Culture | select -exp Name")
		output, err := cmd.Output()
		if err == nil {
			langLocRaw := strings.TrimSpace(string(output))
			langLoc := strings.Split(langLocRaw, "-")
			lang := langLoc[0]
			lang = strings.Split(lang, "-")[0]
			loc := langLoc[1]
			return lang, loc
		}
	case "darwin":
		// Exec shell Get-Culture on MacOS.
		cmd := exec.Command("sh", "osascript -e 'user locale of (get system info)'")
		output, err := cmd.Output()
		if err == nil {
			langLocRaw := strings.TrimSpace(string(output))
			langLoc := strings.Split(langLocRaw, "_")
			lang := langLoc[0]
			lang = strings.Split(lang, "-")[0]
			if len(langLoc) == 1 {
				return lang, defaultLoc
			}
			loc := langLoc[1]
			return lang, loc
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
		lang := langLoc[0]
		lang = strings.Split(lang, "-")[0]
		if len(langLoc) == 1 {
			return lang, defaultLoc
		}
		loc := langLoc[1]
		return lang, loc
	case "linux":
		envlang, ok := os.LookupEnv("LANG")
		if ok {
			langLocRaw := strings.TrimSpace(envlang)
			langLocRaw = strings.Split(envlang, ".")[0]
			langLoc := strings.Split(langLocRaw, "_")
			lang := langLoc[0]
			lang = strings.Split(lang, "-")[0]
			if len(langLoc) == 1 {
				return lang, defaultLoc
			}
			loc := langLoc[1]
			return lang, loc
		}
	}
	if lang == "" {
		langLocRaw := os.Getenv("LC_CTYPE")
		langLocRaw = strings.Split(langLocRaw, ".")[0]
		langLoc := strings.Split(langLocRaw, "_")
		lang := langLoc[0]
		lang = strings.Split(lang, "-")[0]
		if len(langLoc) == 1 {
			return lang, defaultLoc
		}
		loc := langLoc[1]
		return lang, loc
	}
	return defaultLang, defaultLoc
}

func checkHost() {
	for _, v := range hosts {
		_, err := httplib.Get(v).SetTimeout(4*time.Second, 4*time.Second).String()
		if err == nil {
			host = v
			return
		}
	}
}

func getLic(product string, dur int) (isOk bool, result string) {
	req := httplib.Get(host + "/getLic?device=" + getMacMD5() + "&dur=" + fmt.Sprint(dur) + "&product=" + product)
	res, err := req.String()
	if err != nil {
		isOk = false
		result = err.Error()
		return
	}
	code := gjson.Get(res, "code").Int()
	msg := gjson.Get(res, "msg").String()
	result = msg
	if code != 0 {
		isOk = false
		return
	}
	isOk = true
	return
}
