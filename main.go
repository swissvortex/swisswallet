package main

import (
	"fmt"
	"time"
	. "vortex-wallet/constants"
	"vortex-wallet/logger"
	repo "vortex-wallet/repository"
	"vortex-wallet/utils"

	"vortex-wallet/controller"
	"vortex-wallet/service"
)

func main() {
	start := time.Now()

	logger := logger.NewLogger()
	logger.SetLoggingLevel(LOGGING_LEVEL)
	utils := utils.NewSimpleUtils(logger)
	cryptoRepository := repo.NewCryptoRepository(logger)
	service := service.NewService(cryptoRepository, utils, logger)
	controller := controller.NewController(service, utils, logger)
	controller.RunVortexWallet()

	fmt.Printf("Elapsed time: %s\n", time.Now().Sub(start))
}
