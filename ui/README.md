# GoDad 前端（UI）

基于 Vue 3 + TypeScript + Vite + Pinia + Tailwind 的前端工程。

## 架构与约定（重要）

- 鉴权：使用 httpOnly Cookie（`access_token` + `refresh_token`）。前端请求默认携带凭据（`credentials: 'include'`），不再从 localStorage 注入 Authorization 头。
- 刷新：接口返回 401 时，HttpClient 会自动调用 `/auth/refresh-token` 刷新并重放原请求，页面无感。
- 统一响应壳：后端返回 `{ code, message, data }`，前端将 `code=0 或 200` 视为成功，其它视为业务错误并抛出 `AppError`。
- 统一分页形态：前端消费统一结构 `{ items, total, page, size, total_pages }`。通过 `normalizePageResponse()` 适配不同后端分页字段。
- 端点集中：所有端点统一在 `src/api/config.ts` 的 `API_CONFIG.ENDPOINTS`，API 模块调用配置常量，不在页面内发起原始 http 请求。
- 通知实时化：优先使用 SSE `/notifications/stream`，失败时自动降级到轮询；页面不可见时暂停轮询。

## 环境变量

- `VITE_API_BASE_URL` 后端基础地址（默认 `http://127.0.0.1:8888`）。开发时建议保留 Vite 代理 `/api` → 后端。

## 开发

```
pnpm i
pnpm dev
```

## 目录要点

- `src/api/`：API 模块，仅在这里与后端交互，使用 `API_CONFIG.ENDPOINTS` + `HttpClient`。
- `src/api/http.ts`：统一请求封装（携带凭据、401 刷新、统一错误抛出）。
- `src/api/pagination.ts`：分页适配器，输出统一分页结构。
- `src/utils/errorHandler.ts`：统一错误处理，配合 Toast 桥接在非组件环境显示提示。
- `src/components/Toast.vue` + `src/composables/useToast.ts`：唯一的 Toast 实现；`App.vue` 注入桥接供全局使用。

## 测试

```
pnpm test
```

- 单元测试覆盖：
  - `tests/unit/httpClient.test.ts`：401 → 刷新 → 重试流程；成功与业务错误解析
  - `tests/unit/pagination.test.ts`：常见分页结构映射

## 迁移提示

- 已移除对本地 token 注入的依赖；如需读取用户状态，请通过 Pinia `auth` store。
- 若新增接口：
  1) 在 `API_CONFIG.ENDPOINTS` 增加端点
  2) 在 `src/api/` 新增/修改模块函数（必要时用 `normalizePageResponse`）
  3) 页面只调用 API 模块，避免直接使用 `http` 或 `fetch`
