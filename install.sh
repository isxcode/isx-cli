#!/bin/sh

sudo wget https://github.com/isxcode/isx-cli/releases/download/v1.1.0/isx_darwin_arm64 -O /usr/local/bin/isx_darwin_arm64
sudo mv /usr/local/bin/isx_darwin_arm64 /usr/local/bin/isx
sudo chmod a+x /usr/local/bin/isx
/usr/local/bin/isx set version.number 1.1.0

echo "安装成功"