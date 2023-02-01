## 主要功能

这是一个简单的cli工具，可以从数据库中导出表结构，包括字段的名称、类型、长度、是否允许为空、默认值、注释等。

## 安装

### 从release下载

从[release](https://github.com/Bit0r/schema2file/releases/latest)下载已经编译好的二进制文件，解压后即可使用。

### 从源码编译

```bash
go install github.com/Bit0r/schema2file@latest
```

## 使用

### 示例

```bash
schema2file -u <user> -p <password> -B <database> -o <file.md/xlsx>
```

### 命令行参数

```
-B string
	database name
-P int
	port number (default 3306)
-h string
	host name (default "localhost")
-o string
	output file name (default "schema.md")
-p string
	password
-u string
	user name (default "root")
```
