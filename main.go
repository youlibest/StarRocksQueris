/*
 *@author  chengkenli
 *@project StarRocksQueris
 *@package StarRocksQueris
 *@file    main
 *@date    2024/8/7 14:42
 */

package main

import (
	"StarRocksQueris/logs"
	"StarRocksQueris/run"
	"StarRocksQueris/util"
	"fmt"
)

func printStarRocks() {
	fmt.Println(`            ____  _             ____            _       
           / ___|| |_ __ _ _ __|  _ \ ___   ___| | _____ 
           \___ \| __/ __| ___| |_) / _ \ / __| |/ / __|
            ___) | || (_| | |  |  _ < (_) | (__|   <\__ \ 
           |____/ \__\____|_|  |_| \_\___/ \___|_|\_\___/
                                                            c.k(v1.4)`)
}

func main() {
	// 启动日志服务
	go logs.Logserver()
	// 打印LOGO
	printStarRocks()
	// 加载基础配置 + 规则配置文件（核心：先加载所有配置）
	util.Parms()
	// 初始化MySQL配置库（此时配置已全部加载完成）
	ConfigDB()
	// 启动核心业务逻辑
	run.Run()
}