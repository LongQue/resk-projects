package resk

import (
	"resk-projects/apis/gorpc"
	_ "resk-projects/apis/gorpc"
	_ "resk-projects/apis/web"
	_ "resk-projects/core/accounts"
	_ "resk-projects/core/envelopes"

	"resk-projects/infra"
	"resk-projects/infra/base"
)

func init() {
	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DbxDatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})
	infra.Register(&base.GoRPCStarter{})
	infra.Register(&gorpc.GoRpcApiStarter{})
	infra.Register(&base.IrisServerStarter{})
	infra.Register(&infra.WebApiStarter{})
}
