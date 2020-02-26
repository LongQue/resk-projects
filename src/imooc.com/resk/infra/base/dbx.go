package base

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"github.com/tietang/props/kvs"
	"resk-projects/infra"
)

//dbx 数据库实例
//dbx 数据库实例
var database *dbx.Database

func DbxDatabase() *dbx.Database {
	return database
}

//dbx数据库starter，并设置全局
type DbxDatabaseStarter struct {
	infra.BaseStarter
}

func (s *DbxDatabaseStarter) Setup(ctx infra.StarterContext) {
	conf := ctx.Props()
	//数据库配置
	settings := dbx.Settings{}
	err := kvs.Unmarshal(conf, &settings, "mysql")
	if err != nil {
		panic(err)
	}
	logrus.Infof("%+v\n",settings)
	logrus.Info("mysql.conn url:", settings.ShortDataSourceName())
	db, err := dbx.Open(settings)
	if err != nil {
		panic(err)
	}
	logrus.Info(db.Ping())
	database = db
}

