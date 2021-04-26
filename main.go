package main

import (
	. "vortex-wallet/constants"
	"vortex-wallet/logger"
	repo "vortex-wallet/repository"
	"vortex-wallet/utils"

	"vortex-wallet/controller"
	"vortex-wallet/service"
)

func main() {
	logger := logger.NewLogger()
	logger.SetLoggingLevel(LOGGING_LEVEL)

	utils := utils.NewSimpleUtils(logger)
	cryptoRepository := repo.NewCryptoRepository(logger)
	service := service.NewService(cryptoRepository, logger)
	controller := controller.NewController(service, utils, logger)

	controller.RunVortexWallet()
}
