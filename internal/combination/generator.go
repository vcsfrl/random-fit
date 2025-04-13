package combination

type Builder interface {
	Build() *Combination
}

type StarlarkBuilder struct {
}

func (s *StarlarkBuilder) Build() *Combination {
	return nil
}
