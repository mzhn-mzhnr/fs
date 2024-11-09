package uuid

import "github.com/google/uuid"

type UuidProvider struct {
}

func New() *UuidProvider {
	return &UuidProvider{}
}

func (u *UuidProvider) Provide() string {
	id := uuid.New().String()
	return id
}
