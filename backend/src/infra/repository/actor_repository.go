package repository

import (
	"context"

	"gorm.io/gorm"

	"kondait-backend/domain/entity"
	domainrepo "kondait-backend/domain/repository"
)

type actorRepository struct {
	db *gorm.DB
}

func NewActorRepository(db *gorm.DB) domainrepo.IActorRepository {
	return &actorRepository{
		db: db,
	}
}

func (repo *actorRepository) FetchBySub(ctx context.Context, sub string) (*entity.Actor, error) {
	_ = repo.db.WithContext(ctx)
	_ = sub
	// TODO: implement actor lookup by sub.
	return &entity.Actor{
		Code: "",
	}, nil
}

type actorRepositoryMock struct{}

func NewActorRepositoryMock() domainrepo.IActorRepository {
	return &actorRepositoryMock{}
}

func (repo *actorRepositoryMock) FetchBySub(ctx context.Context, sub string) (*entity.Actor, error) {
	_ = ctx
	return &entity.Actor{
		Code: "11111111-1111-1111-1111-111111111111",
	}, nil
}
