# 至行云

### 团队开发规范脚手架

##### 安装包下载

- [isx_linux_amd64](https://isxcode.oss-cn-shanghai.aliyuncs.com/zhixingyun/isx_linux_amd64)
- [isx_windows_amd64.exe](https://isxcode.oss-cn-shanghai.aliyuncs.com/zhixingyun/isx_windows_amd64.exe)
- [isx_windows_arm64.exe](https://isxcode.oss-cn-shanghai.aliyuncs.com/zhixingyun/isx_windows_arm64.exe)
- [isx_darwin_amd64](https://isxcode.oss-cn-shanghai.aliyuncs.com/zhixingyun/isx_darwin_amd64)
- [isx_darwin_arm64](https://isxcode.oss-cn-shanghai.aliyuncs.com/zhixingyun/isx_darwin_arm64)

##### 在线安装

```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/isxcode/isx-cli/main/install.sh)"
```

##### 国内阿里云安装

```bash
sh -c "$(curl -fsSL https://isxcode.oss-cn-shanghai.aliyuncs.com/zhixingyun/install.sh)"
```

##### 开发流程

```bash
# 1.登陆github账号
isx login
# 2.下载项目
isx clone
# 3.选择开发项目
isx choose
# 4.切换开发分支
isx checkout <issue_number>
# 5.格式化代码
isx format
# 6.提交代码到本地
isx git add .
isx git commit -m "your commit message"
# 7.推送到origin仓库
isx push
# 8.提交pr
isx pr <issue_number>
# 9.如果提交的pr，无法成功rebase合并
isx pull 
# 10.修复冲突,重新提交推送代码
isx push -f
```

##### 使用说明

```bash
# 查看所有命令
isx -h
# 命令详解†
isx login -h
```

```text
 ____ _____ __ __           __  _      ____ 
|    / ___/|  |  |         /  ]| |    |    |
 |  (   \_ |  |  | _____  /  / | |     |  | 
 |  |\__  ||_   _||     |/  /  | |___  |  | 
 |  |/  \ ||     ||_____/   \_ |     | |  | 
 |  |\    ||  |  |      \     ||     | |  | 
|____|\___||__|__|       \____||_____||____|

欢迎使用isx-cli脚手架
代码仓库：https://github.com/isxcode/isx-cli

Usage:
  isx [command]

Available Commands:
  autotest    isx autotest                                                     | 自动化测试
  build       isx build                                                        | 使用docker编译项目代码
  checkout    isx checkout <issue_number>                                      | 切换开发分支
  choose      isx choose                                                       | 切换开发项目
  clean       isx clean                                                        | 删除项目缓存
  clone       isx clone                                                        | 下载项目代码
  completion  Generate the autocompletion script for the specified shell
  config      isx config                                                       | 查看配置
  db          isx db list | isx db <issue_number>                              | 查看当前db(暂不开放)
  delete      isx delete <issue_number>                                        | 删除组织分支
  fork        isx fork                                                         | Fork当前项目为同名个人仓库
  format      isx format                                                       | 格式化代码
  git         isx git <git command>                                            | 项目内执行git命令
  gradle      isx gradle <gradle_command>                                      | 项目内执行gradle命令
  help        Help about any command
  install     isx install                                                      | 安装项目依赖
  issue       isx issue                                                        | 列出当前仓库分配给您的issue
  login       isx login                                                        | 登录github账号
  now         isx now                                                          | 查看项目信息
  package     isx package                                                      | 源码编译打包
  pr          isx pr <issue_number>                                            | 提交pr
  pull        isx pull                                                         | 拉取组织代码
  push        isx push                                                         | 格式化代码后,提交代码
  remove      isx remove                                                       | 删除本地项目
  run         isx run [frontend/backend/web] [port]                            | 使用docker运行项目
  server      isx server                                                       | 本地启动后端程序
  set         isx set <config_key> <value>                                     | 设置配置参数
  start       isx start                                                        | 启动项目
  sync        isx sync <branch_name>                                           | 同步Github个人仓库指定分支
  upgrade     isx upgrade                                                      | 升级isx-cli脚手架
  version     isx version                                                      | 查看版本号
  web         isx web                                                          | 本地启动前端服务
  website     isx website                                                      | 本地启动官网

Flags:
      --config string   config file (default is $HOME/.isx/isx-config.yml)
  -h, --help            help for isx

Use "isx [command] --help" for more information about a command.
```