# GoDad 后端鉴权与API约定

## 鉴权策略

- 授权介质：`access_token`（短期）、`refresh_token`（长期），均通过 httpOnly Cookie 下发。
- 刷新：`POST /api/auth/refresh-token` 从 `refresh_token` Cookie 刷新 `access_token`。
- 兼容：中间件同时支持 `Authorization: Bearer <token>` 与 `access_token` Cookie，便于迁移与调试。
- CORS：`Access-Control-Allow-Credentials: true`，并将 `Origin` 精确到前端地址（开发默认 `http://127.0.0.1:3333`）。

## 统一响应壳

```
{ code: number, message: string, data: any }
```

- 成功：`code = 200`
- 业务错误：使用 4xx/5xx 业务码作为 `code`，HTTP 状态可以保持 200（或语义对应状态码）

## 分页约定

后端可返回以下任一结构，前端将做统一适配：

- 扁平：`{ data: T[], total, page, size }`
- 标准：`{ data: { items: T[], total, page, size, total_pages } }`
- 通知：`{ data: { notifications: T[], pagination: { total, current_page, per_page, total_pages } } }`
- 话题：`{ data: { items: T[], total, page, page_size, total_page } }`

尽量向“标准”靠拢；前端已提供适配器以缩短迁移周期。

## 通知实时化

- SSE：`GET /api/notifications/stream`，按用户推送未读统计与通知事件；支持 `Last-Event-ID`。
- 降级：若 SSE 不可用，前端降级轮询 `/api/notifications/stats` 与 `/api/notifications`。

## 环境

- `SERVER_ENV=development` 时允许 DB 自动迁移与本地 CORS；生产需收紧 CORS 与 Cookie Secure/SameSite。

