## JetBrains IDE 激活

[English](./README.md)

Jetbra Active 是一个 JetBrains IDE 激活工具，可以帮助你快速激活 JetBrains IDE
---
### 使用方式

在 MacOS/Linux 中，请打开终端；在 Windows 中，请打开 Git Bash。然后执行以下命令来安装：
> 部分电脑可能会误报毒，需要关闭杀毒软件/电脑管家/安全防护再进行

* 方式1：通过 ghp.ci 代理脚本
```shell
bash <(curl -Lk https://ghp.ci/https://github.com/kingparks/jetbra-activate/releases/download/latest/install.sh) githubReadme
```
* 方式2：通过 GitHub 脚本
```shell
bash <(curl -Lk https://github.com/kingparks/jetbra-activate/releases/download/latest/i.sh) githubReadme
```
* 方式3：手动下在二进制文件
> 从 [release](https://github.com/kingparks/jetbra-activate/releases) 页下载对应操作系统的二进制文件
 ```shell
# MaxOS/Linux
sudo mv jetbra_xx_xxx /usr/local/bin/jetbra;
chmod +x /usr/local/bin/jetbra;
jetbra githubReadme;
# Windows 
# double click jetbra_xx_xxx.exe
```

* 方式4：通过 go install 安装方式
```shell
go run github.com/kingparks/jetbra-activate@latest githubGoReadme;
```

---
### 功能特点

> 适用于IDEs全系列软件，如：IntelliJ IDEA、AppCode、CLion、DataGrip、GoLand、PhpStorm、PyCharm、Rider、RubyMine、WebStorm，当然也适用Windows/Mac/Linux平台。 此激活方法支持登录账号，支持在线更新，支持跨平台，支持最新版(正版激活)

![img_7.png](./img/img_2.png)

### Star History
![Star History Chart](https://api.star-history.com/svg?repos=kingparks/jetbra-activate&type=Date)