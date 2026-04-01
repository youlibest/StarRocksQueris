package conn

import (
	"StarRocksQueris/tools"
	"StarRocksQueris/util"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

func StarRocks(app string) (*gorm.DB, error) {
	var avg tools.SrAvgs
	// connect
	if len(tools.UniqueMaps(util.ConnectBody)) == 0 {
		return nil, errors.New("config db is null")
	}
	for _, m := range tools.UniqueMaps(util.ConnectBody) {
		appVal, ok1 := m["app"].(string)
		if !ok1 || appVal != app {
			continue
		}
		feip, ok2 := m["feip"].(string)
		feport, ok3 := m["feport"].(int)
		if !ok3 {
			// 尝试 int32 类型
			if feport32, ok := m["feport"].(int32); ok {
				feport = int(feport32)
				ok3 = true
			}
		}
		user, ok4 := m["user"].(string)
		password, ok5 := m["password"].(string)

		if ok2 && ok3 && ok4 && ok5 {
			avg = tools.SrAvgs{
				Host: feip,
				Port: feport,
				User: user,
				Pass: password,
			}
			break
		}
	}

	if avg.Host == "" {
		return nil, errors.New("avg is null")
	}

	newLogger := logger.New(nil,
		logger.Config{
			SlowThreshold: time.Second * 1000, // 控制慢SQL阈值
		},
	)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/information_schema?parseTime=true&charset=utf8mb4&loc=Local",
		avg.User,
		avg.Pass,
		avg.Host,
		avg.Port,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: newLogger,
	})
	if err != nil {
		util.Loggrs.Error(app, " ->: ", err.Error())
	}
	return db, err
}

func StarRocksApp(app, host string) (*gorm.DB, error) {
	var avg tools.SrAvgs

	if len(host) == 0 {
		return nil, errors.New("登录信息有误，或为空。")
	}
	for _, m := range tools.UniqueMaps(util.ConnectBody) {
		appVal, ok1 := m["app"].(string)
		if !ok1 || appVal != app {
			continue
		}
		feport, ok2 := m["feport"].(int)
		if !ok2 {
			// 尝试 int32 类型
			if feport32, ok := m["feport"].(int32); ok {
				feport = int(feport32)
				ok2 = true
			}
		}
		user, ok3 := m["user"].(string)
		password, ok4 := m["password"].(string)

		if ok2 && ok3 && ok4 {
			avg = tools.SrAvgs{
				Host: host, // 使用传入的 host 参数
				Port: feport,
				User: user,
				Pass: password,
			}
			break
		}
	}

	if avg.Host == "" {
		return nil, errors.New("avg is null")
	}

	newLogger := logger.New(nil,
		logger.Config{
			SlowThreshold: time.Second * 1000, // 控制慢SQL阈值
		},
	)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/information_schema?charset=utf8mb4&parseTime=True&loc=Local", avg.User, avg.Pass, host, avg.Port)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名

		},
		Logger: newLogger,
	})
	if err != nil {
		util.Loggrs.Error(app, " ->: ", err.Error())
		return db, err
	}
	return db, err
}

func StarRocksItem(item *tools.SrAvgs) (*gorm.DB, error) {
	newLogger := logger.New(nil,
		logger.Config{
			SlowThreshold: time.Second * 1000, // 控制慢SQL阈值
		},
	)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/information_schema?charset=utf8mb4&parseTime=True&loc=Local", item.User, item.Pass, item.Host, 9030)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名

		},
		Logger: newLogger,
	})
	if err != nil {
		util.Loggrs.Error(item.Host, " ->: ", err.Error())
		return db, err
	}
	return db, err
}

func StarRocksShort(item *util.ReData) (*gorm.DB, error) {
	// connect
	if len(tools.UniqueMaps(util.ConnectBody)) == 0 {
		return nil, errors.New("config db is null")
	}
	var host string
	var port int
	for _, m := range tools.UniqueMaps(util.ConnectBody) {
		if m["app"].(string) == item.App {
			host = m["feip"].(string)
			port = int(m["feport"].(int32))
		}
	}
	if host == "" {
		return nil, errors.New("avg is null")
	}
	//password, err := util.AesDecrypt2(item.Password, util.ENCKEY)
	//if err != nil {
	//	util.Loggrs.Error(fmt.Sprintf("%s", err.Error()))
	//}
	newLogger := logger.New(nil,
		logger.Config{
			SlowThreshold: time.Second * 1000, // 控制慢SQL阈值
		},
	)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/information_schema?parseTime=true&charset=utf8mb4&loc=Local",
		item.Username,
		item.Password,
		host,
		port,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: newLogger,
	})
	if err != nil {
		util.Loggrs.Error(item.App, " ->: ", err.Error())
	}
	return db, err
}
