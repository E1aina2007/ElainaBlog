package mail

import (
	"ElainaBlog/config"
	"fmt"
	"net/smtp"
)

// SendVerificationCode 通过 SMTP 发送验证码邮件（HTML 格式）
func SendVerificationCode(to string, code string) error {
	cfg := config.GlobalConfig.Smtp
	expireSeconds := config.GlobalConfig.Verification.ExpireTime
	subject := "来自 ElainaBlog 的邮箱验证码"

	htmlBody := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="margin:0;padding:0;background:#f4f7fa;font-family:'Microsoft YaHei','PingFang SC',sans-serif;">
  <div style="max-width:480px;margin:40px auto;background:#ffffff;border-radius:12px;box-shadow:0 2px 12px rgba(0,0,0,0.08);overflow:hidden;">
    <div style="background:linear-gradient(135deg,#7ED7C1,#5BC4AD);padding:28px 32px;text-align:center;">
      <h1 style="color:#fff;margin:0;font-size:20px;font-weight:600;">ElainaBlog</h1>
      <p style="color:rgba(255,255,255,0.85);margin:8px 0 0;font-size:13px;">邮箱验证码</p>
    </div>
    <div style="padding:32px;">
      <p style="color:#333;font-size:15px;line-height:1.6;margin:0 0 20px;">你好，你正在注册 ElainaBlog 账号，请使用以下验证码完成验证：</p>
      <div style="text-align:center;margin:24px 0;">
        <span style="display:inline-block;background:#f0f7f4;border:2px dashed #7ED7C1;border-radius:8px;padding:14px 32px;font-size:32px;font-weight:700;color:#2d8a6e;letter-spacing:6px;">%s</span>
      </div>
      <p style="color:#888;font-size:13px;line-height:1.6;margin:0 0 8px;">验证码有效期为 <strong style="color:#555;">%d 秒</strong>，请尽快完成验证。</p>
      <p style="color:#aaa;font-size:12px;line-height:1.6;margin:0;">若非本人操作，请忽略此邮件，不会对你的账号产生任何影响。</p>
    </div>
    <div style="background:#f9fafb;padding:16px 32px;text-align:center;border-top:1px solid #eee;">
      <p style="color:#bbb;font-size:11px;margin:0;">此邮件由系统自动发送，请勿直接回复</p>
    </div>
  </div>
</body>
</html>`, code, expireSeconds)

	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n%s\r\n%s", cfg.From, to, subject, mime, htmlBody)

	auth := smtp.PlainAuth("", cfg.From, cfg.Verification, cfg.Host)
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	return smtp.SendMail(addr, auth, cfg.From, []string{to}, []byte(msg))
}
