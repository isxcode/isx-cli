<h1 align="center">
  至行云
</h1>

<h3 align="center">
  打造开发规范脚手架
</h3>

##### 使用手册

```bash
isx -h
```

##### 安装包下载

- [isx_linux_amd64](https://isxcode.oss-cn-shanghai.aliyuncs.com/zhixingyun/isx_linux_amd64)
- [isx_windows_amd64.exe](https://isxcode.oss-cn-shanghai.aliyuncs.com/zhixingyun/isx_windows_amd64.exe)
- [isx_darwin_amd64](https://isxcode.oss-cn-shanghai.aliyuncs.com/zhixingyun/isx_darwin_amd64)
- [isx_darwin_arm64](https://isxcode.oss-cn-shanghai.aliyuncs.com/zhixingyun/isx_darwin_arm64)

##### 在线安装

```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/isxcode/isx-cli/main/install.sh)"
```

##### brew安装

```bash
brew tap isxcode/isxcode && brew install isx
```

##### Go编译

```bash
cd /Users/ispong/isxcode/isx-cli
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ./target/isx_darwin_arm64 main.go
```

##### Docker编译

> mac(M系列)

```bash
docker run --rm \
  -v "/Users/ispong/isxcode/isx-cli":/usr/src/myapp \
  -w /usr/src/myapp \
  -e GOOS=darwin \
  -e GOARCH=arm64 \
  -e CGO_ENABLED=0 \
  golang:1.21 go build -v -o ./target/isx_darwin_arm64
```

> mac(Intel系列)

```bash
docker run --rm \
  -v "/Users/ispong/isxcode/isx-cli":/usr/src/myapp \
  -w /usr/src/myapp \
  -e GOOS=darwin \
  -e GOARCH=amd64 \
  -e CGO_ENABLED=0 \
  golang:1.21 go build -v -o ./target/isx_darwin_amd64
```

> linux

```bash
docker run --rm \
  -v "/home/ispong/isxcode/isx-cli":/usr/src/myapp \
  -w /usr/src/myapp \
  -e GOOS=linux \
  -e GOARCH=amd64 \
  -e CGO_ENABLED=0 \
  golang:1.21 go build -v -o ./target/isx_linux_amd64
```

##### 安装命令

```bash
sudo mv /Users/ispong/isxcode/isx-cli/target/isx_darwin_arm64 /usr/local/bin/isx
```

##### 开发流程

```bash
# 1.选择开发项目
isx choose
# 2.切出开发分支
isx checkout <issue_number>
# 3.提交代码到本地
git commit -m "your commit message"
# 4.格式化代码 (自动推送到origin仓库)
isx format
# 5.提交pr
isx pr <issue_number>
# 6.如果提交的pr，无法成功rebase合并，需要更新代码二次提交
isx pull && isx format
```