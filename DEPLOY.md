# StarRocksQueris 部署指南

## 一、编译步骤

### 1. 环境要求
- Go 1.23.8 或更高版本
- MySQL 5.7+ 或 MariaDB（用于配置库）
- StarRocks 集群（被监控目标）

### 2. 服务器上编译

```bash
# 1. 进入项目目录
cd /path/to/StarRocksQueris

# 2. 设置 Go 环境变量
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct  # 国内服务器使用代理

# 3. 下载依赖
go mod tidy

# 4. 编译
go build -o StarRocksQueris .

# 5. 检查编译结果
ls -lh StarRocksQueris
```

### 3. 如果编译遇到问题

```bash
# 清理缓存重新编译
go clean -cache
go mod tidy
go build -o StarRocksQueris .
```

## 二、配置文件说明

### 1. 主配置文件

复制 `.StarRocksQueris.yaml` 为可执行文件同目录下的隐藏文件：

```bash
cp .StarRocksQueris.yaml /path/to/executable/.StarRocksQueris.yaml
```

配置文件内容示例：

```yaml
configdb:
  Host: 127.0.0.1      # MySQL配置库地址
  Port: 3306           # MySQL端口
  User: root           # MySQL用户名
  Pass: your_password  # MySQL密码
  Db: your_database    # MySQL数据库名（必填，用于初始化连接）
  Schema:
    App: starrocks_information_slowconfig      # 慢查询配置表（只需表名）
    Connect: starrocks_information_connections # 集群连接表（只需表名）
    Robot: starrocks_information_larkrobot     # 飞书机器人表（只需表名）
    WeComRobot: starrocks_information_wecomrobot # 企业微信机器人表（只需表名）

mode:
  cronsyntax: '*/1 * * * *'  # 定时任务表达式，每分钟执行

logger:
  LogPath: /var/log/StarRocksQueris  # 日志目录
  LogLevel: 0
  MaxSize: 0
  MaxBackups: 0
  MaxAge: 0
  Compress: false
  JsonFormat: false
  ShowLine: true
  LogInConsole: true
```

### 2. 规则配置文件

确保 `config/starrocks_rule.yaml` 存在：

```bash
mkdir -p config
cp config/starrocks_rule.yaml /path/to/executable/config/starrocks_rule.yaml
```

规则配置文件内容：

```yaml
# 慢查询规则配置
slow_query:
  alert_timeout_seconds: 600      # 慢查询告警超时时间（秒）
  kill_timeout_seconds: 1500      # 慢查询查杀时间（秒）
  concurrency_limit: 80           # 并发度限制
  full_scan_rows_limit: 200000000 # 全表扫描行数限制（2亿）
  catalog_scan_rows_limit: 100000000
  be_memory_limit_gb: 200
  scan_rows_limit: 10000000000    # 100亿
  scan_bytes_limit_tb: 5
  resource_group_cpu_core: 10
  resource_group_mem_gb: 50
  resource_group_concurrency: 3
  feishu_proxy: ""                # 飞书代理地址
  email_suffix: "@company.com"    # 企业邮箱后缀
  data_registration_port: 8030

# 集群连接规则
cluster_connection:
  default_fe_port: 9030
  license_expire_remind_days: 30
  license_check_switch: 0

# 短查询规则
short_query:
  default_init_push: 0
  default_status: 0
  default_core: 4
  default_memory_gb: 8

# 飞书机器人配置
lark_robot:
  default_status: 0

# 企业微信机器人配置
wecom_robot:
  default_status: 0               # 0=关闭，1=开启
  webhook_key: "your-webhook-key" # 企业微信Webhook Key
  msg_type: "markdown"            # 消息类型：text/markdown
  mention_all: false              # 是否@所有人
  mentioned_mobile_list: []       # @指定手机号列表
```

## 三、数据库初始化

### 1. 创建配置库

在 MySQL 中执行：

```sql
-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS srqueris CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE srqueris;
```

### 2. 创建表结构

执行 `exec.sql` 中的 SQL 语句：

source /xx/xx/exec.sql


## 四、启动服务

### 1. 直接启动

```bash
./StarRocksQueris
```

### 2. 后台启动

```bash
nohup ./StarRocksQueris > /var/log/StarRocksQueris/app.log 2>&1 &
```

### 3. 使用 Systemd 管理（推荐）

创建服务文件 `/etc/systemd/system/starrocks-queris.service`：

