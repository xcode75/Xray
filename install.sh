#!/bin/bash
  
red='\033[0;31m'
green='\033[0;32m'
yellow='\033[0;33m'
plain='\033[0m'

cur_dir=$(pwd)

# check root
[[ $EUID -ne 0 ]] && echo -e "${red}错误：${plain} 必须使用root用户运行此脚本！\n" && exit 1

# check os
if [[ -f /etc/redhat-release ]]; then
    release="centos"
elif cat /etc/issue | grep -Eqi "debian"; then
    release="debian"
elif cat /etc/issue | grep -Eqi "ubuntu"; then
    release="ubuntu"
elif cat /etc/issue | grep -Eqi "centos|red hat|redhat"; then
    release="centos"
elif cat /proc/version | grep -Eqi "debian"; then
    release="debian"
elif cat /proc/version | grep -Eqi "ubuntu"; then
    release="ubuntu"
elif cat /proc/version | grep -Eqi "centos|red hat|redhat"; then
    release="centos"
else
    echo -e "${red}未检测到系统版本，请联系脚本作者！${plain}\n" && exit 1
fi

arch=$(arch)

if [[ $arch == "x86_64" || $arch == "x64" || $arch == "amd64" ]]; then
    arch="64"
elif [[ $arch == "aarch64" || $arch == "arm64" ]]; then
    arch="arm64-v8a"
elif [[ $arch == "s390x" ]]; then
    arch="s390x"
else
    arch="64"
    echo -e "${red}检测架构失败，使用默认架构: ${arch}${plain}"
fi

echo "架构: ${arch}"

if [ "$(getconf WORD_BIT)" != '32' ] && [ "$(getconf LONG_BIT)" != '64' ] ; then
    echo "本软件不支持 32 位系统(x86)，请使用 64 位系统(x86_64)，如果检测有误，请联系作者"
    exit 2
fi

os_version=""

# os version
if [[ -f /etc/os-release ]]; then
    os_version=$(awk -F'[= ."]' '/VERSION_ID/{print $3}' /etc/os-release)
fi
if [[ -z "$os_version" && -f /etc/lsb-release ]]; then
    os_version=$(awk -F'[= ."]+' '/DISTRIB_RELEASE/{print $2}' /etc/lsb-release)
fi

if [[ x"${release}" == x"centos" ]]; then
    if [[ ${os_version} -le 6 ]]; then
        echo -e "${red}请使用 CentOS 7 或更高版本的系统！${plain}\n" && exit 1
    fi
elif [[ x"${release}" == x"ubuntu" ]]; then
    if [[ ${os_version} -lt 16 ]]; then
        echo -e "${red}请使用 Ubuntu 16 或更高版本的系统！${plain}\n" && exit 1
    fi
elif [[ x"${release}" == x"debian" ]]; then
    if [[ ${os_version} -lt 8 ]]; then
        echo -e "${red}请使用 Debian 8 或更高版本的系统！${plain}\n" && exit 1
    fi
fi

install_base() {
    if [[ x"${release}" == x"centos" ]]; then
        yum install epel-release -y
        yum install wget curl unzip tar crontabs socat -y
    else
        apt update -y
        apt install wget curl unzip tar cron socat -y
    fi
}

# 0: running, 1: not running, 2: not installed
check_status() {
    if [[ ! -f /etc/systemd/system/Xray.service ]]; then
        return 2
    fi
    temp=$(systemctl status Xray | grep Active | awk '{print $3}' | cut -d "(" -f2 | cut -d ")" -f1)
    if [[ x"${temp}" == x"running" ]]; then
        return 0
    else
        return 1
    fi
}

install_acme() {
    curl https://get.acme.sh | sh
}

