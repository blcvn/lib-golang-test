package types

import "github.com/google/uuid"

type Task struct {
	Id    uuid.UUID
	Path  string
	Data  []byte
	UseId uint64
}

type TaskResponse struct {
	Id      string
	Code    int
	Message string
	Path    string
	Data    interface{}
}
