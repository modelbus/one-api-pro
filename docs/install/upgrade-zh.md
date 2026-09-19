---
title: 版本升级
description: 升级到新版本前要做什么，升级过程中要注意什么。
category: install
order: 6
---

# 版本升级

> 升级前 / 升级中 / 升级后的标准流程。

## 升级前

1. **备份数据库**：详见 [备份与恢复](./backup-restore)
2. **阅读 CHANGELOG**：看升级版本对应 [CHANGELOG](https://github.com/modelbus/one-api-pro/blob/main/CHANGELOG) 的「升级注意事项」段落
3. **检查破坏性变更**：不兼容的字段、不兼容的接口、不兼容的配置项

## 升级方式

### Docker / docker-compose

```bash
docker compose pull one-api-pro
docker compose up -d
```

只重启 one-api-pro 容器。数据卷（`./data` / `./mysql-data`）不变。

### 源码部署

```bash
git pull
cd web/default-pro && pnpm install && pnpm build && cd ../..
go build -o one-api-pro .
# 停旧进程，启动新进程
systemctl stop one-api-pro
systemctl start one-api-pro
```

### 二进制部署

直接替换二进制文件，重启服务。数据文件不变。

## 升级过程中

容器启动时 `AutoMigrate` 会**自动建表 / 加列**。日志会显示迁移过程：

```
[migrate] adding column foo.bar
[migrate] creating index idx_xxx
```

如果失败，看日志中的具体错误（通常是字段类型冲突）。

## 升级后

1. **看一遍 [CHANGELOG](https://github.com/modelbus/one-api-pro/blob/main/CHANGELOG)** 里的「升级注意事项」
2. **登录后台 → 运营仪表盘**，对比升级前后的请求成功率
3. **抽样几条调用日志**，确认调用行为符合预期
4. **新功能**：有些版本需要在后台「系统设置」里手动配置（CHANGELOG 会指出）

## 回滚

若升级后出现严重问题：

```bash
# Docker
docker compose down
# 把备份的 mysql-data / data 恢复回去
mv ./mysql-data ./mysql-data.broken
cp -r ./backup-YYYYMMDD/mysql-data ./
docker compose up -d
```

升级过程中没有动过数据库 schema 的话，可以直接降级回滚。

## 相关文档

- [备份与恢复](./backup-restore)
- [CHANGELOG](../changelog/index)