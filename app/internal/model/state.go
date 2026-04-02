package model

import "time"

type ContractState struct {
	ID        int64
	Value     uint64
	UpdatedAt time.Time
}
