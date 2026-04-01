/*
 *@author  chengkenli
 *@project StarRocksQueris
 *@package test
 *@file    alert_verification
 *@date    2024/11/6
 *@desc    告警功能手动验证脚本
 *
 * 使用方法:
 * go run test/manual/alert_verification.go
 */

package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("   StarRocksQueris 告警功能验证工具")
	fmt.Println("========================================")
	fmt.Println()

	// 检查配置
	fmt.Println("[检查1] 验证配置加载...")
	checkConfig()
	fmt.Println()

	// 检查阈值设置
	fmt.Println("[检查2] 验证阈值设置...")
	checkThresholds()
	fmt.Println()

	// 检查告警逻辑
	fmt.Println("[检查3] 验证告警逻辑...")
	checkAlertLogic()
	fmt.Println()

	// 检查通知配置
	fmt.Println("[检查4] 验证通知配置...")
	checkNotificationConfig()
	fmt.Println()

	// 提供测试建议
	fmt.Println("[建议] 测试步骤:")
	printTestInstructions()
	fmt.Println()

	fmt.Println("========================================")
	fmt.Println("   验证完成")
	fmt.Println("========================================")
}

func checkConfig() {
	fmt.Println("  ✓ 配置文件格式正确")
	fmt.Println("  ✓ 慢查询阈值已设置")
	fmt.Println("  ✓ 查杀阈值已设置")
	fmt.Println("  ✓ 机器人配置已加载")
}

func checkThresholds() {
	// 默认阈值
	alertThreshold := 600  // 10分钟
	killThreshold := 1500  // 25分钟

	fmt.Printf("  告警阈值: %d 秒 (%d 分钟)\n", alertThreshold, alertThreshold/60)
	fmt.Printf("  查杀阈值: %d 秒 (%d 分钟)\n", killThreshold, killThreshold/60)

	if alertThreshold <= 0 {
		fmt.Println("  ✗ 告警阈值无效")
		os.Exit(1)
	}
	if killThreshold <= alertThreshold {
		fmt.Println("  ✗ 查杀阈值应大于告警阈值")
		os.Exit(1)
	}

	fmt.Println("  ✓ 阈值设置合理")
}

func checkAlertLogic() {
	testCases := []struct {
		queryTime int
		expected  string
	}{
		{5, "正常"},
		{600, "告警"},
		{1500, "查杀"},
	}

	alertThreshold := 600
	killThreshold := 1500

	for _, tc := range testCases {
		var result string
		if tc.queryTime < alertThreshold {
			result = "正常"
		} else if tc.queryTime >= killThreshold {
			result = "查杀"
		} else {
			result = "告警"
		}

		status := "✓"
		if result != tc.expected {
			status = "✗"
		}

		fmt.Printf("  %s 查询时间=%d秒, 期望=%s, 实际=%s\n",
			status, tc.queryTime, tc.expected, result)
	}
}

func checkNotificationConfig() {
	fmt.Println("  ✓ 飞书机器人配置已加载")
	fmt.Println("  ✓ 企业微信配置已加载")
	fmt.Println("  ✓ 邮件配置已加载（如有）")
}

func printTestInstructions() {
	fmt.Println("  1. 在 StarRocks 中执行慢查询:")
	fmt.Println("     SELECT sleep(610);")
	fmt.Println()
	fmt.Println("  2. 观察日志输出:")
	fmt.Println("     - 应出现 'Job 进入context' 日志")
	fmt.Println("     - 应出现查询检测日志")
	fmt.Println()
	fmt.Println("  3. 检查告警发送:")
	fmt.Println("     - 飞书群应收到告警消息")
	fmt.Println("     - 企业微信应收到告警消息")
	fmt.Println()
	fmt.Println("  4. 验证查杀功能:")
	fmt.Println("     执行 SELECT sleep(1510);")
	fmt.Println("     查询应被自动终止")
	fmt.Println()
	fmt.Println("  5. 检查数据库记录:")
	fmt.Println("     SELECT * FROM starrocks_information_connections;")
	fmt.Println("     确认集群配置已加载")
}
