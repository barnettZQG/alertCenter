package notice

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"alertCenter/core/db"
	"alertCenter/models"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

type WeNoticeServer struct {
	weChan   chan *WeMessage
	stopChan chan bool
}

var WeTag map[string]*models.WeTag

func (e *WeNoticeServer) SendAlert(alert *models.Alert) error {
	url := beego.AppConfig.String("url")
	me := e.GetWeAlertByTag(alert.Receiver.Name)
	message := ``
	message += string(alert.Annotations.LabelSet["description"]) + "\n\t"
	message += `------------------\n\r`
	message += `[<a href=\"` + url + `/alertList?receiver＝` + alert.Receiver.Name + `\">点击查看详情</a>]`
	me = strings.Replace(me, "CONTENT", message, -1)
	e.weChan <- &WeMessage{
		message:  me,
		alert:    alert,
		errCount: 0,
	}
	return nil
}
func (e *WeNoticeServer) StartWork() error {
	beego.Info("Wexin notice server init begin")
	defer beego.Info("Wexin notice server init over")

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
		beego.Info("mail work stop success")
	}()
	beego.Info("wexin notice server start success")
	return nil
}
func (e *WeNoticeServer) StopWork() error {
	if e.stopChan != nil {
		e.stopChan <- true
		close(e.stopChan)
	}
	if e.weChan != nil {
		close(e.weChan)
	}
	return nil
}

func (e *WeNoticeServer) SendWeChatMessage(mestr string) error {
	//util.Debug("send weiChat message :" + mestr)
	MessageURI := beego.AppConfig.String("weURI") + "/cgi-bin/message/send?access_token=ACCESS_TOKEN"
	body := bytes.NewBufferString(mestr) //.NewReader(me)
	client := &http.Client{}
	req, err := http.NewRequest("POST", MessageURI, body)
	if err != nil {
		beego.Error("create wechat request faild." + err.Error())
		return err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("token", beego.AppConfig.String("weToken"))
	resp, err := client.Do(req)
	if err != nil {
		beego.Error("Send weChat message error," + err.Error())
		return err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		beego.Error("Send weChat message error," + err.Error())
		return err
	} else {
		beego.Info("Send WeChat result feedback：" + string(content))
	}
	return nil
}
func (e *WeNoticeServer) GetWeAlertByTag(tags string) string {
	AgentId := beego.AppConfig.String("weiAgentId")
	tagId, err := e.GetTagIdByTag(tags)
	if err == nil {
		return `{"totag": " ` + strconv.Itoa(tagId) + ` ","msgtype": "text","agentid": ` + AgentId + `,"text": {"content": "CONTENT"},"safe":0}`
	}
	return ""
}
func (e *WeNoticeServer) GetTagIdByTag(tag string) (int, error) {
	session := db.GetMongoSession()
	defer session.Close()
	coll := session.GetCollection("WeiTag")
	if coll == nil {
		return 0, errors.New("get collection WeiTag faild")
	}
	weTag := &models.WeTag{}
	err := coll.Find(bson.M{"tagname": tag}).Select(nil).One(&weTag)
	if err != nil {
		beego.Error("get weiTag by name error." + err.Error())
		return 0, err
	}
	return weTag.TagId, nil
}
func (e *WeNoticeServer) GetAllTags() bool {
	TagListURI := beego.AppConfig.String("weURI") + "/cgi-bin/tag/list?access_token=ACCESS_TOKEN"
	client := &http.Client{}
	req, err := http.NewRequest("GET", TagListURI, nil)
	if err != nil {
		beego.Error("Get mongo session error." + err.Error())
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", beego.AppConfig.String("weToken"))
	resp, err := client.Do(req)
	if err != nil {
		beego.Error("GET taglist error," + err.Error())
		return false
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		beego.Error("GET taglist error," + err.Error())
		return false
	}
	WeiTagResult := &models.WeiTagResult{}
	err = json.Unmarshal(content, WeiTagResult)
	if err != nil {
		beego.Error("Parse the taglist data error ." + err.Error())
		return false
	}
	for _, tag := range WeiTagResult.TagList {
		WeTag[tag.TagName] = &tag
	}

	session := db.GetMongoSession()
	defer session.Close()
	if ok := session.RemoveAll("WeTag"); ok {

		var data []interface{}
		for _, tag := range WeiTagResult.TagList {
			data = append(data, tag)
		}
		beego.Debug("Got the wetag number is " + strconv.Itoa(len(data)))
		return session.Insert("WeTag", data...)
	}
	return false
}

type WeMessage struct {
	alert    *models.Alert
	errCount int
	message  string
}

func GetWeTagByName(name string) *models.WeTag {
	if WeTag != nil {
		return WeTag[name]
	}
	return nil
}
