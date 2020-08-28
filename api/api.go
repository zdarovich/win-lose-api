package api

import (
	"fmt"
)

type (
	// api struct
	api struct {
		APIRouter IRouter
	}
	// IAPI interface
	IAPI interface {
		Run(port uint)
	}
)

// New get new configured api
func New(router IRouter) IAPI {

	return &api{
		APIRouter: router,
	}
}

// Run starts the api
func (api *api) Run(port uint) {

	apiRouter := api.APIRouter.GetEngine()
	apiRouter.Run(fmt.Sprintf(":%d", port))
}
