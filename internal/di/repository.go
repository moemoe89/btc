package di

import (
	"github.com/moemoe89/btc/internal/entities/repository"
	"github.com/moemoe89/btc/internal/infrastructure/datastore"
)

// GetBaseRepo returns BaseRepo instance.
func GetBaseRepo() *datastore.BaseRepo {
	return datastore.NewBaseRepo(datastore.GetDatabase())
}

// GetBTCRepo returns BTCRepo instance.
func GetBTCRepo() repository.BTCRepo {
	return datastore.NewBTCRepo(GetBaseRepo())
}
