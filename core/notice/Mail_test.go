package notice

import (
	"fmt"
	"testing"
)

func Test_SendMail(t *testing.T) {
	fmt.Println("in SendMail_Test")
	server := &MailNoticeServer{}
	message := server.GetMessage("hhhhh", "hhhhhh", "zengqingguo@goyoo.com")
	server.SendMail(message.message)
}
