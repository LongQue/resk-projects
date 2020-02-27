package resk

import (
	_ "resk-projects/apis/web"
	_ "resk-projects/core/accounts"
	"resk-projects/infra"
	"resk-projects/infra/base"
)

func init() {
	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DbxDatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})
	infra.Register(&base.IrisServerStarter{})
	infra.Register(&infra.WebApiStarter{})
}