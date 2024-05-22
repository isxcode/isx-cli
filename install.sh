#!/bin/sh

ARCH=$(uname -m)
OS_TYPE=$OSTYPE

if [[ "$OS_TYPE" == "linux-gnu"* ]]; then
    if [ "$ARCH" == "x86_64" ]; then
        sudo curl -ssL https://isxcode.oss-cn-shanghai.aliyuncs.com/zhixingyun/isx_linux_amd64 -o /usr/local/bin/isx_linux_amd64
        sudo mv /usr/local/bin/isx_linux_amd64 /usr/local/bin/isx
        sudo chmod a+x /usr/local/bin/isx
    fi
elif [[ "$OS_TYPE" == "darwin"* ]]; then
    if [ "$ARCH" == "x86_64" ]; then
        sudo curl -ssL https://isxcode.oss-cn-shanghai.aliyuncs.com/zhixingyun/isx_darwin_amd64 -o /usr/local/bin/isx_darwin_amd64
        sudo mv /usr/local/bin/isx_darwin_amd64 /usr/local/bin/isx
        sudo chmod a+x /usr/local/bin/isx
    elif [ "$ARCH" == "aarch64" ]; then
        sudo curl -ssL https://isxcode.oss-cn-shanghai.aliyuncs.com/zhixingyun/isx_darwin_arm64 -o /usr/local/bin/isx_darwin_arm64
        sudo mv /usr/local/bin/isx_darwin_arm64 /usr/local/bin/isx
        sudo chmod a+x /usr/local/bin/isx
    fi
else
    echo "Unsupported OS or architecture: $OS_TYPE / $ARCH"
    exit 1
fi

echo "安装成功"