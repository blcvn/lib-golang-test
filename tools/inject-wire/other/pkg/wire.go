package pkg

import (
	"github.com/google/wire"
)

var OtherSet = wire.NewSet(ProvideOther)
