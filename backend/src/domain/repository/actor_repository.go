package repository

import (
	"context"
	"kondait-backend/domain/entity"
)

type IActorRepository interface {
	FetchBySub(ctx context.Context, sub string) (*entity.Actor, error)
}
