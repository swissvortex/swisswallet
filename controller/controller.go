package controller

import (
	"errors"
	. "vortex-wallet/constants"
	"vortex-wallet/logger"
	"vortex-wallet/model"
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

	//  || arguments.SaltIsEmpry()
	if (arguments.GetCurrencyCode() == CurrencyCode["unknown"]) || arguments.PasswordIsEmpry() || !c.simpleUtils.IsEmptyArray(nonFlagArguments) {
		err := errors.New("Wrong arguments")
		c.logger.LogOnBadRequestErrorWithContext(c.logger.GetContext(), err)
		c.simpleUtils.PrintHelpParamsAndExit(mode)
	}

	c.SwitchFunctionByMode(mode, arguments)
	c.logger.LogOnExitWithContext(c.logger.GetContext())
}

func (c *Controller) SwitchFunctionByMode(mode string, arguments *model.Arguments) {
	c.logger.LogOnEntryWithContext(c.logger.GetContext(), mode, arguments)

	var mapModeToFunction = map[string]func(model.Arguments) error{
		GENERATE_MODE: c.service.GenerateWallet,
		DECRYPT_MODE:  c.service.DecryptWallet,
		ENCRYPT_MODE:  c.service.EncryptWallet,
	}

	err := mapModeToFunction[mode](*arguments)
	if err != nil {
		c.logger.LogOnErrorWithContext(c.logger.GetContext(), err)
		c.simpleUtils.ExitWithError(err)
	}

	c.logger.LogOnExitWithContext(c.logger.GetContext())
}
