package main

import (
	repo "vortex-wallet/repository"
	"vortex-wallet/utils"

	"vortex-wallet/controller"
	"vortex-wallet/service"
)

func main() {
	cryptoRepository := repo.NewCryptoRepository()
	service := service.NewService(cryptoRepository)
	utils := utils.NewSimpleUtils()
	controller := controller.NewController(service, utils)
	controller.RunVortexWallet()
}
