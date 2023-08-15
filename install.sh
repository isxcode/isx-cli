#!/bin/bash

# 获取当前路径
BASE_PATH=$(cd "$(dirname "$0")" || exit ; pwd)

# 安装nuitka
python3 -m pip install nuitka

# 安装依赖
python3 -m pip install -r requirements.txt

# 清空打包文件
rm -rf ${BASE_PATH}/build
rm -f isx

# 打包
python3 -m nuitka --output-dir=./build --output-filename=isx main.py

# 配置环境变量
echo 'export PATH=$PATH:'"${BASE_PATH}" >> ~/.zshrc

# 生效配置
source ~/.zshrc