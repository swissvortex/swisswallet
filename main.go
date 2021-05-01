package main

import (
	"fmt"
	. "swisswallet/constants"
	"swisswallet/logger"
	repo "swisswallet/repository"
	"swisswallet/utils"
	"time"

	"swisswallet/controller"
	"swisswallet/service"
)

func main() {
	start := time.Now()

	logger := logger.NewLogger()
	logger.SetLoggingLevel(LOGGING_LEVEL)
	utils := utils.NewSimpleUtils(logger)
	cryptoRepository := repo.NewCryptoRepository(logger)
	service := service.NewService(cryptoRepository, utils, logger)
	controller := controller.NewController(service, utils, logger)
	controller.RunSwissWallet()

	fmt.Printf("Elapsed time: %s\n", time.Now().Sub(start))
}
