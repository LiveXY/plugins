# Plugins

本项目收集了作为插件使用的多个通用组件，用于增强应用程序的功能。

## 目录结构

- `dameng`: 达梦数据库适配
- `excel`: Excel 文件处理
- `mysql`: MySQL 数据库适配
- `redis`: Redis 缓存适配（支持单机与集群）
- `local-fs`: 本地文件系统操作
- `pgsql`: PostgreSQL 数据库适配
- ... 其他组件

## 开发环境

- **Go版本**: >= 1.25
- **依赖管理**: Go Modules

## 编译与测试

本项目使用 `Makefile` 进行管理。

```bash
# 更新所有依赖
make update-depend

# 编译所有插件 (.so 文件)
make build-plugins

# 安全检查
make govulncheck
make gosec
```
