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

#### 最佳实践

1. 启动 web 服务

2. 使用 `cron` 等定时执行工具执行 `./bin/ucasnj-smi check` (注意执行时 `$PWD` 为 `history.db` 所在文件夹，否则历史记录不会被记录)

### 4. Docker 部署 ucasnj-smi-server
```bash
# 构建镜像
docker build -t ucasnj-smi .

# 运行容器 (挂载配置文件和数据库)
docker run -d \
  --name ucasnj-smi-server \
  -p 8080:8080 \
  crpi-z6352oddczzzx18w.cn-hangzhou.personal.cr.aliyuncs.com/jiny3/ucasnj-smi:latest
```

## 目录结构

```bash
├── cmd
│   ├── check.go
│   ├── root.go
│   ├── server.go
│   └── utils.go
├── dockerfile
├── go.mod
├── go.sum
├── history.db  # 历史记录数据库
├── library
│   ├── check.go
│   ├── datahub.go
│   ├── encrypt_aes.go
│   └── pushplus.go
├── LICENSE
├── main.go
├── ops.log     # json 格式 log
├── README.md
├── service
│   ├── history.go
│   └── static.go
├── static
│   ├── echarts.js
│   ├── index.html
│   └── my.js
└── users.toml  # 用户配置
```

## License

MIT
