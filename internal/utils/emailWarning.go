package utils

import (
	"go.uber.org/zap"
	"net/smtp"
)

func SendWarningEmail(ctx MyContext, userID, oldIP, newIP string) {
	ctx.Logger.Warn("IP address change detected",
		zap.String("user_id", userID),
		zap.String("email", "to@example.com"),
		zap.String("old_ip", oldIP),
		zap.String("new_ip", newIP),
	)

	from := "from@example.com"
	to := []string{
		"to@example.com",
	}
	msg := []byte("From: " + from + "\r\n" +
		"To: " + to[0] + "\r\n" +
		"Subject: IP address change detected\r\n\r\n" +
		"IP address change detected\n" +
		"Previous IP: " + oldIP + "\n" +
		"Current IP: " + newIP)

	auth := smtp.PlainAuth("", "username", "password", "host")

	_ = smtp.SendMail("example.addr", auth, from, to, msg)

	ctx.Logger.Infof("email sent to %s", to)
}
