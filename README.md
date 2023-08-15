### isx-cli

##### 使用手册

```bash
isx -h
```

##### 在线安装

```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/isxcode/isx-cli/main/install.sh)"
```

##### 本地编译

```bash
docker run --rm \
  -v "/Users/ispong/isxcode/isx-cli":/usr/src/myapp \
  -w /usr/src/myapp \
  -e GOOS=darwin \
  -e GOARCH=arm64 \
  -e CGO_ENABLED=0 \
  golang:1.21 go build -v -o ./target/isx
```

##### 本地安装

````bash
sudo mv /Users/ispong/isxcode/isx-cli/target/isx /usr/local/bin/isx
````