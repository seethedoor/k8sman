# k8sman - Kubernetes 集群管理仪表盘

## 项目介绍

k8sman 是一个功能强大的 Kubernetes 集群管理仪表盘，提供直观的 Web 界面用于监控和管理 Kubernetes 集群。它支持多种集群操作，包括资源管理、日志查看、终端访问等功能，帮助用户更高效地管理 Kubernetes 集群。

## 功能特性

### 核心功能
- **集群概览**：查看集群整体状态和资源使用情况
- **资源管理**：管理和监控 Pod、Deployment、Service、ConfigMap 等 Kubernetes 资源
- **日志查看**：实时查看容器日志，支持多容器日志聚合
- **终端访问**：直接在浏览器中访问容器终端
- **事件监控**：查看集群中的事件和告警
- **资源详情**：查看资源的详细信息和配置

### 技术特点
- **响应式设计**：适配不同屏幕尺寸，提供良好的用户体验
- **实时数据**：通过 WebSocket 实现实时数据更新
- **安全认证**：支持 Kubernetes 集群认证
- **模块化架构**：前后端分离，易于扩展和维护

## 技术栈

### 后端
- **语言**：Go 1.22.0
- **框架**：Gin
- **Kubernetes 客户端**：client-go
- **构建工具**：Go 标准工具链

### 前端
- **框架**：Vue 3 + TypeScript
- **状态管理**：Pinia
- **路由**：Vue Router
- **UI 组件**：Element Plus
- **图表**：ECharts
- **终端模拟器**：xterm.js
- **构建工具**：Vite
- **样式**：Tailwind CSS

## 环境要求

### 后端
- Go 1.22.0 或更高版本
- Kubernetes 集群访问配置（kubeconfig）

### 前端
- Node.js 18.0 或更高版本
- npm 或 yarn 包管理器

## 安装和启动

### 后端启动

#### 方法 1：使用编译好的二进制文件

```bash
# 进入 backend 目录
cd backend

# 运行编译好的服务器
./k8s-dashboard-server
```

#### 方法 2：从源码构建

```bash
# 进入 backend 目录
cd backend

# 构建服务器
go build -o k8s-dashboard-server ./cmd/server

# 运行服务器
./k8s-dashboard-server
```

### 前端启动

```bash
# 进入 frontend 目录
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build
```

## 项目结构

```
├── backend/                # 后端代码
│   ├── cmd/                # 命令行入口
│   │   └── server/         # 服务器入口
│   ├── internal/           # 内部包
│   │   ├── handler/        # HTTP 处理器
│   │   ├── kubernetes/     # Kubernetes 客户端
│   │   └── middleware/     # 中间件
│   ├── pkg/                # 公共包
│   │   ├── config/         # 配置管理
│   │   └── response/       # 响应处理
│   ├── go.mod              # Go 模块文件
│   ├── go.sum              # 依赖校验和
│   └── k8s-dashboard-server # 编译后的二进制文件
├── frontend/               # 前端代码
│   ├── public/             # 静态资源
│   ├── src/                # 源代码
│   │   ├── api/            # API 调用
│   │   ├── components/     # 组件
│   │   ├── composables/    # 组合式函数
│   │   ├── pages/          # 页面
│   │   ├── router/         # 路由
│   │   ├── stores/         # 状态管理
│   │   ├── views/          # 视图
│   │   └── types/          # 类型定义
│   ├── package.json        # npm 配置文件
│   └── vite.config.ts      # Vite 配置
├── README.md               # 项目说明
└── merge_summary.md        # 合并摘要
```

## 如何使用

1. **启动后端服务器**：按照上述方法启动后端服务器
2. **启动前端开发服务器**：按照上述方法启动前端开发服务器
3. **访问仪表盘**：在浏览器中访问前端服务器地址（默认为 http://localhost:5173）
4. **登录认证**：使用 Kubernetes 集群的认证信息登录
5. **管理集群**：通过仪表盘界面管理和监控 Kubernetes 集群

## 配置说明

### 后端配置

后端服务器默认使用 `~/.kube/config` 作为 Kubernetes 集群配置文件。如果需要指定其他配置文件，可以通过环境变量设置：

```bash
# 指定 kubeconfig 文件路径
export KUBECONFIG=/path/to/kubeconfig

# 启动服务器
./k8s-dashboard-server
```

### 前端配置

前端 API 地址配置在 `frontend/src/api/index.ts` 文件中，可以根据实际部署情况修改：

```typescript
// frontend/src/api/index.ts
const apiClient = axios.create({
  baseURL: 'http://localhost:8000/api', // 修改为实际的后端 API 地址
  timeout: 10000,
});
```

## 贡献指南

1. **Fork 项目**：在 GitHub 上 fork 项目到自己的仓库
2. **克隆仓库**：克隆 fork 后的仓库到本地
3. **创建分支**：创建一个新的分支用于开发
4. **修改代码**：根据需要修改代码
5. **提交更改**：提交更改并推送到远程仓库
6. **创建 PR**：在 GitHub 上创建 Pull Request

## 许可证

本项目采用 MIT 许可证，详情请查看 [LICENSE](LICENSE) 文件。

## 联系方式

如有问题或建议，欢迎通过以下方式联系：

- GitHub Issues：在项目仓库中创建 Issue
- 邮件：请在 GitHub 个人资料中查看联系邮箱

---

**k8sman** - 让 Kubernetes 集群管理更简单！