## isx-cli工具

> 开发isxcode组织项目的脚手架。

- 支持项目初始化
- 支持模块代码自动生成
- 支持Git分支管理
- 支持项目快速启动

#### 前提

- 仅限macOS系统
- 安装zsh

#### 安装

```bash
zsh install.sh
```

#### 使用

```bash
# 工具配置重置
isx reset
# 用户登录
isx login
# 退出登录
isx logout
# 下载代码
isx clone
# 查看所有项目列表
isx list
# 切换项目
isx develop spark-yun
# 查看当前开发项目
isx show
# 使用idea打开项目
isx idea 
# 使用vscode打开项目
isx vscode
# 清理项目打包内容
isx clean 
# 启动项目
isx start 
# 打包项目
isx package
# 打包项目镜像
isx docker 
# 发布项目镜像
isx deploy 
# 启动本地官网
isx website
# 创建开发需求分支
isx branch #89
# 多项目
isx git <command>
# 自动提交pr
isx pr 88 '<message>'
# 初始化项目
isx init <app-name>
# 初始化前端模块
isx frontend <module-name>
# 初始化后端模块
isx backend <module-name>
```