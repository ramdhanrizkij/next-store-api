package worker

import (
	"context"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

type TaskProcessor interface {
	Start() error
	Shutdown()
	ProcessTaskSendEmailVerification(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	mailer Mailer
}

func NewRedisTaskProcessor(redisOpt asynq.RedisConnOpt, mailer Mailer) TaskProcessor {
	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Printf("[WORKER ERROR] Task '%s' failed: %v", task.Type(), err)
			}),
		},
	)

	return &RedisTaskProcessor{
		server: srv,
		mailer: mailer,
	}
}

func (p *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeSendEmailVerification, p.ProcessTaskSendEmailVerification)

	log.Println("[WORKER] Asynq worker server started listening for tasks...")
	return p.server.Run(mux)
}

func (p *RedisTaskProcessor) Shutdown() {
	log.Println("[WORKER] Shutting down Asynq worker server...")
	p.server.Shutdown()
}

func (p *RedisTaskProcessor) ProcessTaskSendEmailVerification(ctx context.Context, task *asynq.Task) error {
	payload, err := DecodeSendEmailVerificationPayload(task.Payload())
	if err != nil {
		return fmt.Errorf("failed to decode verification payload: %w", err)
	}

	log.Printf("[WORKER] Processing verification email for user: %s (%s)", payload.Name, payload.Email)

	if err := p.mailer.SendVerificationEmail(ctx, payload.Email, payload.Name, payload.Token, payload.VerificationURL); err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	log.Printf("[WORKER] Successfully sent verification email to %s", payload.Email)
	return nil
}
