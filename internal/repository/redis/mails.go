package redis

import (
	"context"
	"github.com/go-redis/redis/v9"
	"time"
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

func (r *MailsRepo) SetVerify(ctx context.Context, verifyKey string, userId int) error {
	err := r.cache.Set(ctx, verifyKey, userId, r.cacheTTL).Err()

	return err
}

func (r *MailsRepo) ExistVerify(ctx context.Context, verifyKey string) (bool, error) {
	if err := r.cache.Exists(ctx, verifyKey).Err(); err != nil {
		return false, nil
	} else {
		return true, nil
	}
}
