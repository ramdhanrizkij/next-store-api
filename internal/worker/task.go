package worker

import "encoding/json"

const (
	TypeSendEmailVerification = "task:send_email_verification"
)

type SendEmailVerificationPayload struct {
	UserID          string `json:"user_id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Token           string `json:"token"`
	VerificationURL string `json:"verification_url"`
}

func (p *SendEmailVerificationPayload) Encode() ([]byte, error) {
	return json.Marshal(p)
}

func DecodeSendEmailVerificationPayload(data []byte) (*SendEmailVerificationPayload, error) {
	var payload SendEmailVerificationPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	return &payload, nil
}
