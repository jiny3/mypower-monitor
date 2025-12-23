# ucasnj-smi

UCASNJ Dormitory Power Monitor

## 功能特性

- 🔋 **电量查询**: 支持多用户批量查询宿舍剩余电量
- 🔔 **消息推送**: 支持 PushPlus 微信推送
- 📊 **历史记录**: 自动记录电量数据，支持历史趋势查看
- 🌐 **Web服务**: 提供可视化 Web 界面展示电量变化
- 🐳 **Docker支持**: 支持 Docker 部署，方便在服务器或 NAS 上运行

## 快速开始

### 1. 编译安装

```bash
# 克隆项目
git clone https://github.com/jiny3/mypower-monitor.git
cd mypower-monitor

# 编译
go build -o bin/ucasnj-smi .
```

### 2. 配置文件

在项目根目录下创建 `users.toml` 文件，填入用户信息：

```toml
[[users]]
account = "2023xxxxxxx"      # 校园网账号
password = "mypassword"      # 校园网密码
room_id = "b905"             # 宿舍号 (如: b905)
token = "xxxxxxxxxxxxxxx"    # (可选) PushPlus Token
to = "xxxxxxxxxxxxxxx"       # (可选) PushPlus 好友ID
```

### 3. 使用说明

#### 命令行工具

```bash
# 查看帮助
./bin/ucasnj-smi --help

# 执行一次电量检查 (读取配置文件)
./bin/ucasnj-smi check

# 临时检查指定用户 (忽略配置文件)
./bin/ucasnj-smi check -a "2023xxxx" -p "password" -r "b905"

# 启动 Web 服务 (默认端口 8080)
./bin/ucasnj-smi server
# 指定端口启动
./bin/ucasnj-smi server 9090
```

### 4. Docker 部署

```bash
# 构建镜像
docker build -t ucasnj-smi .

# 运行容器 (挂载配置文件和数据库)
docker run -d \
  --name ucasnj-smi-server \
  -p 8080:8080 \
  jiny14/ucasnj-smi
```

## 目录结构

- `cmd/`: 命令行入口及逻辑
- `service/`: Web 服务逻辑
- `library/`: 核心功能库 (爬虫、加密、推送)
- `static/`: 前端静态资源
- `users.toml`: 用户配置文件

## License

MIT
