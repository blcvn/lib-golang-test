package foobar

type Bar string

func provideBar(f Fooer) string {
	// f will be a *MyFooer.
	return f.Foo()
}
