# GO采集性能指标

## build注意事项
1.go build和平台会关系，需要现在环境变量中进行设置

2.设置`GOOS`环境变量

> unset GOOS

> export GOOS=linux

3.执行build操作

> go build

## docker build
1.进入项目目录
cd plum_exporter

2.执行build操作
```shell
docker run --rm -ti -v $(pwd):/app quay.io/prometheus/golang-builder -i "github.com/seekplum/plum_exporter" -p "linux/amd64"
```
