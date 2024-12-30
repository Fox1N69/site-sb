package smtp

import (
	"fmt"
	"net/smtp"
)

type SmtpClient interface {
	SendMailRecoveryPassword(email string, recoveryLink string) error
}

type smtpClient struct {
	smtpServer   string
	smtpPort     string
	smtpUsername string
	smtpPassword string
}

func NewSMTPClient(smtpServer, smtpPort, smtpUsername, smtpPassword string) SmtpClient {
	return &smtpClient{
		smtpServer:   smtpServer,
		smtpPort:     smtpPort,
		smtpUsername: smtpUsername,
		smtpPassword: smtpPassword,
	}
}

func (s *smtpClient) SendMailRecoveryPassword(email string, recoveryLink string) error {
	auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpServer)

	to := []string{email}
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: Password Recovery\r\n"+
		"\r\n"+
		"Click %s to recover your password.\r\n", email, recoveryLink))

	err := smtp.SendMail(fmt.Sprintf("%s:%s", s.smtpServer, s.smtpPort), auth, s.smtpUsername, to, msg)
	if err != nil {
		return err
	}

	return nil
}
