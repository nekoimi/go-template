# Go Template

基于 Go 的后端项目模板，集成常用组件，开箱即用。

## 技术栈

| 组件 | 说明 |
|---|---|
| [Gin](https://github.com/gin-gonic/gin) | HTTP 框架 |
| [GORM](https://gorm.io/) | ORM |
| [PostgreSQL](https://www.postgresql.org/) | 数据库 |
| [Zap](https://github.com/uber-go/zap) | 结构化日志 |
| [Viper](https://github.com/spf13/viper) | 配置管理 |
| [golang-jwt](https://github.com/golang-jwt/jwt) | JWT 认证 |
| [gorilla/websocket](https://github.com/gorilla/websocket) | WebSocket |
| [robfig/cron](https://github.com/robfig/cron) | 定时任务 |
| [snowflake](https://github.com/bwmarrin/snowflake) | 分布式 ID |
| [gin-swagger](https://github.com/swaggo/gin-swagger) | API 文档 |
| [MinIO](https://min.io/) | 对象存储 (可选) |

## 项目结构

```
├── cmd/
│   ├── server/main.go          # HTTP 服务入口
│   └── scheduler/main.go       # 独立定时任务入口
├── configs/                    # 配置文件
├── deployments/                # Dockerfile & docker-compose
├── docs/                       # Swagger 文档
├── internal/
│   ├── app/                    # 应用初始化
│   ├── config/                 # 配置结构与加载
│   ├── dto/                    # 请求/响应 DTO
│   ├── handler/
│   │   ├── middleware/         # 中间件 (CORS, JWT, 限流, 日志, Recovery, RequestID)
│   │   └── v1/                 # API Handler
│   ├── infrastructure/         # 基础设施 (数据库, 日志)
│   ├── model/                  # GORM 数据模型
│   ├── pkg/
│   │   ├── errcode/            # 业务错误码
│   │   ├── response/           # 统一响应格式
│   │   ├── snowflake/          # ID 生成器
│   │   └── timeutil/           # 时区与时间类型
│   ├── repository/             # 数据访问层
│   ├── scheduler/              # 定时任务
│   ├── service/                # 业务逻辑层
│   ├── storage/                # 文件存储 (local / minio)
│   └── websocket/              # WebSocket (Hub 模式)
├── migrations/                 # 数据库迁移脚本
├── docker-compose.yml          # 开发基础设施 (PG + MinIO)
├── Makefile
└── go.mod
```

## 快速开始

### 1. 启动开发基础设施

```bash
make dev-up    # 启动 PostgreSQL + MinIO
```

### 2. 运行数据库迁移

```bash
make migrate-up
```

### 3. 启动服务

```bash
make run       # HTTP 服务 http://localhost:8080
```

定时任务可独立运行：

```bash
make run-scheduler
```

## 配置

配置文件位于 `configs/`，通过 `--config` 参数指定。支持环境变量覆盖：

| 环境变量 | 对应配置 |
|---|---|
| `DATABASE_HOST` | database.host |
| `DATABASE_PORT` | database.port |
| `DATABASE_USER` | database.user |
| `DATABASE_PASSWORD` | database.password |
| `DATABASE_NAME` | database.dbname |
| `JWT_SECRET` | jwt.secret |
| `TZ` | server.timezone |
| `SNOWFLAKE_NODE_ID` | snowflake.node_id |

完整配置项见 `configs/config.dev.yaml`。

## API

启动后访问 Swagger UI：`http://localhost:8080/swagger/index.html`

| 方法 | 路径 | 认证 | 说明 |
|---|---|---|---|
| GET | `/health` | - | 存活检查 |
| GET | `/ready` | - | 就绪检查 (DB ping) |
| POST | `/v1/auth/register` | - | 用户注册 |
| POST | `/v1/auth/login` | - | 用户登录 |
| GET | `/v1/users/profile` | JWT | 获取当前用户信息 |
| POST | `/v1/upload/single` | JWT | 上传单个文件 |
| POST | `/v1/upload/multiple` | JWT | 上传多个文件 |
| GET | `/ws/v1/chat?token=` | JWT | WebSocket 连接 |

### 统一响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {},
  "error": null
}
```

错误时 `code` 为业务错误码（如 40100=未授权，40401=用户不存在），`error` 包含错误详情。

## 构建与部署

```bash
make build         # 编译到 bin/
make swagger       # 重新生成 Swagger 文档
make test          # 运行测试

make docker-build  # 构建 Docker 镜像
make docker-up     # 启动完整部署 (app + scheduler + PG + MinIO)
make docker-down   # 停止
```

## License

[MIT](LICENSE)
