# StarRocksQueris 测试文档

## 测试目录结构

```
test/
├── README.md                          # 测试文档
├── integration/                       # 集成测试
│   └── cluster_init_test.go          # 集群初始化集成测试
├── e2e/                              # 端到端测试
│   └── cluster_init_e2e_test.go      # 集群初始化E2E测试
└── manual/                           # 手动测试脚本
    └── cluster_init_manual.go        # 集群初始化手动测试
```

## 测试类型说明

### 1. 单元测试

位于项目根目录和 `util/` 目录下的 `*_test.go` 文件：

- `initialize_test.go` - Initialize.go 的单元测试
- `util/rule_conf_test.go` - 规则配置的单元测试

运行方式：
```bash
# 运行所有单元测试
go test ./...

# 运行指定包的测试
go test ./util

# 运行指定测试函数
go test -run TestLoadRuleConfig

# 显示详细输出
go test -v ./...
```

### 2. 集成测试

位于 `test/integration/` 目录：

- `cluster_init_test.go` - 测试集群初始化流程，无需真实数据库

运行方式：
```bash
go test ./test/integration/...
```

### 3. 端到端测试 (E2E)

位于 `test/e2e/` 目录：

- `cluster_init_e2e_test.go` - 完整的端到端测试，需要真实数据库连接

**注意**：E2E 测试默认被跳过，需要取消 `t.Skip()` 才能运行。

运行方式：
```bash
# 先取消测试文件中的 Skip，然后运行
go test ./test/e2e/... -v
```

### 4. 手动测试

位于 `test/manual/` 目录：

- `cluster_init_manual.go` - 手动验证配置文件加载和显示

运行方式：
```bash
go run test/manual/cluster_init_manual.go
```

## 集群初始化测试场景

### 场景1: 首次启动初始化

**目的**: 验证首次启动时集群配置正确写入数据库

**步骤**:
1. 准备 config/starrocks_rule.yaml，配置 clusters 部分
2. 确保数据库 cconect 表为空或不存在该集群
3. 启动服务
4. 验证数据库中是否正确插入记录

**预期结果**:
- 数据库 cconect 表中新增集群记录
- 所有字段值与配置文件一致
- 日志显示 "集群 xxx 配置已插入"

### 场景2: 重启服务配置更新

**目的**: 验证重启服务时配置更新正确

**步骤**:
1. 修改 config/starrocks_rule.yaml 中的集群配置（如IP、密码）
2. 重启服务
3. 验证数据库中配置是否更新

**预期结果**:
- 数据库中对应集群记录被更新
- 日志显示 "集群 xxx 配置已更新"

### 场景3: 多集群配置

**目的**: 验证多个集群同时初始化

**步骤**:
1. 在配置文件中配置多个集群
2. 启动服务
3. 验证所有集群都正确写入数据库

**预期结果**:
- 所有集群配置都写入数据库
- 每个集群都有独立的记录

### 场景4: 无效配置过滤

**目的**: 验证不完整的配置被正确跳过

**步骤**:
1. 在配置文件中添加不完整的集群配置（缺少 App、Feip 或 User）
2. 启动服务
3. 观察日志输出

**预期结果**:
- 不完整的配置被跳过
- 日志显示警告信息 "集群配置不完整，跳过"
- 其他有效配置正常处理

### 场景5: 默认值填充

**目的**: 验证未指定的字段使用默认值

**步骤**:
1. 配置集群时只填写必填字段（App, Feip, User, Password）
2. 启动服务
3. 验证数据库中的记录

**预期结果**:
- Feport 默认为 9030
- Expire 默认为 30
- Status 默认为 0

## 配置文件验证

### 最小配置示例

```yaml
clusters:
  - app: "my-cluster"
    feip: "192.168.1.100"
    user: "root"
    password: "my-password"
```

### 完整配置示例

```yaml
clusters:
  - app: "prod-cluster"
    nickname: "生产集群"
    alias: "prod"
    feip: "192.168.1.100"
    user: "root"
    password: "secure-password"
    feport: 9030
    address: "http://manager.example.com"
    expire: 30
    status: 1
    fe_log_path: "/var/log/starrocks/fe"
    be_log_path: "/var/log/starrocks/be"
```

## 常见问题排查

### 问题1: 配置文件未找到

**现象**: 日志显示 "规则配置文件未找到，使用默认配置"

**解决**: 确保 config/starrocks_rule.yaml 存在于程序执行目录

### 问题2: 集群配置未写入数据库

**现象**: 启动后数据库中没有集群记录

**排查步骤**:
1. 检查配置文件 clusters 部分是否正确
2. 检查数据库连接是否正常
3. 查看日志是否有错误信息
4. 确认 cconect 表是否存在

### 问题3: 配置更新不生效

**现象**: 修改配置文件后重启，数据库中配置未更新

**排查步骤**:
1. 确认配置文件保存成功
2. 检查集群名称 App 是否变更（App 是主键）
3. 查看日志确认是否执行了更新操作

## 测试 checklist

- [ ] 单元测试全部通过
- [ ] 集成测试全部通过
- [ ] 手动测试脚本正常运行
- [ ] 首次启动集群配置正确写入
- [ ] 重启服务配置更新正确
- [ ] 多集群配置全部生效
- [ ] 无效配置被正确过滤
- [ ] 默认值正确填充
- [ ] 日志输出清晰完整
