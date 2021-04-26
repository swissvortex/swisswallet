package controller

import (
	"errors"
	. "vortex-wallet/constants"
	"vortex-wallet/logger"
	"vortex-wallet/service"
	"vortex-wallet/utils"
)

type Controller struct {
	service     service.Service
	simpleUtils utils.SimpleUtils
	logger      *logger.Logger
}

func NewController(service service.Service, simpleUtils utils.SimpleUtils, logger *logger.Logger) *Controller {
	return &Controller{
		service:     service,
		simpleUtils: simpleUtils,
		logger:      logger,
	}
}

func (c *Controller) RunVortexWallet() {
	arguments, mode, nonFlagArguments := c.simpleUtils.GetArguments()
	c.logger.LogOnEntryWithContext(c.logger.GetContext(), arguments, nonFlagArguments)

	if c.simpleUtils.IsEmptyString(mode) {
		err := errors.New("Missing operation mode")
		c.logger.LogOnBadRequestErrorWithContext(c.logger.GetContext(), err)
		c.simpleUtils.PrintHelpMode()
	}

	err := c.simpleUtils.CheckMode(mode)
	if err != nil {
		c.logger.LogOnBadRequestErrorWithContext(c.logger.GetContext(), err)
		c.simpleUtils.ExitWithError(err)
	}

	if (arguments.GetCurrencyCode() == CurrencyCode["unknown"]) || arguments.PasswordIsEmpry() || arguments.SaltIsEmpry() || !c.simpleUtils.IsEmptyArray(nonFlagArguments) {
		err := errors.New("Wrong arguments")
		c.logger.LogOnBadRequestErrorWithContext(c.logger.GetContext(), err)
		c.simpleUtils.PrintHelpParams(mode)
	}

	err = c.service.GenerateWallet(*arguments)
	if err != nil {
		c.logger.LogOnErrorWithContext(c.logger.GetContext(), err)
		c.simpleUtils.ExitWithError(err)
	}
}
