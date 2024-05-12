#!/bin/sh

sudo curl -ssL https://isxcode.oss-cn-shanghai.aliyuncs.com/zhixingyun/isx_darwin_arm64 -o /usr/local/bin/isx_darwin_arm64
sudo mv /usr/local/bin/isx_darwin_arm64 /usr/local/bin/isx
sudo chmod a+x /usr/local/bin/isx

echo "安装成功"