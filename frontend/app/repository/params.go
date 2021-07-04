package repository

import (
	"fmt"
)

type Params struct {
	Endpoint string
}

func NewParams(addr, usr, pass string) Params {
	return Params{Endpoint: fmt.Sprintf("mongodb://%s:%s@%s:27017", usr, pass, addr)}
}
