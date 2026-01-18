package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"kondait-backend/domain/entity"
	domainrepo "kondait-backend/domain/repository"
	"kondait-backend/infra/model"
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
	var user model.User
	err := repo.db.WithContext(ctx).
		Where("sub = ?", sub).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity.Actor{
		Code: user.Code,
		Sub:  user.Sub,
	}, nil
}

func (repo *actorRepository) Save(ctx context.Context, actor *entity.Actor) error {
	if actor == nil {
		return errors.New("actorRepository: actor is nil")
	}

	db := repo.db.WithContext(ctx)
	if actor.Code != "" {
		var byCode model.User
		err := db.Where("code = ?", actor.Code).First(&byCode).Error
		if err == nil {
			return db.Model(&byCode).Updates(map[string]interface{}{
				"sub": actor.Sub,
			}).Error
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	if actor.Sub != "" {
		var bySub model.User
		err := db.Where("sub = ?", actor.Sub).First(&bySub).Error
		if err == nil {
			updates := map[string]interface{}{
				"sub": actor.Sub,
			}
			if actor.Code != "" {
				updates["code"] = actor.Code
			}
			return db.Model(&bySub).Updates(updates).Error
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	return db.Create(&model.User{
		Code: actor.Code,
		Sub:  actor.Sub,
	}).Error
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

func (repo *actorRepositoryMock) Save(ctx context.Context, actor *entity.Actor) error {
	_ = ctx
	_ = actor
	return nil
}
