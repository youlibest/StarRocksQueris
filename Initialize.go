/*
 *@author  chengkenli
 *@project StarRocksQueris
 *@package main
 *@file    Initialize
 *@date    2024/11/6 22:41
 */

package main

import (
	"StarRocksQueris/conn"
	"StarRocksQueris/util"
	"fmt"
	"os"
)

// 慢查询审计日志表
func audit(tablename string) string {
	stmt := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
ts                date            NOT NULL COMMENT "数据载入日期",
app               varchar(64)     NOT NULL COMMENT "集群名称",
queryId           varchar(64)     NOT NULL COMMENT "查询的唯一ID",
origin            varchar(500)        NULL COMMENT "SQL原始语句中涉及到的表",
domain            varchar(64)         NULL COMMENT "主题域简称",
owner             varchar(64)         NULL COMMENT "主题域owner",
action            int(11)             NULL COMMENT "拦截行为：⓿.状态异常停留清退,①.异常违规参数查杀,②.10分钟慢查询提醒,③.30分钟慢查询查杀,④.全表扫描亿级查杀,⑤.TB级扫描字节查杀,⑥.百亿扫描行数查杀,⑦.CATALOG违规查杀,⑧.GB级消耗内存查杀",
timestamp         datetime        NOT NULL COMMENT "查询开始时间",
queryType         varchar(12)         NULL COMMENT "查询类型（query, slow_query,connection）",
clientIp          varchar(32)         NULL COMMENT "客户端IP",
user              varchar(64)         NULL COMMENT "查询用户名",
authorizedUser    varchar(64)         NULL COMMENT "用户唯一标识，既user_identity",
resourceGroup     varchar(64)         NULL COMMENT "资源组名",
catalog           varchar(32)         NULL COMMENT "数据目录名",
db                varchar(96)         NULL COMMENT "查询所在数据库",
state             varchar(8)          NULL COMMENT "查询状态（EOF，ERR，OK）",
errorCode         varchar(96)         NULL COMMENT "错误码",
queryTime         bigint(20)          NULL COMMENT "查询执行时间（秒）",
scanBytes         bigint(20)          NULL COMMENT "查询扫描的字节数",
scanRows          bigint(20)          NULL COMMENT "查询扫描的记录行数",
returnRows        bigint(20)          NULL COMMENT "查询返回的结果行数",
cpuCostNs         bigint(20)          NULL COMMENT "查询CPU耗时（纳秒）",
memCostBytes      bigint(20)          NULL COMMENT "查询消耗内存（字节）",
stmtId            int(11)             NULL COMMENT "SQL语句增量ID",
isQuery           tinyint(4)          NULL COMMENT "SQL是否为查询（1或0）",
feIp              varchar(32)         NULL COMMENT "执行该语句的FE IP",
stmt              varchar(1048576)    NULL COMMENT "SQL原始语句",
digest            varchar(32)         NULL COMMENT "慢SQL指纹",
planCpuCosts      double              NULL COMMENT "查询规划阶段CPU占用（纳秒）",
planMemCosts      double              NULL COMMENT "查询规划阶段内存占用（字节）",
pendingTimeMs     bigint(20)          NULL COMMENT "查询在队列中等待的时间（毫秒）",
logfile           varchar(200)        NULL COMMENT "日志文件",
optimization      int(11)             NULL COMMENT "是否做过优化：0.没有优化，1~99.代表优化次数",
optimizationItems varchar(1048576)    NULL COMMENT "优化项"
) ENGINE=OLAP
PRIMARY KEY(ts, app, queryId)
COMMENT "慢查询审计日志表"
PARTITION BY date_trunc('day', ts)
DISTRIBUTED BY HASH(ts, queryId)`, tablename)
	return stmt
}

func standard(tablename string) string {
	stmt := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
  slow_query_time int NOT NULL DEFAULT '600' COMMENT '慢查询语句的超时告警时间，单位秒。',
  slow_query_ktime int NOT NULL DEFAULT '1500' COMMENT '慢查询语句的查杀时间，单位秒',
  slow_query_concurrencylimit int NOT NULL DEFAULT '80' COMMENT '慢查询的并发度（比如并发语句超过该值则告警），单位整数',
  slow_query_version varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '程序版本号',
  slow_query_focususer varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '慢查询保护白名单用户，使用英文逗号,隔开',
  slow_query_proxy_feishu varchar(300) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '访问飞书代理地址(使用飞书发送信息时，企业需要代理)',
  slow_query_grafana varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT 'prometheus地址，支持向prometheus中推送记录',
  slow_query_lark_app varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '飞书应用名称（企业版）',
  slow_query_lark_appid varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '飞书应用Appid',
  slow_query_lark_appsecret varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '飞书应用AppSecret',
  slow_query_email_host varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '企业邮箱，服务器，host:port',
  slow_query_email_from varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '企业邮箱，用于发送邮件的邮箱',
  slow_query_email_to varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '企业邮箱，用于接收邮件的邮箱, 逗号分隔',
  slow_query_email_cc varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '企业邮箱，由于抄送邮件给cc的邮箱，逗号分隔',
  slow_query_email_bc varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '企业邮箱，由于密送邮件给bc的邮箱，逗号分隔',
  slow_query_email_suffix varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '企业邮件的后缀名，@xxxxx.com',
  slow_query_email_reference_material varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '邮件中呈现的参考资料了解，支持html，逗号分隔',
  slow_query_frontend_avgs varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'parallel_fragment_exec_instance_num=15,query_mem_limit=274877906944,load_mem_limit=274877906944,exec_mem_limit=274877906944' COMMENT '慢查询需要拦截的参数指标比如，key=value,.... 可填多个',
  slow_query_frontend_fullscan_num int DEFAULT '200000000' COMMENT '慢查询拦截全表扫描的最大行数，默认值2亿',
  slow_query_frontend_insert_catalog_scanrow int DEFAULT '100000000' COMMENT '慢查询拦截catalog扫描数据量超过亿级 + INSERT TABLE FROM CATALOG',
  slow_query_frontend_memoryusage int DEFAULT '200' COMMENT '慢查询拦截单个BE 200GB+级别查询消耗内存',
  slow_query_frontend_scanrows bigint DEFAULT '10000000000' COMMENT '慢查询拦截百亿+级别扫描行数',
  slow_query_frontend_scanbytes int DEFAULT '5' COMMENT '慢查询拦截TB+级别扫描字节消耗',
  slow_query_data_registration_username varchar(100) DEFAULT NULL COMMENT '慢查询记录落表，用户名',
  slow_query_data_registration_password varchar(500) DEFAULT NULL COMMENT '慢查询记录落表，密码',
  slow_query_data_registration_table varchar(500) DEFAULT NULL COMMENT '慢查询记录落表，表名',
  slow_query_data_registration_host varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '慢查询记录落表，主机名（FE IP）',
  slow_query_data_registration_port int DEFAULT '8030' COMMENT '慢查询记录落表，端口(因为这个走的是stream load，所以端口默认8030)',
  slow_query_resource_group_cpu_core_limit int DEFAULT '10' COMMENT '慢查询拦截资源隔离，CPU',
  slow_query_resource_group_mem_limit int DEFAULT '50' COMMENT '慢查询拦截资源隔离，内存',
  slow_query_resource_group_concurrency_limit int DEFAULT '3' COMMENT '慢查询拦截资源隔离，并发度',
  updated_at timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;`, tablename)
	return stmt
}

