/**
*FileName: utils
*Create on 2018/12/6 上午6:00
*Create by mok
*/

package utils

import (
	"github.com/go-gomail/gomail"
	"html/template"
	"bytes"
)

type EmailConfig struct {
	From []string
	To []string
	Subject string
	Tp *template.Template
	Messages map[string]interface{}
}

func SendEmail(c *EmailConfig)error{
	m := gomail.NewMessage()
	m.SetAddressHeader("From",c.From[0],c.From[1])
	m.SetHeader("To",c.To...)
	m.SetHeader("Subject",c.Subject)
	buffer := &bytes.Buffer{}
	err := c.Tp.Execute(buffer,c.Messages)
	if err != nil{
		return err
	}
	m.SetBody("text/html",buffer.String())
	d := gomail.NewDialer("smtp.qq.com",465,"1005914310@qq.com","aboyddtczgrrbeac")
	err  = d.DialAndSend(m)
	return err
}
