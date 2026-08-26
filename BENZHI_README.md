# gear-contact：Go 渐开线直齿啮合 Web 服务（几何内核 + Hertz/Lewis + 前端控制台）

给定模数、齿数与压力角，核算中心距、重合度与接触应力；提供 `/api/mesh`、`/api/hertz` 与嵌入网页。

## 构建 / 运行 / 测试

```text
go build ./...
./gear-contact -http :8080
curl -s http://127.0.0.1:8080/api/example
go run . mesh example/spur-20.json
go test ./...
```

## 评测镜像

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -d -P --name gear-contact-b14 <image-name>:latest
curl -s http://127.0.0.1:$(docker port gear-contact-b14 8080 | cut -d: -f2)/api/example
docker rm -f gear-contact-b14
```
