package notice

import "alertCenter/models"

//NoticeServer 发送通知插件接口
type NoticeServer interface {
	StartWork() error
	StopWork() error
	SendAlert(alert *models.Alert) error
}

//GetNoticeServer 获取通知插件
func GetNoticeServer(name string) NoticeServer {
	if name == "mail" {
		return &MailNoticeServer{}
	}
	if name == "wexin" {
		return &WeNoticeServer{}
	}
	return nil
}
