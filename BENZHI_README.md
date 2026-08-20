# gear-contact

`gear-contact` 是一个渐开线直齿圆柱齿轮（外啮合）啮合核算命令行工具。输入一个 JSON 文件给出模数 m、两齿轮齿数 z1/z2、压力角 α（度）、齿顶高系数 ha*（可覆盖，默认 1）、顶隙系数 c*（默认 0.25）与可选变位系数 x1/x2；输出两轮的分度圆 d=m·z、基圆 db=d·cosα、齿顶圆 da、中心距、法向齿距、作用线长度、重合度 εα、顶隙与根切状态。

## 构建 / 运行 / 测试

```text
go build ./...     # 编译
go run .
go test ./...      # 测试
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
