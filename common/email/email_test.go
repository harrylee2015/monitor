package email

import (
	"testing"

	"github.com/harrylee2015/monitor/model"
)

func TestQQEmail_Send(t *testing.T) {
	e := &model.Email{
		FromMail: "harrylee2015@qq.com",
		PassWd:   "vzbvfmsipfkrhfce",
		Host:     "smtp.qq.com",
		Port:     465,
		ToMail:   "18761806026@163.com",
		Subject:  "test",
		Body:     "testing",
	}
	err := SendMail(e)
	if err != nil {
		t.Error(err)
	}
}
