package mysql

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"

	"liveJob/pkg/common/config"
	"liveJob/pkg/tools/errs/errors"
	"liveJob/pkg/zlogger"
)

var (
	LiveDB   *gorm.DB
	XXLJobDB *gorm.DB
)

func InitLiveDB() error {
	var (
		replicasList []gorm.Dialector
		sourcesList  []gorm.Dialector
	)

	logConfig := &gorm.Config{
		SkipDefaultTransaction:                   true, // 禁用默认事务
		DisableForeignKeyConstraintWhenMigrating: true, // 外键
	}

	if config.Config.App.Env == "dev" {
		logConfig.Logger = zlogger.NewDBLog(
			logger.Config{
				SlowThreshold:             time.Second * 2,
				Colorful:                  true,
				IgnoreRecordNotFoundError: true,
				LogLevel:                  logger.Info,
			},
		)
	}
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
			*config.Config.MysqlMaster.Username,
			*config.Config.MysqlMaster.Password,
			(*config.Config.MysqlMaster.Address)[0],
			*config.Config.MysqlMaster.Database),
		DefaultStringSize: 255,
	}), logConfig)

	if err != nil {
		return err
	}

	for _, item := range *config.Config.MysqlSlave.Address {
		dbDsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			*config.Config.MysqlSlave.Username, *config.Config.MysqlSlave.Password, item, *config.Config.MysqlSlave.Database)
		replicasList = append(replicasList, mysql.Open(dbDsn))
	}

	for _, item := range *config.Config.MysqlMaster.Address {
		dbDsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			*config.Config.MysqlMaster.Username, *config.Config.MysqlMaster.Password, item, *config.Config.MysqlMaster.Database)
		sourcesList = append(sourcesList, mysql.Open(dbDsn))
	}

	err = db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  sourcesList,
		Replicas: replicasList,
		Policy:   dbresolver.RandomPolicy{},
		// TraceResolverMode: true,
		TraceResolverMode: false,
	}).SetConnMaxLifetime(time.Second * time.Duration(*config.Config.MysqlMaster.MaxLifeTime)).
		SetMaxOpenConns(*config.Config.MysqlMaster.MaxOpenConn).
		SetMaxIdleConns(*config.Config.MysqlMaster.MaxIdleConn))

	LiveDB = db

	return nil
}

func InitXXLJobDB() error {
	var (
		replicasList []gorm.Dialector
		sourcesList  []gorm.Dialector
	)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
			*config.Config.MysqlMaster.Username,
			*config.Config.MysqlMaster.Password,
			(*config.Config.MysqlMaster.Address)[0],
			*config.Config.MysqlMaster.Database),
		DefaultStringSize: 255,
	}), &gorm.Config{
		SkipDefaultTransaction:                   true, // 禁用默认事务
		DisableForeignKeyConstraintWhenMigrating: true, // 外键
		Logger: zlogger.NewDBLog(
			logger.Config{
				SlowThreshold:             time.Second * 2,
				Colorful:                  true,
				IgnoreRecordNotFoundError: true,
				LogLevel:                  logger.Info,
			}),
	})

	if err != nil {
		return err
	}

	for _, item := range *config.Config.MysqlSlave.Address {
		dbDsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			*config.Config.MysqlSlave.Username, *config.Config.MysqlSlave.Password, item, *config.Config.MysqlSlave.Database)
		replicasList = append(replicasList, mysql.Open(dbDsn))
	}

	for _, item := range *config.Config.MysqlMaster.Address {
		dbDsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			*config.Config.MysqlMaster.Username, *config.Config.MysqlMaster.Password, item, *config.Config.MysqlMaster.Database)
		sourcesList = append(sourcesList, mysql.Open(dbDsn))
	}

	err = db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  sourcesList,
		Replicas: replicasList,
		Policy:   dbresolver.RandomPolicy{},
		// TraceResolverMode: true,
		TraceResolverMode: false,
	}).SetConnMaxLifetime(time.Second * time.Duration(*config.Config.MysqlMaster.MaxLifeTime)).
		SetMaxOpenConns(*config.Config.MysqlMaster.MaxOpenConn).
		SetMaxIdleConns(*config.Config.MysqlMaster.MaxIdleConn))

	XXLJobDB = db

	return nil
}

func CheckErr(err error) error {
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return err
}
