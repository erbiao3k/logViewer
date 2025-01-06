package webhook

import (
	"gopkg.in/gomail.v2"
	"log"
	"logViewerServer/setting"
)

func SendMail(mailTo []string, mailAlias, subject, body string) error {

	user := setting.Conf.MailSenderAccount
	password := setting.Conf.MailSenderPassword
	host := setting.Conf.MailConfig.MailServerHost
	port := setting.Conf.MailConfig.MailServerPort

	log.Println(123)
	log.Println(user, password)
	m := gomail.NewMessage()

	m.SetHeader("From", m.FormatAddress(user, mailAlias)) //指定邮件发送方、邮件别名

	m.SetHeader("To", mailTo...) //发送给多个用户

	m.SetHeader("Subject", subject) //设置邮件主题

	m.SetBody("text/html", body) //设置邮件正文

	d := gomail.NewDialer(host, port, user, password)

	err := d.DialAndSend(m)

	return err

}
