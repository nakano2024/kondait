package util

import "github.com/google/uuid"

type uuidGenerator struct{}

func NewUuidGenerator() *uuidGenerator {
	return &uuidGenerator{}
}

func (generator *uuidGenerator) Generate() string {
	return uuid.NewString()
}
