package repository

import (
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

func (repo *actorRepository) FetchBySub(sub string) (*entity.Actor, error) {
	_ = repo.db
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

func (repo *actorRepositoryMock) FetchBySub(sub string) (*entity.Actor, error) {
	return &entity.Actor{
		Code: "11111111-1111-1111-1111-111111111111",
	}, nil
}
