#!/bin/sh

ARCH=$(uname -m)
OS_TYPE=$OSTYPE

echo "需要sudo授权"
sudo echo "授权成功"

if [[ "$OS_TYPE" == "linux-gnu"* ]]; then
    if [ "$ARCH" == "x86_64" ]; then
        echo "开始下载"
        sudo curl -ssL https://isxcode.oss-cn-shanghai.aliyuncs.com/zhixingyun/isx_linux_amd64 -o /usr/bin/isx_linux_amd64
        echo "下载完成"
        sudo mv /usr/bin/isx_linux_amd64 /usr/bin/isx
        sudo chmod a+x /usr/bin/isx
        echo "安装成功"
    fi
elif [[ "$OS_TYPE" == "darwin"* ]]; then
    if [ "$ARCH" == "x86_64" ]; then
        echo "开始下载"
        sudo curl -ssL https://isxcode.oss-cn-shanghai.aliyuncs.com/zhixingyun/isx_darwin_amd64 -o /usr/local/bin/isx_darwin_amd64
        echo "下载完成"
        sudo mv /usr/local/bin/isx_darwin_amd64 /usr/local/bin/isx
        sudo chmod a+x /usr/local/bin/isx
        echo "安装成功"
    elif [ "$ARCH" == "arm64" ]; then
        echo "开始下载"
        sudo curl -ssL https://isxcode.oss-cn-shanghai.aliyuncs.com/zhixingyun/isx_darwin_arm64 -o /usr/local/bin/isx_darwin_arm64
        echo "下载完成"
        sudo mv /usr/local/bin/isx_darwin_arm64 /usr/local/bin/isx
        sudo chmod a+x /usr/local/bin/isx
        echo "安装成功"
    fi
else
    echo "Unsupported OS or architecture: $OS_TYPE / $ARCH"
    exit 1
fi