/*
 *@author  chengkenli
 *@project StarRocksQueris
 *@package util
 *@file    rule_conf
 *@date    2024/10/22
 */
package util

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
)

// 慢查询规则配置结构体
type SlowQueryRule struct {
	AlertTimeoutSeconds     int    `mapstructure:"alert_timeout_seconds"`
	KillTimeoutSeconds      int    `mapstructure:"kill_timeout_seconds"`
	ConcurrencyLimit        int    `mapstructure:"concurrency_limit"`
	FullScanRowsLimit       int    `mapstructure:"full_scan_rows_limit"`
	CatalogScanRowsLimit    int    `mapstructure:"catalog_scan_rows_limit"`
	BEMemoryLimitGB         int    `mapstructure:"be_memory_limit_gb"`
	ScanRowsLimit           int64  `mapstructure:"scan_rows_limit"`
	ScanBytesLimitTB        int    `mapstructure:"scan_bytes_limit_tb"`
	ResourceGroupCpuCore    int    `mapstructure:"resource_group_cpu_core"`
	ResourceGroupMemGB      int    `mapstructure:"resource_group_mem_gb"`
	ResourceGroupConcurrency int   `mapstructure:"resource_group_concurrency"`
	FeishuProxy             string `mapstructure:"feishu_proxy"`
	EmailSuffix             string `mapstructure:"email_suffix"`
	DataRegistrationPort    int    `mapstructure:"data_registration_port"`
}

// 集群连接配置结构体
type ClusterConnRule struct {
	DefaultFEPort          int `mapstructure:"default_fe_port"`
	LicenseExpireRemindDays int `mapstructure:"license_expire_remind_days"`
	LicenseCheckSwitch     int `mapstructure:"license_check_switch"`
}

// 短查询配置结构体
type ShortQueryRule struct {
	DefaultInitPush int `mapstructure:"default_init_push"`
	DefaultStatus   int `mapstructure:"default_status"`
	DefaultCore     int `mapstructure:"default_core"`
	DefaultMemoryGB int `mapstructure:"default_memory_gb"`
}

// 飞书机器人配置结构体
type LarkRobotRule struct {
	DefaultStatus int `mapstructure:"default_status"`
}

// 企业微信应用配置结构体（使用CorpID+AgentId+Secret）
type WeComAppRule struct {
	DefaultStatus     int      `mapstructure:"default_status"`
	CorpID            string   `mapstructure:"corp_id"`
	AgentID           string   `mapstructure:"agent_id"`
	Secret            string   `mapstructure:"secret"`
	MsgType           string   `mapstructure:"msg_type"`
	MentionAll        bool     `mapstructure:"mention_all"`
	MentionedUserList []string `mapstructure:"mentioned_user_list"`
}

// 集群配置结构体
type ClusterConfig struct {
	App         string `mapstructure:"app"`
	Nickname    string `mapstructure:"nickname"`
	Alias       string `mapstructure:"alias"`
	Feip        string `mapstructure:"feip"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	Feport      int    `mapstructure:"feport"`
	Address     string `mapstructure:"address"`
	Expire      int    `mapstructure:"expire"`
	Status      int    `mapstructure:"status"`
	FeLogPath   string `mapstructure:"fe_log_path"`
	BeLogPath   string `mapstructure:"be_log_path"`
}

// 飞书机器人配置结构体
type LarkRobotConfig struct {
	Type   string `mapstructure:"type"`
	Key    string `mapstructure:"key"`
	Robot  string `mapstructure:"robot"`
	Status int    `mapstructure:"status"`
}

// 企业微信机器人配置结构体
type WeComRobotConfig struct {
	Type                string `mapstructure:"type"`
	Key                 string `mapstructure:"key"`
	WebhookKey          string `mapstructure:"webhook_key"`
	MsgType             string `mapstructure:"msg_type"`
	MentionAll          int    `mapstructure:"mention_all"`
	MentionedMobileList string `mapstructure:"mentioned_mobile_list"`
	Status              int    `mapstructure:"status"`
}

// 机器人配置列表
type RobotsConfig struct {
	Lark  []LarkRobotConfig  `mapstructure:"lark"`
	WeCom []WeComRobotConfig `mapstructure:"wecom"`
}

// 全局规则配置
type RuleConfig struct {
	SlowQuery         SlowQueryRule   `mapstructure:"slow_query"`
	ClusterConnection ClusterConnRule `mapstructure:"cluster_connection"`
	ShortQuery        ShortQueryRule  `mapstructure:"short_query"`
	LarkRobot         LarkRobotRule   `mapstructure:"lark_robot"`
	WeComApp          WeComAppRule    `mapstructure:"wecom_app"`
	Clusters          []ClusterConfig `mapstructure:"clusters"`
	Robots            RobotsConfig    `mapstructure:"robots"`
}

var GlobalRuleConfig RuleConfig

// 加载规则配置文件
func LoadRuleConfig() error {
	// 获取执行目录，拼接配置路径
	execDir, _ := filepath.Abs(filepath.Dir("./"))
	ruleConfPath := filepath.Join(execDir, "config", "starrocks_rule.yaml")

	// 初始化viper读取规则配置
	ruleViper := viper.New()
	ruleViper.SetConfigFile(ruleConfPath)
	ruleViper.SetConfigType("yaml")

	// 设置默认值
	ruleViper.SetDefault("slow_query.alert_timeout_seconds", 600)
	ruleViper.SetDefault("slow_query.kill_timeout_seconds", 1500)
	ruleViper.SetDefault("slow_query.concurrency_limit", 80)
	ruleViper.SetDefault("slow_query.full_scan_rows_limit", 200000000)
	ruleViper.SetDefault("slow_query.catalog_scan_rows_limit", 100000000)
	ruleViper.SetDefault("slow_query.be_memory_limit_gb", 200)
	ruleViper.SetDefault("slow_query.scan_rows_limit", 10000000000)
	ruleViper.SetDefault("slow_query.scan_bytes_limit_tb", 5)
	ruleViper.SetDefault("slow_query.resource_group_cpu_core", 10)
	ruleViper.SetDefault("slow_query.resource_group_mem_gb", 50)
	ruleViper.SetDefault("slow_query.resource_group_concurrency", 3)
	ruleViper.SetDefault("slow_query.data_registration_port", 8030)
	ruleViper.SetDefault("cluster_connection.default_fe_port", 9030)
	ruleViper.SetDefault("cluster_connection.license_expire_remind_days", 30)
	ruleViper.SetDefault("cluster_connection.license_check_switch", 0)
	ruleViper.SetDefault("short_query.default_init_push", 0)
	ruleViper.SetDefault("short_query.default_status", 0)
	ruleViper.SetDefault("short_query.default_core", 4)
	ruleViper.SetDefault("short_query.default_memory_gb", 8)
	ruleViper.SetDefault("lark_robot.default_status", 0)
	ruleViper.SetDefault("wecom_app.default_status", 0)
	ruleViper.SetDefault("wecom_app.msg_type", "markdown")
	ruleViper.SetDefault("wecom_app.mention_all", false)

	// 读取配置文件
	if err := ruleViper.ReadInConfig(); err != nil {
		// 如果配置文件不存在，使用默认值继续运行
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			Loggrs.Warn("规则配置文件未找到，使用默认配置")
		} else {
			return fmt.Errorf("读取规则配置文件失败: %v", err)
		}
	}

	// 解析到结构体
	if err := ruleViper.Unmarshal(&GlobalRuleConfig); err != nil {
		return fmt.Errorf("解析规则配置失败: %v", err)
	}

	Loggrs.Info("规则配置文件加载成功！")
	return nil
}