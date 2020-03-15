package jobs

import (
	log "github.com/sirupsen/logrus"
	"resk-projects/infra"
	"time"
)

type RefundExpiredJobStarter struct {
	infra.BaseStarter
	ticker *time.Ticker
}

func (r *RefundExpiredJobStarter) Init(ctx infra.StarterContext) {
	d := ctx.Props().GetDurationDefault("jobs.refund.interval", time.Minute)
	r.ticker = time.NewTicker(d)
}

func (r *RefundExpiredJobStarter) Start(ctx infra.StarterContext) {
	log.Info("进入定时任务")
	go func() {
		for {
			c := <-r.ticker.C
			log.Info("过期红包退款开始...", c)

		}
	}()

}

func (r *RefundExpiredJobStarter) Stop(ctx infra.StarterContext) {
	r.ticker.Stop()
}

//这里把测试和debug留给同学们
