package pkg

type Other struct {
	X int
}

// ProvideFoo returns a Foo.
func ProvideOther() Other {
	return Other{X: 42}
}
