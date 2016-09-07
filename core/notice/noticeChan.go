package notice

import (
	"alertCenter/models"
	"fmt"
)

var NoticChans map[string]chan *models.Alert

func init() {
	NoticChans = make(map[string]chan *models.Alert)
}

//GetChannelByMark 获取发送报警通道
func GetChannelByMark(mark string) (chan *models.Alert, bool) {
	result, ok := NoticChans[mark]
	return result, ok

}

func CreateChanByMark(mark string) error {
	if _, ok := NoticChans[mark]; ok {
		return fmt.Errorf("Channel already exist.")
	} else {
		NoticChans[mark] = make(chan *models.Alert)
	}
	return nil
}

func DeleteChanByMark(mark string) {
	delete(NoticChans, mark)
}
