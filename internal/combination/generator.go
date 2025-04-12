package combination

type Generator interface {
	Generate() *Combination
}

type StarlarkGenerator struct {
}

func (s *StarlarkGenerator) Generate() *Combination {
	return nil
}
