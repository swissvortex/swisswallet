package controller

import (
	. "vortex-wallet/enum"
	"vortex-wallet/service"
	"vortex-wallet/utils"
)

type Controller struct {
	service     service.Service
	simpleUtils utils.SimpleUtils
}

func NewController(service service.Service, simpleUtils utils.SimpleUtils) *Controller {
	return &Controller{
		service:     service,
		simpleUtils: simpleUtils,
	}
}

func (c *Controller) RunVortexWallet() {
	arguments, nonFlagArguments := c.simpleUtils.GetArguments()

	if arguments.GetCurrencyCode() == CurrencyCode["unknown"] || arguments.PassphraseIsEmpry() || arguments.SaltIsEmpry() || !c.simpleUtils.IsEmptyArray(nonFlagArguments) {
		c.simpleUtils.PrintHelp()
	}

	c.service.GenerateWallet(*arguments)
}
