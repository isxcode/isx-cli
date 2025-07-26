# 至爻数据开发规范脚手架

##### 在线安装

```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/isxcode/isx-cli/main/install.sh)"
```

##### 国内安装

```bash
sh -c "$(curl -fsSL https://gitee.com/isxcode/isx-cli/raw/main/install.sh)"
```

##### 安装包下载

- [isx_linux_amd64](https://isxcode.oss-cn-shanghai.aliyuncs.com/zhixingyun/isx_linux_amd64)
- [isx_windows_amd64.exe](https://isxcode.oss-cn-shanghai.aliyuncs.com/zhixingyun/isx_windows_amd64.exe)
- [isx_windows_arm64.exe](https://isxcode.oss-cn-shanghai.aliyuncs.com/zhixingyun/isx_windows_arm64.exe)
- [isx_darwin_amd64](https://isxcode.oss-cn-shanghai.aliyuncs.com/zhixingyun/isx_darwin_amd64)
- [isx_darwin_arm64](https://isxcode.oss-cn-shanghai.aliyuncs.com/zhixingyun/isx_darwin_arm64)

##### 开发流程

```bash
# 1.登陆github账号
isx login
# 2.下载项目
isx clone
# 3.选择开发项目
isx choose
# 4.选择任务
isx issue
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
# 命令详解
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

至行云-至爻数据开发规范脚手架
代码仓库：https://github.com/isxcode/isx-cli

Available Commands:
  backend     isx backend                             | 本地启动后端
  backup      isx backup <comment>                    | 备份项目资源
  checkout    isx checkout <issue_number>             | 切换分支
  choose      isx choose                              | 切换项目
  clean       isx clean                               | 清除项目缓存
  clone       isx clone                               | 下载代码
  config      isx config                              | 查看脚手架配置
  delete      isx delete <issue_number>               | 删除远程分支
  docker      isx docker                              | 构建Docker镜像
  format      isx format                              | 代码格式化
  frontend    isx frontend                            | 本地启动前端
  git         isx git <git_command>                   | 执行git命令
  gradle      isx gradle <gradle_command>             | 执行gradle命令
  install     isx install                             | 安装依赖
  issue       isx issue                               | 选择任务
  login       isx login                               | 用户登录
  logout      isx logout                              | 退出登录
  now         isx now                                 | 查看当前信息
  package     isx package                             | 源码打包
  pr          isx pr <issue_number>                   | 提交pr
  pull        isx pull                                | 拉取代码
  push        isx push                                | 提交代码
  remove      isx remove                              | 删除本地项目
  rollback    isx rollback                            | 回滚项目资源
  set         isx set <config_key> <value>            | 修改脚手架配置
  start       isx start                               | 启动项目
  upgrade     isx upgrade                             | 升级脚手架
  upload      isx upload <target>                     | 发布本地安装包
  version     isx version                             | 查看脚手架版本
  website     isx website                             | 本地启动官网
```