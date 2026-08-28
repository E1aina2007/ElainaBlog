package mail

import (
	"ElainaBlog/internal/config"
	"bytes"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"
)

//go:embed templates/verification_code.html
var verificationCodeHTML string

var verificationCodeTmpl = template.Must(template.New("verification_code").Parse(verificationCodeHTML))

// SendVerificationCode 通过 SMTP 发送验证码邮件（HTML 模板渲染）
func SendVerificationCode(smtpCfg config.SmtpConfig, expireSeconds int, to string, code string) error {
	subject := "来自 ElainaBlog 的邮箱验证码"

	htmlBody, err := renderVerificationCode(code, expireSeconds)
	if err != nil {
		return fmt.Errorf("渲染邮件模板失败: %w", err)
	}

	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n%s\r\n%s", smtpCfg.From, to, subject, mime, htmlBody)

	auth := smtp.PlainAuth("", smtpCfg.From, smtpCfg.Verification, smtpCfg.Host)
	addr := fmt.Sprintf("%s:%d", smtpCfg.Host, smtpCfg.Port)

	return smtp.SendMail(addr, auth, smtpCfg.From, []string{to}, []byte(msg))
}

// renderVerificationCode 渲染验证码邮件 HTML
func renderVerificationCode(code string, expireSeconds int) (string, error) {
	var buf bytes.Buffer
	data := struct {
		Code          string
		ExpireSeconds int
	}{
		Code:          code,
		ExpireSeconds: expireSeconds,
	}
	if err := verificationCodeTmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// EmailToAvatarHash 根据邮箱生成头像文件名哈希
// 同一邮箱始终生成相同的哈希值，用于头像文件命名
func EmailToAvatarHash(email string) string {
	email = strings.ToLower(strings.TrimSpace(email))
	hash := sha256.Sum256([]byte(email))
	return hex.EncodeToString(hash[:8]) // 取前8字节，生成16位十六进制字符串
}
