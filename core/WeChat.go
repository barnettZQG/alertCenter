package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/barnettZQG/alertCenter/models"
	"github.com/barnettZQG/alertCenter/util"
	"gopkg.in/mgo.v2/bson"
)

type WeAlertSend struct {
	weChan   chan *WeMessage
	stopChan chan bool
}

var WeTag map[string]*models.WeTag

func (e *WeAlertSend) SendAlert(alert *models.Alert) {
	url := beego.AppConfig.String("url")
	me := e.GetWeAlertByTag(alert.Receiver.Name)
	message := ``
	message += string(alert.Annotations["description"]) + "\n\t"
	message += `------------------\n\r`
	message += `[<a href=\"` + url + `/alertList?receiver＝` + alert.Receiver.Name + `\">点击查看详情</a>]`
	me = strings.Replace(me, "CONTENT", message, -1)
	e.weChan <- &WeMessage{
		message:  me,
		alert:    alert,
		errCount: 0,
	}
}

func (e *WeAlertSend) SendWeChatMessage(mestr string) error {
	//util.Debug("send weiChat message :" + mestr)
	MessageURI := beego.AppConfig.String("weURI") + "/cgi-bin/message/send?access_token=ACCESS_TOKEN"
	body := bytes.NewBufferString(mestr) //.NewReader(me)
	client := &http.Client{}
	req, err := http.NewRequest("POST", MessageURI, body)
	if err != nil {
		util.Error("create wechat request faild." + err.Error())
		return err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("token", beego.AppConfig.String("weToken"))
	resp, err := client.Do(req)
	if err != nil {
		util.Error("Send weChat message error," + err.Error())
		return err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.Error("Send weChat message error," + err.Error())
		return err
	} else {
		util.Info("Send WeChat result feedback：" + string(content))
	}
	return nil
}
func (e *WeAlertSend) GetWeAlertByTag(tags string) string {
	AgentId := beego.AppConfig.String("weiAgentId")
	tagId, err := e.GetTagIdByTag(tags)
	if err == nil {
		return `{"totag": " ` + strconv.Itoa(tagId) + ` ","msgtype": "text","agentid": ` + AgentId + `,"text": {"content": "CONTENT"},"safe":0}`
	}
	return ""
}
func (e *WeAlertSend) GetTagIdByTag(tag string) (int, error) {
	session := GetMongoSession()
	defer session.Close()
	coll := session.GetCollection("WeiTag")
	if coll == nil {
		return 0, errors.New("get collection WeiTag faild")
	}
	weTag := &models.WeTag{}
	err := coll.Find(bson.M{"tagname": tag}).Select(nil).One(&weTag)
	if err != nil {
		util.Error("get weiTag by name error." + err.Error())
		return 0, err
	}
	return weTag.TagId, nil
}
func (e *WeAlertSend) GetAllTags() bool {
	TagListURI := beego.AppConfig.String("weURI") + "/cgi-bin/tag/list?access_token=ACCESS_TOKEN"
	client := &http.Client{}
	req, err := http.NewRequest("GET", TagListURI, nil)
	if err != nil {
		util.Error("create get taglist request faild." + err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", beego.AppConfig.String("weToken"))
	resp, err := client.Do(req)
	if err != nil {
		util.Error("GET taglist error," + err.Error())
		return false
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.Error("GET taglist error," + err.Error())
		return false
	}
	WeiTagResult := &models.WeiTagResult{}
	err = json.Unmarshal(content, WeiTagResult)
	if err != nil {
		util.Error("Parse the taglist data error ." + err.Error())
		return false
	}
	for _, tag := range WeiTagResult.TagList {
		WeTag[tag.TagName] = &tag
	}

	session := GetMongoSession()
	defer session.Close()
	if session == nil {
		util.Error("Get mongo session error." + err.Error())
		return false
	}
	if ok := session.RemoveAll("WeTag"); ok {

		var data []interface{}
		for _, tag := range WeiTagResult.TagList {
			data = append(data, tag)
		}
		util.Debug("Got the wetag number is " + strconv.Itoa(len(data)))
		return session.Insert("WeTag", data...)
	}
	return false
}

type WeMessage struct {
	alert    *models.Alert
	errCount int
	message  string
}

func (e *WeAlertSend) StartWork() error {
	beego.Info("Wexin sender init begin")
	defer beego.Info("Wexin sender init over")

	weCount, err := beego.AppConfig.Int("weCount")
	if err != nil {
		beego.Error("weCount's type is not int ." + err.Error())
		return err
	}
	weReCount, err := beego.AppConfig.Int("weReCount")
	if err != nil {
		beego.Error("weReCount's type is not int ." + err.Error())
		return err
	}
	if WeTag == nil {
		WeTag = make(map[string]*models.WeTag, 0)
		if ok := e.GetAllTags(); !ok {
			return errors.New("get all weTags faild")
		}
	}
	if e.weChan == nil {
		e.weChan = make(chan *WeMessage, weCount)
	}
	if e.stopChan == nil {
		e.stopChan = make(chan bool)
	}
	go func() {

		util.Info("wexin work start success")
		for {
			select {
			case m, ok := <-e.weChan:
				if !ok {
					return
				}
				if err := e.SendWeChatMessage(m.message); err != nil {
					m.errCount++
					if m.errCount < weReCount {
						//5秒后重试
						go func(m *WeMessage) {
							time.Sleep(time.Second * 5)
							e.weChan <- m
						}(m)
					}
				}
			case stop := <-e.stopChan:
				if stop {
					goto exit
				}
			}
		}
	exit:
		util.Info("mail work stop success")
	}()

	return nil
}
func (e *WeAlertSend) StopWork() {
	if e.stopChan != nil {
		e.stopChan <- true
		close(e.stopChan)
	}
	if e.weChan != nil {
		close(e.weChan)
	}
}

func GetWeTagByName(name string) *models.WeTag {
	if WeTag != nil {
		return WeTag[name]
	}
	return nil
}
