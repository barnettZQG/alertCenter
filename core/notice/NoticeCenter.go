package notice

import "github.com/astaxie/beego"
import "strings"

var (
	cacheServer map[string]NoticeServer
)

//StartCenter 初始化并开始发送通知操作
func StartCenter() (err error) {
	beego.Info("Starting NoticeServer is begin")
	server := beego.AppConfig.String("NoticeServer")
	servers := strings.Split(server, ",")
	cacheServer = make(map[string]NoticeServer, 0)
	for _, s := range servers {
		noticeServer := GetNoticeServer(s)
		if noticeServer != nil {
			err = noticeServer.StartWork()
			if err != nil {
				return
			}
			cacheServer[s] = noticeServer
		}
	}
	beego.Info("Starting NoticeServer success")

	return
}


