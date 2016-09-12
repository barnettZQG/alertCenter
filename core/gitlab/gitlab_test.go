package gitlab

import (
	"testing"
	"fmt"
	"github.com/astaxie/beego"
)

func Test_GetGitlabAccessToken(t *testing.T) {
	err := beego.LoadAppConfig("ini", "/Users/qwding/gopath/src/alertCenter/conf/app.conf")
	if err!=nil{
		t.Error(err)
		return
	}
	access, err := GetGitlabAccessToken("3ef631e677a61af5bb0cea434746b8f665584858ec45201dd8c9b37e8abcfa5c")
	if err != nil {
		t.Error(err)
	} else {
		fmt.Printf("access token %#v\n", access)
	}
}