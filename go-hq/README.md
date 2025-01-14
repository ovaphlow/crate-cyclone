# CRATE-CYCLONE HQ

## 脚本

同步依赖

```shell
make tidy
```

## 构建

### 构建 Linux 版本

```shell
make build-linux
```

### 交叉编译到 Windows

```shell
make build-windows
```

### 同时编译 Linux 和 Windows

```shell
make build
```

## 清理构建文件

```shell
make clean
```

## 测试

在工程根目录执行 `utility` 目录的测试

```shell
make test-utility
```

## 运行程序

```shell
make run
```

## 设置库路径

请确保 `libcrate_shared.so` 文件存放在 `./lib` 目录下。

## 在 Windows 系统上使用 Makefile

在 Windows 系统上使用 Makefile，您需要安装一个兼容的工具，例如 Git Bash 或 Cygwin。安装完成后，您可以在这些终端中运行 `make` 命令。