package testx

import (
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
	"resk-projects/infra"
	"resk-projects/infra/base"
)

func init()  {
	file:=kvs.GetCurrentFilePath("../brun/test/config.ini",1)
	//加载和解析配置文件
	conf:=ini.NewIniFileCompositeConfigSource(file)
	base.InitLog(conf)

	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DbxDatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})

	app:=infra.New(conf)
	app.Start()
}