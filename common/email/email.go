package email

import (
	"github.com/harrylee2015/monitor/model"
	"gopkg.in/gomail.v2"
)

func SendMail(email *model.Email) error {
	//定义邮箱服务器连接信息，如果是阿里邮箱 pass填密码，qq邮箱填授权码
	// mailConn := map[string]string{
	// 	"user": "zhangqiang@xxxx.com",
	// 	"pass": "xxxx",
	// 	"host": "smtp.mxhichina.com",
	// 	"port": "465",
	// }
	m := gomail.NewMessage()
	m.SetHeader("From", "chain33 system monitor"+"<"+email.FromMail+">") //这种方式可以添加别名，即“XD Game”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	m.SetHeader("To", email.ToMail)                                      //发送给多个用户
	m.SetHeader("Subject", email.Subject)                                //设置邮件主题
	m.SetBody("text/html", email.Body)                                   //设置邮件正文

	d := gomail.NewDialer(email.Host, email.Port, email.FromMail, email.PassWd)

	err := d.DialAndSend(m)
	return err

}
