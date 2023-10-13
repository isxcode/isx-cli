## isx-cli

> 开发isxcode组织项目的脚手架

- 项目管理
- 项目快速启动
- 项目初始化
- 代码自动生成
- Git分支管理

#### 安装

```bash
./install.sh
```

#### 下载最新版本isx二进制

> 若系统不支持，下载源码本地构建。

- [https://openfly.oss-cn-shanghai.aliyuncs.com/isx/isx_macos_m2_v1.0.0](https://openfly.oss-cn-shanghai.aliyuncs.com/isx/isx_macos_m2_v1.0.0)

#### 使用说明

```bash
# 重新用户旧配置
isx reset

# 用户登录
isx login

# 退出登录或者token过期替换
isx logout

# 下载项目代码准备开发
isx clone

# 查看本地项目列表情况
isx list

# 选择开发项目
isx choose

# 根据缺陷编号，切换开发分支
isx branch 88

# 安装项目依赖，准备启动项目
isx install

# 启动项目
isx start

# 项目打包
isx package

# 提交代码前，一定要先格式化代码
isx format

# 获取缺陷编号的真实分支名称，留给提交分支使用
isx get 88

# 执行相关git命令，会在每个字模块中执行一次，减少重复操作
isx <git-command>
isx git add .
isx git commit -m ":sparkles: 添加需求"
isx git push origin 0.0.7-#88
isx git push origin 0.0.7-#88 -f
isx git rebase upstream/0.0.7-#88

# 一剑提交pr，自动提交子模块的pr，减少重复操作
isx pr 88 '提交pr内容'
```

