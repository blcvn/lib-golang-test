package types

type IptableTask struct {
	Cmd       string `json: cmd`
	Interface string `json: interface`
	Chain     string `json: chain`
	Table     string `json: table`
	Protocol  string `json: protocol`
	SPort     string `json: sport`
	DPort     string `json: dport`
	Action    string `json: action`
	Remote    string `json: remote`
}