```ini
[Unit]
Description=StarRocks Queris Service
After=network.target

[Service]
Type=simple
User=starrocks
Group=starrocks
WorkingDirectory=/opt/starrocks-queris
ExecStart=/opt/starrocks-queris/StarRocksQueris
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
# 重新加载 systemd
systemctl daemon-reload

# 启动服务
systemctl start starrocks-queris

# 设置开机自启
systemctl enable starrocks-queris

# 查看状态
systemctl status starrocks-queris

# 查看日志
journalctl -u starrocks-queris -f
```

## 五、企业微信配置

### 配置方式：企业微信应用（CorpID + AgentId + Secret）

本项目使用**企业微信应用**方式发送告警消息，需要在企业微信管理后台创建应用。

### 1. 获取企业微信应用凭证

1. 登录企业微信管理后台：https://work.weixin.qq.com/wework_admin
2. 进入"应用管理" -> "应用" -> "创建应用"
3. 填写应用信息，创建后获取以下凭证：
   - **企业ID (CorpID)**：在"我的企业"页面获取
   - **AgentId**：在应用详情页获取
   - **Secret**：在应用详情页获取

### 2. 配置企业微信告警

编辑 `config/starrocks_rule.yaml`：

```yaml
wecom_app:
  default_status: 1                          # 0=关闭，1=开启
  corp_id: "wx7d3e81e155049ca3"             # 企业ID
  agent_id: "1000027"                        # 应用ID
  secret: "DvpGuNYmQEIxn-DBzY9Q7_pzEATOE8b3na3WOuijZ7s"  # 应用密钥
  msg_type: "markdown"                       # 消息类型：text/markdown
  mention_all: false                         # 是否@所有人
  mentioned_user_list:                       # @指定用户（企业微信中的用户账号）
    - "zhangsan"
    - "lisi"
```

### 3. 验证配置

重启服务后，可以在 StarRocks 中执行测试查询：

```sql
-- 执行一个耗时查询测试告警
SELECT SLEEP(700);
```

如果配置正确，将在 10 分钟后收到企业微信应用消息。

### 注意事项

1. **应用可见范围**：确保应用对需要接收告警的成员可见
2. **IP白名单**：如果企业微信设置了IP白名单，需要将服务器IP加入白名单
3. **消息频率**：企业微信应用消息有频率限制，默认20条/秒

---

### 旧版配置（Webhook机器人方式，已废弃）

如需使用 Webhook 机器人方式，请参考历史版本文档。

```sql
UPDATE `starrocks_information_wecomrobot`
SET `webhook_key` = 'xxxxxxxxxxxxxxxx',
    `msg_type` = 'markdown',
    `status` = 1
WHERE `type` = 'global';
```

重启服务生效：

```bash
systemctl restart starrocks-queris
```

## 六、验证部署

### 1. 检查服务状态

```bash
# 查看进程
ps aux | grep StarRocksQueris

```

### 2. 查看日志

```bash
# 应用日志
tail -f /var/log/StarRocksQueris/app.log

# 慢查询日志
tail -f /var/log/StarRocksQueris/slowquery.log
```

### 3. 测试告警

可以在 StarRocks 中执行一个慢查询测试告警功能：

```sql
-- 执行一个耗时查询
SELECT SLEEP(700);
```

如果配置正确，应该在 10 分钟后收到企业微信告警。

## 七、常见问题

### 1. 编译报错：package not found

```bash
# 确保在正确的目录下编译
cd /path/to/StarRocksQueris

# 清理 Go 缓存
go clean -cache
go mod tidy
go build -o StarRocksQueris .
```

### 2. 启动报错：配置文件未找到

确保配置文件路径正确：
- 主配置文件：`.StarRocksQueris.yaml`（与可执行文件同目录）
- 规则配置文件：`config/starrocks_rule.yaml`

### 3. 数据库连接失败

检查：
- MySQL 服务是否启动
- 用户名密码是否正确
- 防火墙是否放行端口

### 4. 企业微信告警未收到

检查：
- `wecom_robot.default_status` 是否为 1
- `webhook_key` 是否正确
- 查看应用日志是否有发送记录

## 八、目录结构

部署后的目录结构：

```
/opt/starrocks-queris/
├── StarRocksQueris              # 可执行文件
├── .StarRocksQueris.yaml        # 主配置文件
├── config/
│   └── starrocks_rule.yaml      # 规则配置文件
└── logs/                        # 日志目录（自动创建）
    ├── app.log
    └── slowquery.log
```
