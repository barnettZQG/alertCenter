package notice

import (
	"alertCenter/models"
	"fmt"
	"testing"
)

// func Test_SendMail(t *testing.T) {
// 	fmt.Println("in SendMail_Test")
// 	server := &MailNoticeServer{}
// 	message := server.GetMessage("hhhhh", "hhhhhh", "zengqingguo@goyoo.com")
// 	server.SendMail(message.message)
// }

func Test_GetBody(t *testing.T) {
	server := &MailNoticeServer{}
	fmt.Println(server.GetBody(&models.Alert{}))
}
