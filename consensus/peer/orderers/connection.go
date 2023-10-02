package orderers

type Endpoint struct {
	Address   string
	RootCerts [][]byte
	Refreshed chan struct{}
}
