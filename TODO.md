# Signal ZXH Blog — TODO

项目待办清单，按优先级排列。完成一项后将 `- [ ]` 改为 `- [x]`。

---

## P0 — 分类/标签功能收尾（推荐优先）

> 最新功能刚做到一半，收尾后前后端闭环、演示效果最好。

### 前端

- [x] 首页增加分类筛选（侧边栏或顶栏 tab，调用 `GET /categories/:id/posts`）
- [x] 首页支持按标签筛选（后端接口 `/tags/:id/posts` + 前端集成）
- [x] 管理后台增加「分类管理」页面（增删改，调用 `/api/categories`）
- [x] 管理后台增加「标签管理」页面（增删改，调用 `/api/tags`）
- [x] 后台编辑文章时支持修改分类与标签（目前创建已支持，编辑需确认）

### 测试

- [x] 补 `handler/category_test.go`
- [x] 补 `handler/tag_test.go`
- [x] 补 `service/category_test.go`
- [x] 补 `service/tag_test.go`
- [x] 补 `db/category_test.go` / `db/tag_test.go`（sqlmock）

### 文档

- [x] 更新 README：补充分类/标签 API 说明
- [x] 更新 Swagger 注释并重新生成 docs（如有变更）

---

## P1 — 技术债清理

> 减少维护成本，让新人 clone 后可直接上手。

### 路由与 API

- [x] 统一认证路由：目前 `/posts`（`router/auth.go`）与 `/api/posts`（`router/api.go`）并存，择一保留并迁移前端调用
- [x] 修正 Swagger 路径：`/probe` → `/api/tools/http`，`/agent` → `/api/tools/agent`
- [x] 确认公开/需认证接口划分一致（category、tag 的 GET/POST 权限）

### 配置与部署

- [ ] 添加 `.env.example`（README 已引用但仓库缺失）
- [ ] 统一环境变量命名（README 的 `DBHOST` vs 代码的 `DB_HOST` / `DB_USER`）
- [ ] 清理 `docker-compose.yml` 中未使用的 Postgres 服务，或补充使用说明
- [ ] 核对 docker-compose 环境变量映射是否与 `db/mysql.go` 一致

---

## P2 — Agent 工具完善

> 将 `agent/` 从占位实现变为可用功能。

- [ ] `GetPosts()` 接入 PostService，返回真实文章列表
- [ ] `GetPostByID()` 接入 PostService，移除硬编码 id `"26"`
- [ ] 支持按分类/标签查询（扩展 `RouteTool` 路由逻辑）
- [ ] 补 Agent 相关单元测试
- [ ] （可选）接入 LLM API，实现自然语言问答

---

## P3 — 生产环境加固

> 在线演示长期运行所需。

- [ ] 配置 HTTPS（Nginx / Caddy 反向代理）
- [ ] MySQL 定时备份策略
- [ ] Redis 持久化与备份说明
- [ ] JWT_SECRET、ADMIN_PASSWORD 强制从环境变量读取，禁止默认值上线
- [ ] 日志轮转配置
- [ ] （可选）Prometheus metrics 或健康检查扩展（依赖状态）

---

## P4 — 体验与内容功能（按需）

### 写作与展示

- [ ] Markdown 渲染（文章内容）
- [ ] 代码高亮
- [ ] 文章摘要自动生成（替代截断正文）

### 发现与互动

- [ ] 全文搜索（标题 + 内容，或 Elasticsearch / MySQL FULLTEXT）
- [ ] RSS 订阅
- [ ] Sitemap（SEO）
- [ ] 评论系统（自建或 Giscus / Utterances）

### 后台增强

- [ ] 文章草稿 / 发布状态
- [ ] 文章排序 / 置顶

---

## 参考路线

| 路线 | 顺序 |
|------|------|
| **A — 产品向** | P0 → Markdown → 搜索 → RSS |
| **B — 工程向** | P1 → P0 测试 → P3 → 数据库迁移工具 |
| **C — 特色向** | P2 → LLM 接入 → 智能问答/推荐 |

---

## 已完成（归档）

- [x] 文章 CRUD + 分页 + Redis 缓存
- [x] JWT 认证
- [x] 单元测试（post / middleware / cache / db / utils）
- [x] Swagger API 文档
- [x] GitHub Actions CI/CD
- [x] Docker 镜像自动构建
- [x] 分类/标签后端 API（基础版）
- [x] 首页、详情页、后台展示分类/标签
- [x] HTTP 探测工具、番茄钟
- [x] 管理后台内联编辑文章
