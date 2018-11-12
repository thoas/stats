package main

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/thoas/stats"
)

var Stats = stats.New()

type StatsController struct {
	beego.Controller
}

func (this *StatsController) Get() {
	this.Data["json"] = Stats.Data()
	this.Ctx.Output.SetStatus(200)
	this.ServeJson()
}

func main() {
	beego.InsertFilter("*", beego.BeforeRouter, func(ctx *context.Context) {
		startTime := time.Now()
		ctx.Input.SetData("stats_timer", startTime)
	})
	beego.InsertFilter("*", beego.FinishRouter, func(ctx *context.Context) {
		Stats.End(ctx.Input.GetData("stats_timer").(time.Time), stats.WithStatusCode(ctx.Output.Status))
	})

	beego.Router("/stats", &StatsController{})

	beego.Run()
}
