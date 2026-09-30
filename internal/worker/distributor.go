package worker

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskSendEmailVerification(
		ctx context.Context,
		payload *SendEmailVerificationPayload,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisConnOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{client: client}
}

func (d *RedisTaskDistributor) DistributeTaskSendEmailVerification(
	ctx context.Context,
	payload *SendEmailVerificationPayload,
	opts ...asynq.Option,
) error {
	jsonPayload, err := payload.Encode()
	if err != nil {
		return fmt.Errorf("failed to encode email verification payload: %w", err)
	}

	task := asynq.NewTask(TypeSendEmailVerification, jsonPayload, opts...)
	info, err := d.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	_ = info
	return nil
}
