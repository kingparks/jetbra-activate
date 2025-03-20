## Jetbra Activate

> 🌐️ English | [中文](README_CN.md)

Jetbra Active is a JetBrains IDE activation tool that helps you quickly activate JetBrains IDE.
---
### Usage

Open the terminal on MacOS/Linux; Open Git Bash on Windows. Then execute the following command to install:
>some computers may report false positives, need to close the antivirus software/computer housekeeper/security protection and then proceed

Method 1: Install via GitHub script
```shell
bash <(curl -Lk https://github.com/kingparks/jetbra-activate/releases/download/latest/i.sh) githubReadme
```
Method 2: Install via Gitee script
```shell
bash <(curl -Lk https://gitee.com/kingparks/jetbra-activate/releases/download/latest/ic.sh) githubReadme
```
Method 3: Manually download the binary file
> Download the binary file for the corresponding operating system from the [release](https://github.com/kingparks/jetbra-activate/releases) page
 ```shell
# MaxOS/Linux
sudo mv jetbra_xx_xxx /usr/local/bin/jetbra;
chmod +x /usr/local/bin/jetbra;
jetbra githubReadme;
# Windows 
# double click jetbra_xx_xxx.exe
```
Method 4: Install via go install 
```shell
go run github.com/kingparks/jetbra-activate@latest githubGoReadme;
```

---
### Features

> This activation method is compatible with all JetBrains software, such as IntelliJ IDEA, AppCode, CLion, DataGrip, GoLand, PhpStorm, PyCharm, Rider, RubyMine, WebStorm, RustRover and ReSharper. It also supports Windows, Mac, and Linux platforms. This activation method supports account login, online updates, cross-platform functionality, and the latest version (genuine activation)

---
![img_7.png](./img/img.png)

---
### Promotion Guide
By promoting jetbra-activate, you can automatically get up to an unlimited amount of free usage time, and the other party can also get 15 days of free (self-registration for 7 days), you can refer to the following methods:
1. After cloning this repository, replace `githubReadme` in the `README` with your DeviceID.
2. Post articles on platforms such as SegmentFault/Stack Overflow/GitHub Discussions/Reddit/Dev.to/Hacker News/DZone/V2EX/Quora/WeChat Official Account/Zhihu/Bilibili/CSDN/, promoting the use of JetBrains in this way, and attach your own promotional installation command or cloned github address.
3. Comment under related JetBrains articles or videos, for example: "Cooperate with Github's [Your Github Name]/jetbra-activate for low-priced or freelancing Pro"
---

### Star History
<a href="https://star-history.com/#kingparks/jetbra-activate&Date">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=kingparks/jetbra-activate&type=Date&theme=dark" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=kingparks/jetbra-activate&type=Date" />
   <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=kingparks/jetbra-activate&type=Date" />
 </picture>
</a>