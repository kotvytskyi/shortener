package mongo

import (
	"fmt"

	"github.com/kotvytskyi/frontend/pkg/config"
)

type Params struct {
	Endpoint string
}

func NewParams(cfg config.MongoConfig) Params {
	return Params{Endpoint: fmt.Sprintf("mongodb://%s:%s@%s:27017", cfg.User, cfg.Password, cfg.Address)}
}
