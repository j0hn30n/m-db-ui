# MongoDB管理工具

一个使用Go语言开发的MongoDB可视化管理工具，提供简洁的Web界面来管理MongoDB数据库。

## 功能特性

- 📊 **数据库管理**: 创建、查看、删除数据库
- 🗂️ **集合管理**: 创建、查看、删除集合
- 📄 **文档管理**: 增删改查文档数据
- 🔍 **高级查询**: 支持复杂的MongoDB查询
- 📈 **统计信息**: 查看服务器和数据库统计
- 🎨 **现代UI**: 基于Bootstrap的响应式界面
- 📱 **移动友好**: 支持移动设备访问

## 快速开始

### 环境要求

- Go 1.19+
- MongoDB 3.6+

### 安装运行

1. 克隆项目
```bash
git clone [https://github.com/j0hn30n/m-db-ui.git](https://github.com/j0hn30n/m-db-ui.git)
cd m-db-ui
```

2. 安装依赖
```bash
go mod download
```

3. 配置环境变量
```bash
# 复制配置文件模板
cp .env.example .env

# 编辑配置文件
vim .env
```

4. 启动服务
```bash
go run main.go
```

5. 访问应用
打开浏览器访问 `http://localhost:8082`


## API文档

### 数据库管理

- `GET /api/v1/databases` - 获取所有数据库
- `GET /api/v1/databases/{name}` - 获取数据库信息
- `DELETE /api/v1/databases/{name}` - 删除数据库

### 集合管理

- `GET /api/v1/databases/{db}/collections` - 获取所有集合
- `POST /api/v1/databases/{db}/collections` - 创建集合
- `DELETE /api/v1/databases/{db}/collections/{collection}` - 删除集合

### 文档管理

- `GET /api/v1/databases/{db}/collections/{collection}/documents` - 获取文档列表
- `POST /api/v1/databases/{db}/collections/{collection}/documents` - 创建文档
- `PUT /api/v1/databases/{db}/collections/{collection}/documents/{id}` - 更新文档
- `DELETE /api/v1/databases/{db}/collections/{collection}/documents/{id}` - 删除文档
- `POST /api/v1/databases/{db}/collections/{collection}/query` - 查询文档

### 统计信息

- `GET /api/v1/stats` - 获取服务器统计信息

## 项目结构

```
m-db-ui/
├── main.go                 # 主程序入口
├── config.yaml            # 配置文件
├── .env.example           # 环境变量模板
├── internal/              # 内部包
│   ├── config/           # 配置管理
│   ├── database/         # 数据库操作
│   ├── handlers/         # HTTP处理器
│   └── services/         # 业务逻辑
├── web/                  # Web资源
│   ├── static/          # 静态文件
│   │   ├── css/
│   │   └── js/
│   └── templates/       # HTML模板
└── docs/                # 文档
```

## 开发指南

### 本地开发

1. 启动MongoDB服务
```bash
# 使用Docker启动MongoDB
docker run -d -p 27017:27017 --name mongodb mongo:latest
```

2. 运行开发服务器
```bash
go run main.go
```

3. 构建生产版本
```bash
go build -o m-db-ui main.go
```

### 添加新功能

1. 在 `internal/database/` 中添加数据库操作
2. 在 `internal/handlers/` 中添加HTTP处理器
3. 在 `web/templates/` 中添加前端页面
4. 在 `web/static/` 中添加CSS/JS文件

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。
