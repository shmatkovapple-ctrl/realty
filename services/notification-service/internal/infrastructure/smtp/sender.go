package smtp

import (
"fmt"
"net/smtp"
)

type EmailSender struct {
host     string
port     string
username string
password string
from     string
}

func NewEmailSender(host, port, username, password string) *EmailSender {
return &EmailSender{
host:     host,
port:     port,
username: username,
password: password,
from:     username,
}
}

func (s *EmailSender) Send(to, subject, body string) error {
if s.username == "" {
return nil
}

auth := smtp.PlainAuth("", s.username, s.password, s.host)

msg := fmt.Sprintf(
"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
s.from, to, subject, body,
)

addr := fmt.Sprintf("%s:%s", s.host, s.port)
if err := smtp.SendMail(addr, auth, s.from, []string{to}, []byte(msg)); err != nil {
return fmt.Errorf("отправка email: %w", err)
}

return nil
}