func cconect(tablename string) string {
	stmt := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
  app           varchar(100)    NOT NULL COMMENT '集群名称(英文)',
  feip          varchar(200)    CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '集群连接地址(必填)F5,VIP,CLB,FE',
  user          varchar(200)    CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '集群登录账号(必填) 建议是管理员角色的账号',
  password      varchar(500)    NOT NULL COMMENT '集群登录密码(必填)',
  feport        int             NOT NULL DEFAULT '9030' COMMENT '集群登录端口，默认9030',
  address       varchar(500)    CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT 'MANAGER地址，如果填了MANAGER地址，那么将触发定时检查LICENSE是否过期(企业级)',
  expire        int             DEFAULT '30' COMMENT 'LICENSE是否过期(企业级)过期提醒倒计时，单位day',
  status        int             NOT NULL DEFAULT '0' COMMENT 'LICENSE是否过期(企业级)开关,0 off, 1 on',
  updated_at    timestamp       NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='StarRocks登录配置，manager地址,(定期检查license过期日期)'`, tablename)
	return stmt
}

func crobot(tablename string) string {
	stmt := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
  type          varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '机器人类型，global,cluster,user',
  key           varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '机器人集群通知标记',
  robot         varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '飞书机器人KEY',
  status        int NOT NULL DEFAULT '0' COMMENT '开关',
  updated_at    timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='慢查询告警推送机器人'`, tablename)
	return stmt
}

