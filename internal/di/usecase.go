package di

import "github.com/moemoe89/btc/internal/usecases"

// GetBTCUsecase returns BTCUsecase instance.
func GetBTCUsecase() usecases.BTCUsecase {
	return usecases.NewBTCUsecase(
		GetBTCRepo(),
		GetTracer().Tracer(),
		GetLogger(),
		GetRedis(),
	)
}
