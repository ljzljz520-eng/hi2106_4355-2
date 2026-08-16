# 培训作业成绩系统

这是一个使用 Go 1.22.12 编写的纯后端命令行项目。数据默认保存在内存中，也可以通过菜单保存到本地 JSON 文件或从文件读取。

## 功能

- 录入学员编号、姓名、作业一、作业二和结课测试成绩
- 校验每项成绩为 0 到 100 的整数
- 查看全部成绩、按编号查询、按字段排序
- 删除和修改学员成绩
- 确定性保存和读取 JSON 数据

## 运行

```bash
go run ./cmd/traininggrades
```

使用仓库内固定数据启动：

```bash
go run ./cmd/traininggrades -fixture ./fixtures/students.json
```

指定菜单保存和读取使用的数据文件：

```bash
go run ./cmd/traininggrades -data ./grades.json
```

## 测试

```bash
CGO_ENABLED=0 go test -count=1 ./...
```

项目不需要数据库、网络服务或外部运行时数据。
