/*
 *@author  chengkenli
 *@project StarRocksQueris
 *@package test
 *@file    cluster_init_manual
 *@date    2024/11/6
 *@desc    集群初始化手动测试脚本
 *
 * 使用方法:
 * 1. 确保 config/starrocks_rule.yaml 已配置集群信息
 * 2. 确保数据库连接配置正确
 * 3. 运行: go run test/manual/cluster_init_manual.go
 */

package main

import (
	"StarRocksQueris/util"
	"fmt"
	"os"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("   StarRocksQueris 集群初始化手动测试")
	fmt.Println("========================================")
	fmt.Println()

	// 1. 加载配置文件
	fmt.Println("[步骤1] 加载规则配置文件...")
	err := util.LoadRuleConfig()
	if err != nil {
		fmt.Printf("✗ 加载规则配置失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ 规则配置文件加载成功")
	fmt.Println()

	// 2. 显示加载的集群配置
	fmt.Println("[步骤2] 显示加载的集群配置...")
	clusters := util.GlobalRuleConfig.Clusters
	if len(clusters) == 0 {
		fmt.Println("⚠ 警告: 配置文件中没有集群配置")
		fmt.Println("  请在 config/starrocks_rule.yaml 中配置 clusters 部分")
	} else {
		fmt.Printf("✓ 发现 %d 个集群配置:\n", len(clusters))
		for i, c := range clusters {
			fmt.Printf("\n  [%d] 集群名称: %s\n", i+1, c.App)
			fmt.Printf("      昵称: %s\n", c.Nickname)
			fmt.Printf("      别名: %s\n", c.Alias)
			fmt.Printf("      FE IP: %s\n", c.Feip)
			fmt.Printf("      用户名: %s\n", c.User)
			fmt.Printf("      密码: %s\n", maskPassword(c.Password))
			fmt.Printf("      FE端口: %d\n", c.Feport)
			fmt.Printf("      Manager地址: %s\n", c.Address)
			fmt.Printf("      License过期提醒: %d 天\n", c.Expire)
			fmt.Printf("      License检查开关: %d\n", c.Status)
			fmt.Printf("      FE日志路径: %s\n", c.FeLogPath)
			fmt.Printf("      BE日志路径: %s\n", c.BeLogPath)
		}
	}
	fmt.Println()

	// 3. 显示默认配置
	fmt.Println("[步骤3] 显示默认配置...")
	fmt.Printf("  默认FE端口: %d\n", util.GlobalRuleConfig.ClusterConnection.DefaultFEPort)
	fmt.Printf("  License过期提醒天数: %d\n", util.GlobalRuleConfig.ClusterConnection.LicenseExpireRemindDays)
	fmt.Printf("  License检查开关: %d\n", util.GlobalRuleConfig.ClusterConnection.LicenseCheckSwitch)
	fmt.Println()

	// 4. 显示企业微信配置
	fmt.Println("[步骤4] 显示企业微信应用配置...")
	wecom := util.GlobalRuleConfig.WeComApp
	fmt.Printf("  状态: %d\n", wecom.DefaultStatus)
	fmt.Printf("  CorpID: %s\n", wecom.CorpID)
	fmt.Printf("  AgentID: %s\n", wecom.AgentID)
	fmt.Printf("  Secret: %s\n", maskPassword(wecom.Secret))
	fmt.Printf("  消息类型: %s\n", wecom.MsgType)
	fmt.Printf("  @所有人: %v\n", wecom.MentionAll)
	fmt.Printf("  @用户列表: %v\n", wecom.MentionedUserList)
	fmt.Println()

	// 5. 显示慢查询规则配置
	fmt.Println("[步骤5] 显示慢查询规则配置...")
	slowQuery := util.GlobalRuleConfig.SlowQuery
	fmt.Printf("  告警超时时间: %d 秒\n", slowQuery.AlertTimeoutSeconds)
	fmt.Printf("  查杀超时时间: %d 秒\n", slowQuery.KillTimeoutSeconds)
	fmt.Printf("  并发限制: %d\n", slowQuery.ConcurrencyLimit)
	fmt.Printf("  全表扫描行数限制: %d\n", slowQuery.FullScanRowsLimit)
	fmt.Printf("  Catalog扫描行数限制: %d\n", slowQuery.CatalogScanRowsLimit)
	fmt.Printf("  BE内存限制: %d GB\n", slowQuery.BEMemoryLimitGB)
	fmt.Printf("  扫描行数限制: %d\n", slowQuery.ScanRowsLimit)
	fmt.Printf("  扫描字节限制: %d TB\n", slowQuery.ScanBytesLimitTB)
	fmt.Println()

	// 6. 验证配置有效性
	fmt.Println("[步骤6] 验证配置有效性...")
	validCount := 0
	invalidCount := 0
	for _, c := range clusters {
		if c.App == "" {
			fmt.Printf("  ✗ 集群配置无效: App 为空\n")
			invalidCount++
			continue
		}
		if c.Feip == "" {
			fmt.Printf("  ✗ 集群 %s 配置无效: Feip 为空\n", c.App)
			invalidCount++
			continue
		}
		if c.User == "" {
			fmt.Printf("  ✗ 集群 %s 配置无效: User 为空\n", c.App)
			invalidCount++
			continue
		}
		fmt.Printf("  ✓ 集群 %s 配置有效\n", c.App)
		validCount++
	}
	fmt.Printf("\n  有效配置: %d, 无效配置: %d\n", validCount, invalidCount)
	fmt.Println()

	// 7. 显示预期数据库操作
	fmt.Println("[步骤7] 预期数据库操作...")
	if len(clusters) > 0 {
		fmt.Println("  启动时将执行以下操作:")
		for _, c := range clusters {
			if c.App != "" && c.Feip != "" && c.User != "" {
				feport := c.Feport
				if feport == 0 {
					feport = util.GlobalRuleConfig.ClusterConnection.DefaultFEPort
				}
				expire := c.Expire
				if expire == 0 {
					expire = util.GlobalRuleConfig.ClusterConnection.LicenseExpireRemindDays
				}
				status := c.Status

				fmt.Printf("\n    - 集群: %s\n", c.App)
				fmt.Printf("      操作: INSERT 或 UPDATE cconect 表\n")
				fmt.Printf("      数据:\n")
				fmt.Printf("        app: %s\n", c.App)
				fmt.Printf("        feip: %s\n", c.Feip)
				fmt.Printf("        user: %s\n", c.User)
				fmt.Printf("        feport: %d\n", feport)
				fmt.Printf("        address: %s\n", c.Address)
				fmt.Printf("        expire: %d\n", expire)
				fmt.Printf("        status: %d\n", status)
			}
		}
	} else {
		fmt.Println("  无集群配置，跳过数据库操作")
	}
	fmt.Println()

	fmt.Println("========================================")
	fmt.Println("   手动测试完成")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("提示: 运行主程序时将自动执行集群初始化")
	fmt.Println("      主程序: go run main.go")
	fmt.Println("      或编译后运行: ./StarRocksQueris")
}

// maskPassword 隐藏密码
func maskPassword(password string) string {
	if password == "" {
		return "(空)"
	}
	if len(password) <= 4 {
		return "****"
	}
	return password[:2] + "****" + password[len(password)-2:]
}
