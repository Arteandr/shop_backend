package redis

import (
	"context"
	"errors"
	"time"

	apperrors "shop_backend/pkg/errors"

	"github.com/go-redis/redis/v9"
)

type MailsRepo struct {
	cache    *redis.Client
	cacheTTL time.Duration
}

func NewMailsRepo(cache *redis.Client) *MailsRepo {
	return &MailsRepo{
		cache:    cache,
		cacheTTL: 60 * time.Minute,
	}
}

func (r *MailsRepo) SetVerify(ctx context.Context, token string, userId int) error {
	err := r.cache.Set(ctx, token, userId, r.cacheTTL).Err()

	return err
}

func (r *MailsRepo) GetVerify(ctx context.Context, token string) (string, error) {
	result, err := r.cache.Get(ctx, token).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", apperrors.ErrUserNotFound
		}

		return "", err
	}

	return result, nil
}

func (r *MailsRepo) CompleteVerify(ctx context.Context, token string) error {
	err := r.cache.Del(ctx, token).Err()

	return err
}
