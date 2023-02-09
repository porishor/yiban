# 易班日报填写

ℹ️这是一个windows命令行程序

1. 运行任务

从`$APPDATA/yiban/user.csv`读取数据，执行日报提交
```shell
$ yiban run
```
2. 使用`client`的函数

```go
helper, err = client.CreateHelperWithPassword("110","aabb")
// 填写task_id，以提交单个任务
helper.Submit("xdqxdi1qnxio1xno")
```