package fakie

import "math/rand"

func (g *Generator) Get(p string, r *rand.Rand) (string, error) {
	return g.GetWithArgs(p, nil, r)
}

func (g *Generator) GetWithArgs(p string, args []string, r *rand.Rand) (string, error) {
	gens := g.generators
	dict := g.dictionary

	if r == nil {
		r = CreateRand()
	}

	if f, ok := gens[p]; ok {
		return f(g, args, r)
	}

	if list, ok := dict[p]; ok {
		i := r.Intn(len(list))
		return list[i], nil
	}

	return "", nil
}

// List all paths.
func (g *Generator) List() []string {
	gens := g.generators
	dict := g.dictionary
	ret := make([]string, 0)

	for k := range gens {
		ret = append(ret, k)
	}

	for k := range dict {
		ret = append(ret, k)
	}

	return ret
}