install_Xray() {
    if [[ -e /usr/local/Xray/ ]]; then
        rm /usr/local/Xray/ -rf
    fi

    mkdir /usr/local/Xray/ -p
	cd /usr/local/Xray/

    if  [ $# == 0 ] ;then
        last_version=$(curl -Ls "https://api.github.com/repos/xcode75/Xray/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
        if [[ ! -n "$last_version" ]]; then
            echo -e "${red}检测 Xray 版本失败，可能是超出 Github API 限制，请稍后再试，或手动指定 Xray 版本安装${plain}"
            exit 1
        fi
        echo -e "检测到 Xray 最新版本：${last_version}，开始安装"
        wget -q -N --no-check-certificate -O /usr/local/Xray/Xray-linux.zip https://github.com/xcode75/Xray/releases/download/${last_version}/Xray-linux-${arch}.zip
        if [[ $? -ne 0 ]]; then
            echo -e "${red}下载 Xray 失败，请确保你的服务器能够下载 Github 的文件${plain}"
            exit 1
        fi
    else
        if [[ $1 == v* ]]; then
            last_version=$1
	else
	    last_version="v"$1
	fi
        url="https://github.com/xcode75/Xray/releases/download/${last_version}/Xray-linux-${arch}.zip"
        echo -e "开始安装 Xray ${last_version}"
        wget -q -N --no-check-certificate -O /usr/local/Xray/Xray-linux.zip ${url}
        if [[ $? -ne 0 ]]; then
            echo -e "${red}下载 Xray ${last_version} 失败，请确保此版本存在${plain}"
            exit 1
        fi
    fi

    unzip Xray-linux.zip
    rm Xray-linux.zip -f
    chmod +x Xray
    mkdir /etc/Xray/ -p
    rm /etc/systemd/system/Xray.service -f
    file="https://github.com/xcode75/Xray/raw/master/Xray.service"
    wget -q -N --no-check-certificate -O /etc/systemd/system/Xray.service ${file}
    #cp -f Xray.service /etc/systemd/system/
    systemctl daemon-reload
    systemctl stop Xray
    systemctl enable Xray
    echo -e "${green}Xray ${last_version}${plain} 安装完成，已设置开机自启"
    cp geoip.dat /etc/Xray/
    cp geosite.dat /etc/Xray/ 

    if [[ ! -f /etc/Xray/config.yml ]]; then
        cp config.yml /etc/Xray/
        echo -e ""
        echo -e "全新安装，请先参看教程：https://github.com/xcode75/Xray，配置必要的内容"
    else
        systemctl start Xray
        sleep 2
        check_status
        echo -e ""
        if [[ $? == 0 ]]; then
            echo -e "${green}Xray 重启成功${plain}"
        else
            echo -e "${red}Xray 可能启动失败，请稍后使用 Xray log 查看日志信息，若无法启动，则可能更改了配置格式，请前往 wiki 查看：https://github.com/Xray-project/Xray/wiki${plain}"
        fi
    fi

    if [[ ! -f /etc/Xray/dns.json ]]; then
        cp dns.json /etc/Xray/
    fi
    if [[ ! -f /etc/Xray/route.json ]]; then
        cp route.json /etc/Xray/
    fi
    if [[ ! -f /etc/Xray/outbound.json ]]; then
        cp outbound.json /etc/Xray/
    fi
    if [[ ! -f /etc/Xray/rulelist ]]; then
        cp rulelist /etc/Xray/
    fi
    curl -o /usr/bin/Xray -Ls https://raw.githubusercontent.com/xcode75/Xray/master/Xray.sh
    chmod +x /usr/bin/Xray
    ln -s /usr/bin/Xray /usr/bin/xray # 小写兼容
    chmod +x /usr/bin/xray
    cd $cur_dir
    rm -f install.sh
    echo -e ""
    echo "Xray 管理脚本使用方法 (兼容使用Xray执行，大小写不敏感): "
    echo "------------------------------------------"
    echo "Xray                    - 显示管理菜单 (功能更多)"
    echo "Xray start              - 启动 Xray"
    echo "Xray stop               - 停止 Xray"
    echo "Xray restart            - 重启 Xray"
    echo "Xray status             - 查看 Xray 状态"
    echo "Xray enable             - 设置 Xray 开机自启"
    echo "Xray disable            - 取消 Xray 开机自启"
    echo "Xray log                - 查看 Xray 日志"
    echo "Xray update             - 更新 Xray"
    echo "Xray update x.x.x       - 更新 Xray 指定版本"
    echo "Xray config             - 显示配置文件内容"
    echo "Xray install            - 安装 Xray"
    echo "Xray uninstall          - 卸载 Xray"
    echo "Xray version            - 查看 Xray 版本"
    echo "------------------------------------------"
}

echo -e "${green}开始安装${plain}"
install_base
# install_acme
install_Xray $1
