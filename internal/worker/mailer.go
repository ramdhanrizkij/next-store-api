package worker

import (
	"context"
	"log"
)

type Mailer interface {
	SendVerificationEmail(ctx context.Context, toEmail, toName, token, verificationURL string) error
}

type LogMailer struct{}

func NewLogMailer() Mailer {
	return &LogMailer{}
}

func (m *LogMailer) SendVerificationEmail(ctx context.Context, toEmail, toName, token, verificationURL string) error {
	log.Println("================================================================================")
	log.Printf("📧 [EMAIL OUTBOX] Verification Email Dispatch")
	log.Printf("   To:       %s <%s>", toName, toEmail)
	log.Printf("   Subject:  Verify Your Email - Next Store")
	log.Printf("   Token:    %s", token)
	log.Printf("   Action:   %s", verificationURL)
	log.Println("================================================================================")
	return nil
}
