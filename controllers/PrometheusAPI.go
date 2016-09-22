package controllers

import (
	"alertCenter/core"
	"alertCenter/core/db"
	"alertCenter/core/service"
	"alertCenter/models"
	"alertCenter/util"
	"encoding/json"

	"github.com/astaxie/beego"
)

type PrometheusAPI struct {
	beego.Controller
}

//ReceivePrometheus 单独验证prometheus
func (e *PrometheusAPI) ReceivePrometheus() {
	ip := e.Ctx.Input.IP()
	configService := &service.GlobalConfigService{
		Session: db.GetMongoSession(),
	}
	if configService.Session != nil {
		defer configService.Session.Close()
	}
	if ok, _ := configService.CheckExist("TrustIP", ip); ok {
		data := e.Ctx.Input.RequestBody
		if data != nil && len(data) > 0 {
			var Alerts []*models.Alert
			err := json.Unmarshal(data, &Alerts)
			if err == nil {
				core.HandleAlerts(Alerts)
				e.Data["json"] = util.GetSuccessJson("receive alert success")
			} else {
				e.Data["json"] = util.GetErrorJson("receive a unknow data")
			}
		}
	} else {
		beego.Debug(ip + " is not trust ip")
		e.Data["json"] = util.GetFailJson("Have no right to access")
	}
	e.ServeJSON()

}
