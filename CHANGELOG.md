# StarRocksQueris 项目修改总结

## 一、配置文件化改造

### 1.1 新增规则配置文件
**文件**: `config/starrocks_rule.yaml`

**新增内容**:
- 慢查询规则配置（告警时间、查杀时间、并发限制、扫描限制等）
- 集群连接规则配置（默认端口、License检查等）
- StarRocks集群连接配置列表（支持多集群）
- 短查询规则配置
- 飞书机器人规则配置
- 企业微信应用配置（CorpID、AgentId、Secret）
- 机器人配置列表（飞书、企业微信）

**便利性提升**:
- ✅ 无需修改代码即可调整所有阈值参数
- ✅ 支持多集群集中管理
- ✅ 企业微信告警支持App方式（原仅支持Webhook）
- ✅ 配置修改后重启即可生效

### 1.2 配置文件结构定义
**文件**: `util/rule_conf.go`

**新增结构体**:
- `SlowQueryRule` - 慢查询规则
- `ClusterConnRule` - 集群连接规则
- `ShortQueryRule` - 短查询规则
- `LarkRobotRule` - 飞书机器人规则
- `WeComAppRule` - 企业微信应用规则
- `ClusterConfig` - 集群配置
- `LarkRobotConfig` - 飞书机器人配置
- `WeComRobotConfig` - 企业微信机器人配置
- `RobotsConfig` - 机器人配置列表
- `RuleConfig` - 全局规则配置

**新增功能**:
- `LoadRuleConfig()` - 加载规则配置文件
- 支持Viper配置管理
- 支持配置默认值

## 二、自动初始化功能

### 2.1 数据库表自动初始化
**文件**: `Initialize.go`

**新增函数**:
- `InitClustersToDB()` - 自动将集群配置写入数据库
- `InitRobotsToDB()` - 自动将机器人配置写入数据库
- `InitSlowConfigToDB()` - 自动将慢查询配置写入数据库
- `LoadSlowConfigToConnectNorm()` - 加载配置到内存
- `LoadRobotAndClusterConfig()` - 加载机器人和集群配置

**便利性提升**:
- ✅ 无需手动执行INSERT语句
- ✅ 启动时自动同步配置文件到数据库
- ✅ 支持配置更新（已存在则更新，不存在则插入）
- ✅ 减少部署复杂度，降低出错概率

### 2.2 配置热更新修复
**文件**: `trigger.go`, `Initialize.go`

**修复内容**:
- 启动 `tigger` 函数监听配置表变更
- 配置变更后自动重新加载到内存

**便利性提升**:
- ✅ 修改数据库配置后无需重启服务
- ✅ 配置变更实时生效

## 三、企业微信告警增强

### 3.1 企业微信App支持
**文件**: `util/rule_conf.go`, `config/starrocks_rule.yaml`

**新增配置**:
```yaml
wecom_app:
  corp_id: ""
  agent_id: ""
  secret: ""
```

**便利性提升**:
- ✅ 支持企业微信官方API方式发送告警
- ✅ 支持@指定用户
- ✅ 支持多种消息类型（text、markdown）



## 四、部署和文档

### 4.1 部署文档更新
**文件**: `DEPLOY.md`

**新增内容**:
- 详细的编译步骤
- 配置文件说明
- 新增`Db`配置项说明
- 表名格式说明

### 4.2 配置检查脚本
**文件**: `check_config.sh`

**功能**:
- 检查配置文件是否存在
- 检查YAML语法
- 检查文件编码和隐藏字符
- 检查BOM头

**便利性提升**:
- ✅ 快速定位配置问题
- ✅ 减少部署错误

## 五、调试功能增强

### 5.1 调试日志
**文件**: `util/conf.go`

**新增内容**:
- 打印执行目录
- 打印配置文件路径
- 打印配置加载错误详情

**便利性提升**:
- ✅ 快速定位配置加载问题
- ✅ 便于排查路径错误


## 六、核心便利性提升

### 6.1 部署简化
**Before**:
1. 编译代码
2. 手动创建数据库表
3. 手动执行INSERT语句配置集群
4. 手动执行INSERT语句配置机器人
5. 手动执行INSERT语句配置阈值
6. 启动服务

**After**:
1. 编译代码
2. 执行exec.sql创建表
3. 修改config/starrocks_rule.yaml配置文件
4. 启动服务（自动初始化所有配置）

### 6.2 配置管理
**Before**:
- 修改配置需要重新编译
- 配置分散在多个地方
- 无版本管理

**After**:
- 修改YAML文件即可
- 配置集中管理
- 支持版本控制
- 支持多环境配置

### 6.3 告警能力
**Before**:
- 仅支持飞书Webhook
- 企业微信仅支持Webhook

**After**:
- 支持飞书Webhook和应用
- 支持企业微信Webhook和App API
- 支持@指定用户
- 支持多种消息格式

