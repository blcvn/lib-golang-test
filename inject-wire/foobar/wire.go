package foobar

import (
	"github.com/google/wire"
)

var NewBarSet = wire.NewSet(
	provideMyFooer,
	wire.Bind(new(Fooer), new(*MyFooer)),
	provideBar)
