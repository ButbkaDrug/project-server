package generator

// a struct intended to be used as singletone to generate contigues ids in different parst of a program
type Generator struct {
	value int
}

// returns a new generator with value staring with 1
func New() *Generator {
	return &Generator{
		value: 1,
	}
}

// generates next value
func (g *Generator) Get() int {
	v := g.value

	g.value++

	return v
}
