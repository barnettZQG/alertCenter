package controllers

import (
	"encoding/json"
	"strconv"
	"time"
	"alertCenter/core"
	"alertCenter/core/db"
	"alertCenter/core/service"
	"alertCenter/models"
	"alertCenter/util"
	"github.com/astaxie/beego"
	"alertCenter/core/gitlab"
)

type APIController struct {
	//beego.Controller
	APIBaseController
	session      *db.MongoSession
	alertService *service.AlertService
	teamServcie  *service.TeamService
}

func (e *APIController) ReceiveAlert() {

	data := e.Ctx.Input.RequestBody
	if data != nil && len(data) > 0 {
		var AlertMessage *models.AlertReceive = &models.AlertReceive{}
		err := json.Unmarshal(data, AlertMessage)
		if err == nil {
			beego.Info("get a alert message,receiver:" + AlertMessage.Receiver)
			go core.HandleMessage(AlertMessage)
			e.Data["json"] = util.GetSuccessJson("receive alert success")
		} else {
			beego.Error("receive a unknow data")
			//util.Info(string(data))
		}
	}
	e.ServeJSON()
}
func (e *APIController) Receive() {
	data := e.Ctx.Input.RequestBody
	if data != nil && len(data) > 0 {
		var Alerts []*models.Alert = make([]*models.Alert, 0)
		err := json.Unmarshal(data, &Alerts)
		if err == nil {
			core.HandleAlerts(Alerts)
			e.Data["json"] = util.GetSuccessJson("receive alert success")
		} else {
			beego.Error("receive a unknow data")
			//util.Info(string(data))
		}
	}
	e.ServeJSON()
}

// func (e *APIController) AddTag() {
// 	we := &core.notice.WeAlertSend{}
// 	if ok := we.GetAllTags(); ok {
// 		e.Data["json"] = util.GetSuccessJson("get weiTag success")
// 	} else {
// 		e.Data["json"] = util.GetFailJson("get weiTag faild")
// 	}
// 	e.ServeJSON()
// }

func (e *APIController) HandleAlert() {
	ID := e.GetString(":ID")
	Type := e.GetString(":type")
	message := e.GetString("message")
	if len(ID) == 0 || len(Type) == 0 {
		e.Data["json"] = util.GetErrorJson("参数格式错误")
	} else {
		session := db.GetMongoSession()
		defer session.Close()
		alertService := service.GetAlertService(session)
		alert := alertService.FindByID(ID)
		if alert == nil {
			e.Data["json"] = util.GetFailJson("报警信息不存在，id信息错误")
		} else {
			if Type == "handle" {
				alert.IsHandle = 1
			} else if Type == "miss" {
				alert.IsHandle = -1
			}
			alert.HandleDate = time.Now()
			alert.HandleMessage = message
			if ok := alertService.Update(alert); ok {
				e.Data["json"] = util.GetSuccessJson("登记成功")
			} else {
				e.Data["json"] = util.GetFailJson("登记失败")
			}
		}
	}
	e.ServeJSON()
}

func (e *APIController) GetAlerts() {
	pageSizeStr := e.Ctx.Request.FormValue("pageSize")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 20
	}

	pageStr := e.Ctx.Request.FormValue("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	receiver := e.APIBaseController.Username

	// if admin. Should show all the alerts.
	user, err := gitlab.GetUserByUsername(receiver)
	if err != nil {
		beego.Error(err)
	} else if user.IsAdmin {
		receiver = "all"
	}

	e.session = db.GetMongoSession()
	if e.session == nil {
		e.Data["json"] = util.GetFailJson("get database session faild.")
		goto over
	} else {
		e.alertService = service.GetAlertService(e.session)
		defer e.session.Close()
		if len(receiver) != 0 && receiver != "all" {
			alerts := e.alertService.FindByUser(receiver, pageSize, page)
			beego.Info("Get", len(alerts), " alerts")
			if alerts == nil {
				e.Data["json"] = util.GetFailJson("get database collection faild or receiver is error ")
				goto over
			} else {
				e.Data["json"] = util.GetSuccessReJson(alerts)
				goto over
			}
		} else if receiver == "all" {
			alerts := e.alertService.FindAll(pageSize, page)
			if alerts == nil {
				e.Data["json"] = util.GetFailJson("get database collection faild")
				goto over
			} else {
				e.Data["json"] = util.GetSuccessReJson(alerts)
				goto over
			}
		} else {
			e.Data["json"] = util.GetErrorJson("api use error,please provide receiver")
			goto over
		}
	}
	over:
	e.ServeJSON()
}
