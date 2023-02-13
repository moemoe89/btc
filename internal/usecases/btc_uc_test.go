package usecases_test

import (
	"github.com/moemoe89/btc/internal/di"
	"github.com/moemoe89/btc/internal/entities/repository"
	"github.com/moemoe89/btc/internal/usecases"
)

type fields struct {
	btcRepo repository.BTCRepo
}

func sut(f fields) usecases.BTCUsecase {
	return usecases.NewBTCUsecase(
		f.btcRepo,
		di.GetTracer().Tracer(),
	)
}
