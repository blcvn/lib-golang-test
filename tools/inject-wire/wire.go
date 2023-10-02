package main

import (
	"github.com/google/wire"
)

// initApp returns a real app.
func initApp() *app {
	wire.Build(appSet)
	return nil
}

// initMockedAppFromArgs returns an app with mocked dependencies provided via
// arguments (Approach A). Note that the argument's type is the interface
// type (timer), but the concrete mock type should be passed.

func initMockedAppFromArgs(mt timer) *app {
	wire.Build(appSetWithoutMocks)
	return nil
}

// initMockedApp returns an app with its mocked dependencies, created
// via providers (Approach B).

func initMockedApp() *appWithMocks {
	wire.Build(mockAppSet)
	return nil
}
