# Signal ZXH

一个基于 Go + Gin + MySQL 的轻量级博客系统，支持文章管理和静态页面展示。

## 功能特性

- 📝 文章 CRUD 操作（创建、读取、更新、删除）
- 🎨 静态页面服务
- 🗄️ MySQL 数据持久化
- 🐳 Docker 容器化部署
- 🚀 RESTful API 设计

## 技术栈

- **后端**: Go 1.24.0 + Gin
- **数据库**: MySQL 9.7
- **容器**: Docker + Docker Compose
- **前端**: 原生 HTML + CSS

## 快速开始

### 前置要求

- Docker
- Docker Compose

### 本地开发

1. 克隆项目
```bash
git clone https://github.com/Signal-zxh/signal-zxh.git
cd signal-zxh
```

2. 配置环境变量
```bash
cp .env.example .env
# 编辑 .env 文件，配置数据库连接信息
```

3. 安装依赖
```bash
go mod download
```

4. 运行项目
```bash
go run main.go
```

  服务将在 http://47.96.119.143:8080 启动

### Docker 部署

1. 配置环境变量
```bash
# 创建 .env 文件
echo "DBPASS=your_password" > .env
echo "MYSQL_ROOT_PASSWORD=your_root_password" >> .env
```

2. 启动服务
```bash
docker-compose up -d
```

3. 查看日志
```bash
docker-compose logs -f signal-zxh
```

## 项目结构

```
signal-zxh/
├── db/              # 数据库初始化
│   └── mysql.go
├── handler/         # HTTP 处理器
│   └── post.go
├── model/           # 数据模型
│   └── post.go
├── static/          # 静态文件
│   ├── index.html
│   └── style.css
├── mysql-conf/      # MySQL 配置
│   └── my.cnf
├── main.go          # 主程序入口
├── Dockerfile       # Docker 镜像构建
└── docker-compose.yml # Docker Compose 配置
```

## API 文档

### 获取所有文章
```http
GET /posts
```

### 获取单篇文章
```http
GET /posts/:id
```

### 创建文章
```http
POST /posts
Content-Type: application/json

{
  "title": "文章标题"
}
```

### 更新文章
```http
PUT /posts/:id
Content-Type: application/json

{
  "title": "更新后的标题"
}
```

### 删除文章
```http
DELETE /posts/:id
```

### 静态页面
```http
GET /           # 首页
GET /static/*   # 静态资源
```

## 配置说明

### 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| DBHOST | 数据库主机 | mysql |
| DBPORT | 数据库端口 | 3306 |
| DBUSER | 数据库用户 | root |
| DBPASS | 数据库密码 | - |
| DBNAME | 数据库名称 | blog |

### MySQL 配置

MySQL 配置文件位于 `mysql-conf/my.cnf`，包含：
- 字符集设置：utf8mb4
- InnoDB 缓冲池大小：256M
- 其他性能优化参数

## 部署说明

### 生产环境建议

1. 修改 MySQL root 密码
2. 配置数据库备份策略
3. 设置资源限制（已在 docker-compose.yml 中配置）
4. 配置 HTTPS
5. 设置日志轮转

### 资源限制

- signal-zxh: 200MB 内存
- MySQL: 300MB 内存

## 开发指南

### 添加新功能

1. 在 `model/` 中定义数据模型
2. 在 `handler/` 中实现业务逻辑
3. 在 `main.go` 中注册路由

### 数据库迁移

手动执行 SQL 或使用迁移工具：

```sql
CREATE TABLE IF NOT EXISTS posts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

## 许可证

MIT License

## 作者

Signal ZXH