# StarRocksQueris 配置文件说明

## 简介

本项目现在支持通过 YAML 配置文件来管理慢查询规则，无需修改数据库中的配置。配置文件位于 `config/starrocks_rule.yaml`。

## 配置文件结构

### 1. 慢查询规则配置 (slow_query)

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| alert_timeout_seconds | 600 | 慢查询告警超时时间（秒） |
| kill_timeout_seconds | 1500 | 慢查询查杀时间（秒） |
| concurrency_limit | 80 | 慢查询并发度限制 |
| full_scan_rows_limit | 200000000 | 全表扫描行数限制（2亿） |
| catalog_scan_rows_limit | 100000000 | Catalog扫描行数限制（1亿） |
| be_memory_limit_gb | 200 | BE内存限制（GB） |
| scan_rows_limit | 10000000000 | 扫描行数限制（100亿） |
| scan_bytes_limit_tb | 5 | 扫描字节限制（TB） |
| resource_group_cpu_core | 10 | 资源组CPU核心限制 |
| resource_group_mem_gb | 50 | 资源组内存限制（GB） |
| resource_group_concurrency | 3 | 资源组并发度限制 |
| feishu_proxy | "" | 飞书代理地址 |
| email_suffix | "" | 企业邮箱后缀 |
| data_registration_port | 8030 | 数据落表端口 |

### 2. 集群连接规则配置 (cluster_connection)

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| default_fe_port | 9030 | 默认FE端口 |
| license_expire_remind_days | 30 | License过期提醒天数 |
| license_check_switch | 0 | License检查开关（0=关闭，1=开启） |

### 3. 短查询规则配置 (short_query)

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| default_init_push | 0 | 默认初始化推送（0=不需要，1=需要） |
| default_status | 0 | 默认状态 |
| default_core | 4 | 默认CPU核心数 |
| default_memory_gb | 8 | 默认内存（GB） |

### 4. 飞书机器人配置 (lark_robot)

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| default_status | 0 | 默认状态（0=关闭，1=开启） |

### 5. 企业微信机器人配置 (wecom_robot)

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| default_status | 0 | 默认状态（0=关闭，1=开启） |
| webhook_key | "" | 企业微信Webhook Key |
| msg_type | "markdown" | 消息类型（text/markdown） |
| mention_all | false | 是否@所有人 |
| mentioned_mobile_list | [] | 需要@的成员手机号列表 |

## 企业微信机器人配置说明

### 获取 Webhook Key

1. 在企业微信群中，点击右上角"..." -> "群设置" -> "群机器人"
2. 点击"添加机器人"，选择"自定义"
3. 复制Webhook地址，格式如：`https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxxxxxxxxx`
4. 提取 `key` 参数的值，填入配置文件中的 `webhook_key`

### 消息类型说明

- **text**: 纯文本消息，支持@成员
- **markdown**: Markdown格式消息，支持部分Markdown语法

### 使用示例

```yaml
wecom_robot:
  default_status: 1
  webhook_key: "your-webhook-key-here"
  msg_type: "markdown"
  mention_all: false
  mentioned_mobile_list:
    - "13800138000"
    - "13900139000"
```

## 数据库表结构更新

### 新增表：starrocks_information_wecomrobot

用于存储企业微信机器人配置，支持按集群和用户维度配置：

```sql
CREATE TABLE `starrocks_information_wecomrobot` (
  `type` varchar(100) NOT NULL COMMENT '机器人类型，global,cluster,user',
  `key` varchar(100) NOT NULL COMMENT '机器人集群通知标记',
  `webhook_key` varchar(500) DEFAULT NULL COMMENT '企业微信机器人Webhook Key',
  `msg_type` varchar(50) DEFAULT 'markdown' COMMENT '消息类型',
  `mention_all` int NOT NULL DEFAULT '0' COMMENT '是否@所有人',
  `mentioned_mobile_list` varchar(500) DEFAULT NULL COMMENT '需要@的成员手机号列表',
  `status` int NOT NULL DEFAULT '0' COMMENT '开关',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB COMMENT='企业微信慢查询告警推送机器人';
```

### 更新表：starrocks_information_slowconfig

新增企业微信相关字段：

```sql
ALTER TABLE `starrocks_information_slowconfig`
ADD COLUMN `slow_query_wecom_webhook_key` varchar(500) DEFAULT NULL COMMENT '企业微信机器人Webhook Key',
ADD COLUMN `slow_query_wecom_msg_type` varchar(50) DEFAULT 'markdown' COMMENT '企业微信消息类型',
ADD COLUMN `slow_query_wecom_mention_all` int DEFAULT '0' COMMENT '企业微信是否@所有人',
ADD COLUMN `slow_query_wecom_status` int DEFAULT '0' COMMENT '企业微信告警开关';
```

## 优先级说明

配置的优先级（从高到低）：

1. 数据库中的配置（starrocks_information_slowconfig）
2. 配置文件中的配置（config/starrocks_rule.yaml）
3. 代码中的默认值

## 重启服务

修改配置文件后，需要重启 StarRocksQueris 服务才能生效。

```bash
# 停止服务
kill -9 $(ps aux | grep StarRocksQueris | grep -v grep | awk '{print $2}')

# 启动服务
./StarRocksQueris
```
