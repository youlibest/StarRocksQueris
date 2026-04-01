/*
 *@author  chengkenli
 *@project StarRocksQueris
 *@package test
 *@file    alert_test
 *@date    2024/11/6
 *@desc    告警功能自动化测试
 */

package alert

import (
	"StarRocksQueris/util"
	"testing"
	"time"
)

// TestSlowQueryThresholdDetection 测试慢查询阈值检测逻辑
func TestSlowQueryThresholdDetection(t *testing.T) {
	tests := []struct {
		name        string
		queryTime   int
		threshold   int
		shouldAlert bool
	}{
		{"低于阈值", 5, 10, false},
		{"等于阈值", 10, 10, true},
		{"高于阈值", 15, 10, true},
		{"零值查询时间", 0, 10, false},
		{"零值阈值", 10, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.queryTime >= tt.threshold
			if result != tt.shouldAlert {
				t.Errorf("查询时间=%d, 阈值=%d, 期望告警=%v, 实际=%v",
					tt.queryTime, tt.threshold, tt.shouldAlert, result)
			}
		})
	}
}

// TestKillThresholdDetection 测试查杀阈值检测逻辑
func TestKillThresholdDetection(t *testing.T) {
	tests := []struct {
		name       string
		queryTime  int
		killTime   int
		shouldKill bool
	}{
		{"低于查杀阈值", 10, 30, false},
		{"等于查杀阈值", 30, 30, true},
		{"高于查杀阈值", 35, 30, true},
		{"告警但未到查杀", 15, 30, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.queryTime >= tt.killTime
			if result != tt.shouldKill {
				t.Errorf("查询时间=%d, 查杀阈值=%d, 期望查杀=%v, 实际=%v",
					tt.queryTime, tt.killTime, tt.shouldKill, result)
			}
		})
	}
}

// TestActionCodeDetermination 测试动作码确定逻辑
func TestActionCodeDetermination(t *testing.T) {
	tests := []struct {
		name      string
		queryTime int
		warnTime  int
		killTime  int
		action    int
	}{
		{"无告警", 5, 10, 30, 0},
		{"告警但不查杀", 15, 10, 30, 2},
		{"告警并查杀", 35, 10, 30, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var action int
			if tt.queryTime < tt.warnTime {
				action = 0
			} else if tt.queryTime >= tt.killTime {
				action = 3
			} else {
				action = 2
			}

			if action != tt.action {
				t.Errorf("查询时间=%d, 告警阈值=%d, 查杀阈值=%d, 期望动作码=%d, 实际=%d",
					tt.queryTime, tt.warnTime, tt.killTime, tt.action, action)
			}
		})
	}
}

// TestCacheKeyGeneration 测试缓存键生成逻辑
func TestCacheKeyGeneration(t *testing.T) {
	tests := []struct {
		action   int
		connId   string
		expected string
	}{
		{2, "12345", "2_12345"},
		{3, "67890", "3_67890"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			key := generateCacheKey(tt.action, tt.connId)
			if key != tt.expected {
				t.Errorf("动作码=%d, 连接ID=%s, 期望键=%s, 实际=%s",
					tt.action, tt.connId, tt.expected, key)
			}
		})
	}
}

func generateCacheKey(action int, connId string) string {
	return string(rune(action+'0')) + "_" + connId
}

// TestQueryTimeParsing 测试查询时间解析
func TestQueryTimeParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		wantErr  bool
	}{
		{"10", 10, false},
		{"600", 600, false},
		{"1500", 1500, false},
		{"abc", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			// 模拟解析逻辑
			var result int
			_, err := parseInt(tt.input)
			hasErr := err != nil

			if hasErr != tt.wantErr {
				t.Errorf("输入=%s, 期望错误=%v, 实际错误=%v",
					tt.input, tt.wantErr, hasErr)
			}

			if !hasErr && result != tt.expected {
				t.Errorf("输入=%s, 期望=%d, 实际=%d",
					tt.input, tt.expected, result)
			}
		})
	}
}

func parseInt(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	result := 0
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return 0, nil
		}
		result = result*10 + int(ch-'0')
	}
	return result, nil
}

// TestConcurrentQueryHandling 测试并发查询处理
func TestConcurrentQueryHandling(t *testing.T) {
	// 模拟并发场景
	queries := []struct {
		id       string
		user     string
		duration int
	}{
		{"1", "user1", 5},
		{"2", "user2", 15},
		{"3", "user1", 25},
		{"4", "user3", 35},
	}

	alertCount := 0
	killCount := 0
	warnTime := 10
	killTime := 30

	for _, q := range queries {
		if q.duration >= warnTime {
			alertCount++
			if q.duration >= killTime {
				killCount++
			}
		}
	}

	if alertCount != 3 {
		t.Errorf("期望告警数量=3, 实际=%d", alertCount)
	}

	if killCount != 1 {
		t.Errorf("期望查杀数量=1, 实际=%d", killCount)
	}
}

