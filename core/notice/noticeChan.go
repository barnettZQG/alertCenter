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
func GetChannelByMark(mark string) chan *models.Alert {
	if result, ok := NoticChans[mark]; ok {
		return result
	}
	return nil
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