func cwecomrobot(tablename string) string {
	stmt := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
  type          varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '机器人类型，global,cluster,user',
  key           varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '机器人集群通知标记',
  webhook_key   varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '企业微信机器人Webhook Key',
  msg_type      varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'markdown' COMMENT '消息类型：text(文本), markdown(markdown格式)',
  mention_all   int NOT NULL DEFAULT '0' COMMENT '是否@所有人，0=否，1=是',
  mentioned_mobile_list varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '需要@的成员手机号列表，逗号分隔',
  status        int NOT NULL DEFAULT '0' COMMENT '开关，0=关闭，1=开启',
  updated_at    timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='企业微信慢查询告警推送机器人'`, tablename)
	return stmt
}

// ConfigDB 初始化配置数据库连接
func ConfigDB() {
	InitGorm()
	util.Loggrs.Info("配置数据库连接成功!")

	// 初始化 ConnectNorm 如果为 nil
	if util.ConnectNorm == nil {
		util.ConnectNorm = &util.DBConfig{}
	}

	// 自动初始化集群配置到数据库
	if err := InitClustersToDB(); err != nil {
		util.Loggrs.Error("初始化集群配置到数据库失败: ", err)
	}

	// 自动初始化机器人配置到数据库
	if err := InitRobotsToDB(); err != nil {
		util.Loggrs.Error("初始化机器人配置到数据库失败: ", err)
	}

	// 自动初始化慢查询配置到数据库
	if err := InitSlowConfigToDB(); err != nil {
		util.Loggrs.Error("初始化慢查询配置到数据库失败: ", err)
	}

	// 加载慢查询配置到 ConnectNorm
	if err := LoadSlowConfigToConnectNorm(); err != nil {
		util.Loggrs.Error("加载慢查询配置到 ConnectNorm 失败: ", err)
	}

	// 启动配置热更新触发器
	go tigger(util.Connect, "starrocks_information_larkrobot", 1)
	go tigger(util.Connect, "starrocks_information_connections", 2)
	go tigger(util.Connect, "starrocks_information_slowconfig", 3)
	util.Loggrs.Info("配置热更新触发器已启动")

	// 立即加载一次机器人和集群配置
	LoadRobotAndClusterConfig()
}

// InitClustersToDB 将配置文件中的集群配置自动写入数据库
func InitClustersToDB() error {
	// 检查是否有集群配置
	if len(util.GlobalRuleConfig.Clusters) == 0 {
		util.Loggrs.Info("配置文件中没有集群配置，跳过自动初始化")
		return nil
	}

	// 获取默认配置
	defaultPort := util.GlobalRuleConfig.ClusterConnection.DefaultFEPort
	if defaultPort == 0 {
		defaultPort = 9030
	}
	defaultExpire := util.GlobalRuleConfig.ClusterConnection.LicenseExpireRemindDays
	if defaultExpire == 0 {
		defaultExpire = 30
	}
	defaultStatus := util.GlobalRuleConfig.ClusterConnection.LicenseCheckSwitch

	// 获取数据库连接
	db := util.Connect
	if db == nil {
		return fmt.Errorf("数据库连接未初始化")
	}

	// 遍历集群配置并写入数据库
	for _, cluster := range util.GlobalRuleConfig.Clusters {
		// 跳过空配置
		if cluster.App == "" || cluster.Feip == "" || cluster.User == "" {
			util.Loggrs.Warnf("集群配置不完整，跳过: app=%s", cluster.App)
			continue
		}

		// 使用默认值填充
		feport := cluster.Feport
		if feport == 0 {
			feport = defaultPort
		}
		expire := cluster.Expire
		if expire == 0 {
			expire = defaultExpire
		}
		status := cluster.Status
		if status == 0 && cluster.Status == 0 {
			status = defaultStatus
		}

		// 检查集群是否已存在
		var count int64
		result := db.Table("starrocks_information_connections").Where("app = ?", cluster.App).Count(&count)
		if result.Error != nil {
			util.Loggrs.Errorf("检查集群 %s 是否存在失败: %v", cluster.App, result.Error)
			continue
		}

		if count > 0 {
			// 更新现有集群配置
			result = db.Table("starrocks_information_connections").Where("app = ?", cluster.App).Updates(map[string]interface{}{
				"feip":     cluster.Feip,
				"user":     cluster.User,
				"password": cluster.Password,
				"feport":   feport,
				"address":  cluster.Address,
				"expire":   expire,
				"status":   status,
			})
			if result.Error != nil {
				util.Loggrs.Errorf("更新集群 %s 配置失败: %v", cluster.App, result.Error)
				continue
			}
			util.Loggrs.Infof("集群 %s 配置已更新", cluster.App)
		} else {
			// 插入新集群配置
			result = db.Table("starrocks_information_connections").Create(map[string]interface{}{
				"app":      cluster.App,
				"feip":     cluster.Feip,
				"user":     cluster.User,
				"password": cluster.Password,
				"feport":   feport,
				"address":  cluster.Address,
				"expire":   expire,
				"status":   status,
			})
			if result.Error != nil {
				util.Loggrs.Errorf("插入集群 %s 配置失败: %v", cluster.App, result.Error)
				continue
			}
			util.Loggrs.Infof("集群 %s 配置已插入", cluster.App)
		}
	}

	util.Loggrs.Info("集群配置自动初始化完成")
	return nil
}

// InitRobotsToDB 将配置文件中的机器人配置自动写入数据库
func InitRobotsToDB() error {
	db := util.Connect
	if db == nil {
		return fmt.Errorf("数据库连接未初始化")
	}

	// 初始化飞书机器人配置
	for _, robot := range util.GlobalRuleConfig.Robots.Lark {
		if robot.Type == "" || robot.Key == "" {
			continue
		}

		var count int64
		result := db.Table("starrocks_information_larkrobot").Where("`type` = ? AND `key` = ?", robot.Type, robot.Key).Count(&count)
		if result.Error != nil {
			util.Loggrs.Errorf("检查飞书机器人 %s/%s 是否存在失败: %v", robot.Type, robot.Key, result.Error)
			continue
		}

		if count > 0 {
			result = db.Table("starrocks_information_larkrobot").Where("`type` = ? AND `key` = ?", robot.Type, robot.Key).Updates(map[string]interface{}{
				"robot":  robot.Robot,
				"status": robot.Status,
			})
			if result.Error != nil {
				util.Loggrs.Errorf("更新飞书机器人 %s/%s 配置失败: %v", robot.Type, robot.Key, result.Error)
				continue
			}
			util.Loggrs.Infof("飞书机器人 %s/%s 配置已更新", robot.Type, robot.Key)
		} else {
			result = db.Table("starrocks_information_larkrobot").Create(map[string]interface{}{
				"type":   robot.Type,
				"key":    robot.Key,
				"robot":  robot.Robot,
				"status": robot.Status,
			})
			if result.Error != nil {
				util.Loggrs.Errorf("插入飞书机器人 %s/%s 配置失败: %v", robot.Type, robot.Key, result.Error)
				continue
			}
			util.Loggrs.Infof("飞书机器人 %s/%s 配置已插入", robot.Type, robot.Key)
		}
	}

	// 初始化企业微信机器人配置
	for _, robot := range util.GlobalRuleConfig.Robots.WeCom {
		if robot.Type == "" || robot.Key == "" {
			continue
		}

		var count int64
		result := db.Table("starrocks_information_wecomrobot").Where("`type` = ? AND `key` = ?", robot.Type, robot.Key).Count(&count)
		if result.Error != nil {
			util.Loggrs.Errorf("检查企业微信机器人 %s/%s 是否存在失败: %v", robot.Type, robot.Key, result.Error)
			continue
		}

		if count > 0 {
			result = db.Table("starrocks_information_wecomrobot").Where("`type` = ? AND `key` = ?", robot.Type, robot.Key).Updates(map[string]interface{}{
				"webhook_key":             robot.WebhookKey,
				"msg_type":                robot.MsgType,
				"mention_all":             robot.MentionAll,
				"mentioned_mobile_list":   robot.MentionedMobileList,
				"status":                  robot.Status,
			})
			if result.Error != nil {
				util.Loggrs.Errorf("更新企业微信机器人 %s/%s 配置失败: %v", robot.Type, robot.Key, result.Error)
				continue
			}
			util.Loggrs.Infof("企业微信机器人 %s/%s 配置已更新", robot.Type, robot.Key)
		} else {
			result = db.Table("starrocks_information_wecomrobot").Create(map[string]interface{}{
				"type":                    robot.Type,
				"key":                     robot.Key,
				"webhook_key":             robot.WebhookKey,
				"msg_type":                robot.MsgType,
				"mention_all":             robot.MentionAll,
				"mentioned_mobile_list":   robot.MentionedMobileList,
				"status":                  robot.Status,
			})
			if result.Error != nil {
				util.Loggrs.Errorf("插入企业微信机器人 %s/%s 配置失败: %v", robot.Type, robot.Key, result.Error)
				continue
			}
			util.Loggrs.Infof("企业微信机器人 %s/%s 配置已插入", robot.Type, robot.Key)
		}
	}

	util.Loggrs.Info("机器人配置自动初始化完成")
	return nil
}

// InitSlowConfigToDB 初始化慢查询配置到数据库
func InitSlowConfigToDB() error {
	db := util.Connect
	if db == nil {
		return fmt.Errorf("数据库连接未初始化")
	}

	// 检查是否已有配置
	var count int64
	result := db.Table("starrocks_information_slowconfig").Count(&count)
	if result.Error != nil {
		util.Loggrs.Errorf("检查慢查询配置是否存在失败: %v", result.Error)
		return result.Error
	}

	// 获取配置值
	slowQuery := util.GlobalRuleConfig.SlowQuery
	wecomApp := util.GlobalRuleConfig.WeComApp

	config := map[string]interface{}{
		"slow_query_time":                                slowQuery.AlertTimeoutSeconds,
		"slow_query_ktime":                               slowQuery.KillTimeoutSeconds,
		"slow_query_concurrencylimit":                    slowQuery.ConcurrencyLimit,
		"slow_query_proxy_feishu":                        slowQuery.FeishuProxy,
		"slow_query_frontend_avgs":                       "parallel_fragment_exec_instance_num=15,query_mem_limit=274877906944,load_mem_limit=274877906944,exec_mem_limit=274877906944",
		"slow_query_frontend_fullscan_num":               slowQuery.FullScanRowsLimit,
		"slow_query_frontend_insert_catalog_scanrow":     slowQuery.CatalogScanRowsLimit,
		"slow_query_frontend_memoryusage":                slowQuery.BEMemoryLimitGB,
		"slow_query_frontend_scanrows":                   slowQuery.ScanRowsLimit,
		"slow_query_frontend_scanbytes":                  slowQuery.ScanBytesLimitTB,
		"slow_query_resource_group_cpu_core_limit":       slowQuery.ResourceGroupCpuCore,
		"slow_query_resource_group_mem_limit":            slowQuery.ResourceGroupMemGB,
		"slow_query_resource_group_concurrency_limit":    slowQuery.ResourceGroupConcurrency,
		"slow_query_data_registration_port":              slowQuery.DataRegistrationPort,
		"slow_query_email_suffix":                        slowQuery.EmailSuffix,
		"slow_query_wecom_webhook_key":                   "",
		"slow_query_wecom_msg_type":                      wecomApp.MsgType,
		"slow_query_wecom_mention_all":                   0,
		"slow_query_wecom_status":                        wecomApp.DefaultStatus,
	}

	if count > 0 {
		// 更新现有配置
		result = db.Table("starrocks_information_slowconfig").Updates(config)
		if result.Error != nil {
			util.Loggrs.Errorf("更新慢查询配置失败: %v", result.Error)
			return result.Error
		}
		util.Loggrs.Info("慢查询配置已更新")
	} else {
		// 插入新配置
		result = db.Table("starrocks_information_slowconfig").Create(config)
		if result.Error != nil {
			util.Loggrs.Errorf("插入慢查询配置失败: %v", result.Error)
			return result.Error
		}
		util.Loggrs.Info("慢查询配置已插入")
	}

	return nil
}

// LoadSlowConfigToConnectNorm 从数据库加载慢查询配置到 ConnectNorm
func LoadSlowConfigToConnectNorm() error {
	db := util.Connect
	if db == nil {
		return fmt.Errorf("数据库连接未初始化")
	}

	// 从数据库读取慢查询配置
	result := db.Table("starrocks_information_slowconfig").First(util.ConnectNorm)
	if result.Error != nil {
		util.Loggrs.Errorf("加载慢查询配置到 ConnectNorm 失败: %v", result.Error)
		return result.Error
	}

	util.Loggrs.Info("慢查询配置已加载到 ConnectNorm")
	util.Loggrs.Infof("SlowQueryTime: %d, SlowQueryKtime: %d", util.ConnectNorm.SlowQueryTime, util.ConnectNorm.SlowQueryKtime)
	return nil
}

// LoadRobotAndClusterConfig 立即加载机器人和集群配置
func LoadRobotAndClusterConfig() {
	db := util.Connect
	if db == nil {
		util.Loggrs.Error("数据库连接未初始化，无法加载机器人和集群配置")
		return
	}

	// 加载机器人配置
	result := db.Raw("select * from starrocks_information_larkrobot where status > 0").Scan(&util.ConnectRobot)
	if result.Error != nil {
		util.Loggrs.Errorf("加载飞书机器人配置失败: %v", result.Error)
	} else {
		util.Loggrs.Infof("已加载 %d 个飞书机器人配置", len(util.ConnectRobot))
	}

	// 加载集群连接配置
	result = db.Raw("select * from starrocks_information_connections where status >= 0").Scan(&util.ConnectBody)
	if result.Error != nil {
		util.Loggrs.Errorf("加载集群连接配置失败: %v", result.Error)
	} else {
		util.Loggrs.Infof("已加载 %d 个集群连接配置", len(util.ConnectBody))
	}
}

// InitGorm 初始化MySQL配置库连接
func InitGorm() {
	db, err := conn.ConnectMySQL()
	if err != nil {
		util.Loggrs.Error("连接配置数据库失败: ", err)
		os.Exit(1)
	}
	util.Connect = db
}
