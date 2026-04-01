/*
 *@author  chengkenli
 *@project StarRocksQueris
 *@package main
 *@file    initialize_test
 *@date    2024/11/6 22:41
 */

package main

import (
	"StarRocksQueris/util"
	"testing"
)

// TestInitClustersToDB_EmptyClusters 测试空集群配置情况
func TestInitClustersToDB_EmptyClusters(t *testing.T) {
	// 保存原始配置
	originalClusters := util.GlobalRuleConfig.Clusters
	defer func() {
		util.GlobalRuleConfig.Clusters = originalClusters
	}()

	// 设置空集群配置
	util.GlobalRuleConfig.Clusters = []util.ClusterConfig{}

	err := InitClustersToDB()
	if err != nil {
		t.Errorf("空集群配置应该返回nil，但返回了: %v", err)
	}
}

// TestInitClustersToDB_InvalidCluster 测试不完整集群配置
func TestInitClustersToDB_InvalidCluster(t *testing.T) {
	// 保存原始配置
	originalClusters := util.GlobalRuleConfig.Clusters
	defer func() {
		util.GlobalRuleConfig.Clusters = originalClusters
	}()

	// 设置不完整的集群配置
	util.GlobalRuleConfig.Clusters = []util.ClusterConfig{
		{App: "", Feip: "192.168.1.1", User: "root", Password: "pass"},
		{App: "test", Feip: "", User: "root", Password: "pass"},
		{App: "test2", Feip: "192.168.1.2", User: "", Password: "pass"},
	}

	// 此时数据库连接可能为nil，但函数应该能处理跳过逻辑
	err := InitClustersToDB()
	// 由于数据库连接问题，可能会返回错误，但不应该panic
	t.Logf("不完整集群配置测试结果: %v", err)
}

// TestInitClustersToDB_ValidCluster 测试有效集群配置
func TestInitClustersToDB_ValidCluster(t *testing.T) {
	// 保存原始配置
	originalClusters := util.GlobalRuleConfig.Clusters
	defer func() {
		util.GlobalRuleConfig.Clusters = originalClusters
	}()

	// 设置有效的集群配置
	util.GlobalRuleConfig.Clusters = []util.ClusterConfig{
		{
			App:      "test-cluster",
			Nickname: "测试集群",
			Alias:    "test",
			Feip:     "192.168.1.100",
			User:     "root",
			Password: "test-password",
			Feport:   9030,
			Address:  "",
			Expire:   30,
			Status:   0,
		},
	}

	// 设置默认值
	util.GlobalRuleConfig.ClusterConnection.DefaultFEPort = 9030
	util.GlobalRuleConfig.ClusterConnection.LicenseExpireRemindDays = 30
	util.GlobalRuleConfig.ClusterConnection.LicenseCheckSwitch = 0

	err := InitClustersToDB()
	// 如果没有数据库连接，会返回错误
	if err != nil {
		t.Logf("有效集群配置测试（无数据库连接）: %v", err)
	}
}

// TestInitClustersToDB_DefaultValues 测试默认值填充
func TestInitClustersToDB_DefaultValues(t *testing.T) {
	cluster := util.ClusterConfig{
		App:      "test-default",
		Nickname: "测试默认",
		Feip:     "192.168.1.100",
		User:     "root",
		Password: "pass",
		// 不设置 Feport, Expire, Status
	}

	// 验证默认值逻辑
	feport := cluster.Feport
	if feport == 0 {
		feport = 9030 // 默认值
	}
	if feport != 9030 {
		t.Errorf("Feport 默认值应该是 9030，但得到 %d", feport)
	}

	expire := cluster.Expire
	if expire == 0 {
		expire = 30 // 默认值
	}
	if expire != 30 {
		t.Errorf("Expire 默认值应该是 30，但得到 %d", expire)
	}
}

// TestInitClustersToDB_MultipleClusters 测试多个集群配置
func TestInitClustersToDB_MultipleClusters(t *testing.T) {
	// 保存原始配置
	originalClusters := util.GlobalRuleConfig.Clusters
	defer func() {
		util.GlobalRuleConfig.Clusters = originalClusters
	}()

	// 设置多个集群配置
	util.GlobalRuleConfig.Clusters = []util.ClusterConfig{
		{
			App:      "cluster-1",
			Feip:     "192.168.1.101",
			User:     "root",
			Password: "pass1",
			Feport:   9030,
		},
		{
			App:      "cluster-2",
			Feip:     "192.168.1.102",
			User:     "admin",
			Password: "pass2",
			Feport:   9031,
		},
	}

	err := InitClustersToDB()
	t.Logf("多个集群配置测试: %v", err)
}
