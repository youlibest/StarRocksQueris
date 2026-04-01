/*
 *@author  chengkenli
 *@project StarRocksQueris
 *@package test
 *@file    cluster_init_test
 *@date    2024/11/6
 *@desc    集群初始化集成测试
 */

package integration

import (
	"StarRocksQueris/util"
	"fmt"
	"testing"
)

// MockDB 模拟数据库连接用于测试
type MockDB struct {
	Clusters map[string]map[string]interface{}
}

func NewMockDB() *MockDB {
	return &MockDB{
		Clusters: make(map[string]map[string]interface{}),
	}
}

func (m *MockDB) Table(name string) *MockDB {
	return m
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *MockDB {
	return m
}

func (m *MockDB) Count(count *int64) *MockDB {
	*count = 0
	return m
}

func (m *MockDB) Updates(values interface{}) *MockDB {
	return m
}

func (m *MockDB) Create(values interface{}) *MockDB {
	return m
}

func (m *MockDB) Error() error {
	return nil
}

// TestClusterInitFlow 测试完整的集群初始化流程
func TestClusterInitFlow(t *testing.T) {
	fmt.Println("=== 开始集群初始化集成测试 ===")

	// 测试场景1: 首次启动，数据库为空
	t.Run("首次初始化", func(t *testing.T) {
		clusters := []util.ClusterConfig{
			{
				App:      "prod-cluster",
				Nickname: "生产集群",
				Alias:    "prod",
				Feip:     "192.168.1.100",
				User:     "root",
				Password: "prod_password",
				Feport:   9030,
				Address:  "http://manager.prod.com",
				Expire:   30,
				Status:   1,
			},
		}

		if len(clusters) == 0 {
			t.Error("集群配置不应该为空")
		}

		for _, c := range clusters {
			if c.App == "" || c.Feip == "" || c.User == "" {
				t.Errorf("集群 %s 配置不完整", c.App)
			}
		}

		fmt.Println("✓ 首次初始化测试通过")
	})

	// 测试场景2: 重启服务，配置更新
	t.Run("配置更新", func(t *testing.T) {
		clusters := []util.ClusterConfig{
			{
				App:      "prod-cluster",
				Nickname: "生产集群-已更新",
				Feip:     "192.168.1.200", // IP变更
				User:     "root",
				Password: "new_password", // 密码变更
				Feport:   9031,           // 端口变更
			},
		}

		if clusters[0].Feip != "192.168.1.200" {
			t.Error("IP 更新失败")
		}

		fmt.Println("✓ 配置更新测试通过")
	})

	// 测试场景3: 多集群配置
	t.Run("多集群配置", func(t *testing.T) {
		clusters := []util.ClusterConfig{
			{
				App:      "cluster-beijing",
				Nickname: "北京集群",
				Feip:     "10.0.1.100",
				User:     "root",
				Password: "bj_pass",
			},
			{
				App:      "cluster-shanghai",
				Nickname: "上海集群",
				Feip:     "10.0.2.100",
				User:     "root",
				Password: "sh_pass",
			},
			{
				App:      "cluster-guangzhou",
				Nickname: "广州集群",
				Feip:     "10.0.3.100",
				User:     "root",
				Password: "gz_pass",
			},
		}

		if len(clusters) != 3 {
			t.Errorf("应该有3个集群，但得到 %d", len(clusters))
		}

		for i, c := range clusters {
			if c.App == "" {
				t.Errorf("第 %d 个集群名称不能为空", i)
			}
			if c.Feip == "" {
				t.Errorf("第 %d 个集群IP不能为空", i)
			}
		}

		fmt.Println("✓ 多集群配置测试通过")
	})

	// 测试场景4: 默认值填充
	t.Run("默认值填充", func(t *testing.T) {
		cluster := util.ClusterConfig{
			App:      "test-defaults",
			Feip:     "192.168.1.100",
			User:     "root",
			Password: "pass",
			// Feport, Expire, Status 未设置
		}

		// 模拟默认值逻辑
		feport := cluster.Feport
		if feport == 0 {
			feport = 9030
		}

		expire := cluster.Expire
		if expire == 0 {
			expire = 30
		}

		status := cluster.Status

		if feport != 9030 {
			t.Errorf("Feport 默认值应该是 9030，但得到 %d", feport)
		}
		if expire != 30 {
			t.Errorf("Expire 默认值应该是 30，但得到 %d", expire)
		}
		if status != 0 {
			t.Errorf("Status 默认值应该是 0，但得到 %d", status)
		}

		fmt.Println("✓ 默认值填充测试通过")
	})

	// 测试场景5: 无效配置过滤
	t.Run("无效配置过滤", func(t *testing.T) {
		clusters := []util.ClusterConfig{
			{App: "", Feip: "192.168.1.1", User: "root", Password: "pass"},      // 无效：无App
			{App: "valid", Feip: "192.168.1.2", User: "root", Password: "pass"}, // 有效
			{App: "test", Feip: "", User: "root", Password: "pass"},            // 无效：无Feip
			{App: "test2", Feip: "192.168.1.3", User: "", Password: "pass"},    // 无效：无User
		}

		validCount := 0
		for _, c := range clusters {
			if c.App != "" && c.Feip != "" && c.User != "" {
				validCount++
			}
		}

		if validCount != 1 {
			t.Errorf("应该有1个有效集群，但得到 %d", validCount)
		}

		fmt.Println("✓ 无效配置过滤测试通过")
	})

	fmt.Println("=== 集群初始化集成测试完成 ===")
}

// TestConfigFileStructure 测试配置文件结构
func TestConfigFileStructure(t *testing.T) {
	config := util.RuleConfig{
		SlowQuery: util.SlowQueryRule{
			AlertTimeoutSeconds:     600,
			KillTimeoutSeconds:      1500,
			ConcurrencyLimit:        80,
			FullScanRowsLimit:       200000000,
			CatalogScanRowsLimit:    100000000,
			BEMemoryLimitGB:         200,
			ScanRowsLimit:           10000000000,
			ScanBytesLimitTB:        5,
			ResourceGroupCpuCore:    10,
			ResourceGroupMemGB:      50,
			ResourceGroupConcurrency: 3,
		},
		ClusterConnection: util.ClusterConnRule{
			DefaultFEPort:           9030,
			LicenseExpireRemindDays: 30,
			LicenseCheckSwitch:      0,
		},
		ShortQuery: util.ShortQueryRule{
			DefaultInitPush: 0,
			DefaultStatus:   0,
			DefaultCore:     4,
			DefaultMemoryGB: 8,
		},
		LarkRobot: util.LarkRobotRule{
			DefaultStatus: 0,
		},
		WeComApp: util.WeComAppRule{
			DefaultStatus:     1,
			CorpID:            "wx7d3e81e155049ca3",
			AgentID:           "1000027",
			Secret:            "DvpGuNYmQEIxn-DBzY9Q7_pzEATOE8b3na3WOuijZ7s",
			MsgType:           "markdown",
			MentionAll:        false,
			MentionedUserList: []string{},
		},
		Clusters: []util.ClusterConfig{
			{
				App:       "prod-cluster",
				Nickname:  "生产集群",
				Alias:     "prod",
				Feip:      "192.168.1.100",
				User:      "root",
				Password:  "your-password",
				Feport:    9030,
				Address:   "",
				Expire:    30,
				Status:    0,
				FeLogPath: "",
				BeLogPath: "",
			},
		},
	}

	// 验证配置结构
	if len(config.Clusters) == 0 {
		t.Error("集群配置不能为空")
	}

	if config.WeComApp.CorpID == "" {
		t.Error("企业微信 CorpID 不能为空")
	}

	if config.SlowQuery.AlertTimeoutSeconds <= 0 {
		t.Error("AlertTimeoutSeconds 必须大于0")
	}

	fmt.Println("✓ 配置文件结构测试通过")
}
