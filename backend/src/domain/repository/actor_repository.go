package repository

import "kondait-backend/domain/entity"

type IActorRepository interface {
	FetchBySub(sub string) (*entity.Actor, error)
}
