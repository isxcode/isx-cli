## isx-cli

> 开发isxcode组织项目的脚手架

- 项目管理
- 项目快速启动
- 项目初始化
- 代码自动生成
- Git分支管理

#### 前提

- 仅支持macOS系统

#### 安装

```bash
./install.sh
```

#### 下载最新版本isx二进制

> 若系统不支持，下载源码本地构建。

- [https://openfly.oss-cn-shanghai.aliyuncs.com/isx](https://openfly.oss-cn-shanghai.aliyuncs.com/isx)

#### 使用说明

```bash
# 查看版本
isx version

# 配置重置
isx reset

# 用户登录
isx login

# 退出登录
isx logout

# 下载代码
isx clone

# 查看所有项目列表
isx list

# 切换开发项目
isx code

# 查看当前开发项目
isx now

# 快速进入项目目录
isx home

# 使用idea开发项目
isx idea 

# 使用vscode开发项目
isx vscode

# 清理项目打包内容
isx clean

# 安装依赖
isx install

# 打包项目
isx package

# 启动项目
isx start 

# 单独启动前端
isx web

# 打包项目docker镜像
isx docker 

# 发布项目镜像
isx deploy 

# 启动本地官网
isx website

# 切开发#88需求分支
isx branch 88

# 获取编号分支信息
isx get 88

# 多项目git命令
isx git <command>

# 自动提交#88的pr
isx pr 88 '<message>'

# 删除本地项目
isx remove

# 格式化开发项目代码
isx format
```

```bash
isx login      ok 
isx logout     ok
isx reset      ok
isx clone      ok
isx list       ok
isx choose     ok
isx install    ok 
isx branch 88  
isx get 88 
isx start
isx package
isx commit ':sparkles: xxxx'
isx push 88
isx pr 88 '提交合并'
isx git command 
isx add project 
```

#### 未支持命令

```bash
# 初始化项目
isx init <app-name>
# 初始化前端模块代码
isx frontend <module-name>
# 初始化后端模块代码
isx backend <module-name>
```