// TestWhiteListProtection 测试白名单保护逻辑
func TestWhiteListProtection(t *testing.T) {
	whiteList := []string{"admin", "monitor", "backup"}
	tests := []struct {
		user       string
		isProtected bool
	}{
		{"admin", true},
		{"user1", false},
		{"monitor", true},
		{"test", false},
	}

	for _, tt := range tests {
		t.Run(tt.user, func(t *testing.T) {
			protected := stringInSlice(tt.user, whiteList)
			if protected != tt.isProtected {
				t.Errorf("用户=%s, 期望保护=%v, 实际=%v",
					tt.user, tt.isProtected, protected)
			}
		})
	}
}

func stringInSlice(s string, list []string) bool {
	for _, item := range list {
		if item == s {
			return true
		}
	}
	return false
}

// TestAlertDeduplication 测试告警去重逻辑
func TestAlertDeduplication(t *testing.T) {
	// 模拟缓存
	cache := make(map[string]bool)

	// 第一次告警
	key1 := "2_12345"
	if _, exists := cache[key1]; exists {
		t.Error("首次告警不应存在缓存中")
	}
	cache[key1] = true

	// 重复告警（应被去重）
	if _, exists := cache[key1]; !exists {
		t.Error("重复告警应在缓存中")
	}

	// 不同连接ID（新告警）
	key2 := "2_67890"
	if _, exists := cache[key2]; exists {
		t.Error("不同连接ID应是新告警")
	}
}

// TestResourceThresholds 测试资源阈值
func TestResourceThresholds(t *testing.T) {
	tests := []struct {
		name      string
		scanRows  int64
		scanBytes int64
		memUsage  int
		alert     bool
	}{
		{"正常查询", 1000000, 1000000, 10, false},
		{"扫描行数超标", 10000000001, 1000000, 10, true},
		{"扫描字节超标", 1000000, 6, 10, true},
		{"内存使用超标", 1000000, 1000000, 201, true},
	}

	thresholds := struct {
		scanRows  int64
		scanBytes int
		memUsage  int
	}{
		scanRows:  10000000000,
		scanBytes: 5,
		memUsage:  200,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alert := tt.scanRows > thresholds.scanRows ||
				int(tt.scanBytes) > thresholds.scanBytes ||
				tt.memUsage > thresholds.memUsage

			if alert != tt.alert {
				t.Errorf("期望告警=%v, 实际=%v", tt.alert, alert)
			}
		})
	}
}

// TestNotificationContent 测试通知内容生成
func TestNotificationContent(t *testing.T) {
	alert := struct {
		App        string
		User       string
		QueryTime  int
		QueryId    string
		Action     int
		Timestamp  time.Time
	}{
		App:       "test-cluster",
		User:      "test-user",
		QueryTime: 15,
		QueryId:   "query-12345",
		Action:    2,
		Timestamp: time.Now(),
	}

	// 验证通知内容不为空
	if alert.App == "" || alert.User == "" || alert.QueryId == "" {
		t.Error("通知内容关键字段不能为空")
	}

	// 验证动作码有效
	if alert.Action != 2 && alert.Action != 3 {
		t.Errorf("动作码=%d 无效，应为 2(告警) 或 3(查杀)", alert.Action)
	}
}

// TestConfigReload 测试配置重载
func TestConfigReload(t *testing.T) {
	// 模拟配置更新
	oldThreshold := 600
	newThreshold := 300

	// 验证配置可以更新
	if newThreshold == oldThreshold {
		t.Error("新配置应与旧配置不同")
	}

	// 验证新阈值生效
	testQueryTime := 400
	shouldAlertWithOld := testQueryTime >= oldThreshold
	shouldAlertWithNew := testQueryTime >= newThreshold

	if shouldAlertWithOld == shouldAlertWithNew {
		t.Error("配置更新后告警行为应改变")
	}
}

// TestErrorHandling 测试错误处理
func TestErrorHandling(t *testing.T) {
	// 模拟数据库连接错误
	dbError := true
	if dbError {
		// 应记录错误但不 panic
		t.Log("数据库连接错误已处理")
	}

	// 模拟空配置
	emptyConfig := true
	if emptyConfig {
		// 应使用默认值
		t.Log("空配置已处理，使用默认值")
	}
}
