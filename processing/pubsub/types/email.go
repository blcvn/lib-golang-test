package types

type EmailTask struct {
	Email        string `json:"email"`
	UserId       uint64 `json:"userId"`
	ActivateCode string `json:"activateCode"`
	Reset        bool   `json:"reset"`
}
