//+build !wireinject

package main

import (
	"context"
	// "inject-wire/foobar"
	"inject-wire/foobarbaz"
)

// Injectors from main.go:

func initializeBaz(ctx context.Context) (foobarbaz.Baz, error) {
	foo := foobarbaz.ProvideFoo()
	bar := foobarbaz.ProvideBar(foo)
	baz, err := foobarbaz.ProvideBaz(ctx, bar)
	if err != nil {
		return foobarbaz.Baz{}, err
	}
	return baz, nil
}

// func initializeNewBar(ctx context.Context) (string, error) {
// 	myFooer := foobar.provideMyFooer()
// 	string2 := foobar.provideBar(myFooer)
// 	return string2, nil
// }

// Injectors from wire.go:

// initApp returns a real app.
func initApp() *app {
	mainTimer := _wireRealTimeValue
	mainGreeter := greeter{
		T: mainTimer,
	}
	mainApp := &app{
		g: mainGreeter,
	}
	return mainApp
}

var (
	_wireRealTimeValue = realTime{}
)

// initMockedAppFromArgs returns an app with mocked dependencies provided via
// arguments (Approach A). Note that the argument's type is the interface
// type (timer), but the concrete mock type should be passed.
func initMockedAppFromArgs(mt timer) *app {
	mainGreeter := greeter{
		T: mt,
	}
	mainApp := &app{
		g: mainGreeter,
	}
	return mainApp
}

// initMockedApp returns an app with its mocked dependencies, created
// via providers (Approach B).
func initMockedApp() *appWithMocks {
	mainMockTimer := newMockTimer()
	mainGreeter := greeter{
		T: mainMockTimer,
	}
	mainApp := app{
		g: mainGreeter,
	}
	mainAppWithMocks := &appWithMocks{
		app: mainApp,
		mt:  mainMockTimer,
	}
	return mainAppWithMocks
}
