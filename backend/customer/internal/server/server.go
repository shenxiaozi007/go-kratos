package server

import (
	"customer/internal/service"

	"github.com/google/wire"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer, service.NewCustomerService)
