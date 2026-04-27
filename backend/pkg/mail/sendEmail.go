package mail

import (
	"ElainaBlog/config"
	"fmt"
	"net/smtp"
)

// SendVerificationCode 通过 SMTP 发送验证码邮件
func SendVerificationCode(to string, code string) error {
	cfg := config.GlobalConfig.Smtp
	subject := "来自{BlogName}的邮箱验证码"
	body := fmt.Sprintf("您的验证码是%s，请在%d秒内进行验证，若非本人操作，请无视此邮件。", code, config.GlobalConfig.Verification.ExpireTime)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s", cfg.From, to, subject, body)

	auth := smtp.PlainAuth("", cfg.From, cfg.Verification, cfg.Host)
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	return smtp.SendMail(addr, auth, cfg.From, []string{to}, []byte(msg))
}

