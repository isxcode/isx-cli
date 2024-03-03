## isx-cli

##### 使用手册

```bash
isx -h
```

##### 在线安装

```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/isxcode/isx-cli/main/install.sh)"
```

##### brew安装

```bash
brew tap isxcode/isxcode && brew install isx
```

##### 本地安装

```bash
# 编译代码
docker run --rm \
  -v "/Users/ispong/isxcode/isx-cli":/usr/src/myapp \
  -w /usr/src/myapp \
  -e GOOS=darwin \
  -e GOARCH=arm64 \
  -e CGO_ENABLED=0 \
  golang:1.21 go build -v -o ./target/isx
  
# 安装命令
sudo mv /Users/ispong/isxcode/isx-cli/target/isx /usr/local/bin/isx
```

##### 开发流程

```bash
# 切出分支
isx checkout <issue_number>
# 提交代码
git commit -m "your commit message"
# 格式化代码
# 并推送到origin仓库
isx format
# 提交pr
isx pr <issue_number>
# 如果提交的pr，无法成功rebase合并
isx pull && isx format
```