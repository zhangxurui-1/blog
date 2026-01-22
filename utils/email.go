package utils

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"server/global"
	"strconv"
	"strings"

	"github.com/jordan-wright/email"
)

func Email(To, subject, body string) error {
	to := strings.Split(To, ";")
	return send(to, subject, body)
}

// send sends an email
func send(to []string, subject, body string) error {
	emailGfg := global.Config.Email

	from := emailGfg.From
	nickname := emailGfg.Nickname
	secret := emailGfg.Secret
	host := emailGfg.Host
	port := emailGfg.Port
	isSSL := emailGfg.IsSSL

	// 使用 PlainAuth 创建认证信息
	auth := smtp.PlainAuth("", from, secret, host)
	// 创建新的电子邮件对象
	e := email.NewEmail()
	if nickname != "" {
		e.From = fmt.Sprintf("%s <%s>", nickname, from)
	} else {
		e.From = from
	}

	// 设置收件人、主题和邮件内容
	e.To = to
	e.Subject = subject
	e.HTML = []byte(body)

	var err error
	// 构建邮件服务器的地址，格式为 host:port
	hostAddr := host + ":" + strconv.Itoa(port)

	if isSSL {
		err = e.SendWithTLS(hostAddr, auth, &tls.Config{ServerName: host})
	} else {
		err = e.Send(hostAddr, auth)
	}
	return err
}
