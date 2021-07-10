package generator

import "context"

type KeyRepository interface {
	Create(ctx context.Context, key string) error
}

type Generator struct {
	repository KeyRepository
}

func NewGenerator(repository KeyRepository) *Generator {
	gen := &Generator{}
	gen.repository = repository
	return gen
}

func (g *Generator) Generate(ctx context.Context) error {
	randKey := generateRandomString(6)
	err := g.repository.Create(ctx, randKey)
	return err
}
