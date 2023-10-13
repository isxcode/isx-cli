#!/bin/bash

# 安装依赖
python3 -m pip install -r requirements.txt

# 打包
python3 -m nuitka --follow-imports --output-dir=./build --output-filename=isx main.py