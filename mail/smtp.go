package mail

import (
	"crypto/tls"
	"os"
	"strconv"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

func mailClient() (*mail.SMTPClient, error) {
	var (
		EMAIL_USER_SECRET = os.Getenv("EMAIL_USER_SECRET")
		EMAIL_PASSWORD    = os.Getenv("EMAIL_PASSWORD")
		SMTP_HOST         = os.Getenv("SMTP_HOST")
		SMTP_PORT         = os.Getenv("SMTP_PORT")
	)
	server := mail.NewSMTPClient()

	// SMTP Server
	server.Host = SMTP_HOST
	server.Port, _ = strconv.Atoi(SMTP_PORT)
	server.Username = EMAIL_USER_SECRET
	server.Password = EMAIL_PASSWORD
	server.Encryption = mail.EncryptionSTARTTLS

	// Variable to keep alive connection
	server.KeepAlive = false

	// Timeout for connect to SMTP Server
	server.ConnectTimeout = 10 * time.Second

	// Timeout for send the data and wait respond
	server.SendTimeout = 10 * time.Second

	// Set TLSConfig to provide custom TLS configuration. For example,
	// to skip TLS verification (useful for testing):
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// SMTP client
	smtpClient, err := server.Connect()
	if err != nil {
		return nil, err
	}

	return smtpClient, nil
}
