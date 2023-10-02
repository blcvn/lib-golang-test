package foobarbaz

import (
	"inject-wire/other/pkg"

	"github.com/google/wire"
)

var SuperSet = wire.NewSet(ProvideFoo, ProvideBar, ProvideBaz)
var MegaSet = wire.NewSet(SuperSet, pkg.OtherSet)
