package db

import (
	"fmt"
	"os"
	"time"

	"github.com/dzrock1989/perkakas/configs"

	postgres2 "go.elastic.co/apm/module/apmgormv2/v2/driver/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConn struct {
	Info       configs.ConnInfo
	SilentMode bool
}

func (conn *DBConn) Open() (db *gorm.DB, err error) {
	if configs.Config.IsUseELK {
		os.Setenv("ELASTIC_APM_ENVIRONMENT", configs.Config.Env)
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable application_name=%s",
			conn.Info.Host,
			conn.Info.User,
			conn.Info.Pass,
			conn.Info.Name,
			conn.Info.Port,
			os.Getenv("ELASTIC_APM_SERVICE_NAME")+"-"+os.Getenv("HOSTNAME"),
		)
		db, err = gorm.Open(
			postgres2.Open(dsn),
			&gorm.Config{},
		)
	} else {
		db, err = gorm.Open(
			postgres.New(postgres.Config{
				DSN: fmt.Sprintf(
					"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
					conn.Info.User,
					conn.Info.Pass,
					conn.Info.Name,
					conn.Info.Host,
					conn.Info.Port,
				),
				PreferSimpleProtocol: true,
			}),
			&gorm.Config{},
		)
	}

	if err != nil {
		return
	}

	db.Logger = db.Logger.LogMode(logger.Silent)
	if !conn.SilentMode {
		db.Logger = db.Logger.LogMode(logger.Info)
	}

	dbGorm, err := db.DB()
	if err != nil {
		return
	}
	dbGorm.SetMaxOpenConns(conn.Info.MaxOpenConn)
	dbGorm.SetMaxIdleConns(conn.Info.MaxIdleConn)
	dbGorm.SetConnMaxLifetime(time.Duration(conn.Info.MaxLifeTime) * time.Minute)

	return
}
