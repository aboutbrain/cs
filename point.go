package cs

type Context struct {
	ConceptId int
}

type Concept struct {
	ContextId int
}

type Point struct {
	Memory [][BitsPerPoint]byte
	/*InputArrayMap [BitsPerPoint]int
	Concept       [PointMemoryCapacity]Concept
	Context       [PointContextCapacity]Context*/
}

/*func (p *Point) SetConcept(concept Concept, context Context) {
	p.Concept[context.ConceptId] = concept
	p.Context[concept.ContextId] = context
}

func (p *Point) GetConceptByContext(context Context) Concept {
	return p.Concept[context.ConceptId]
}

func (p *Point) GetContextByConcept(concept Concept) Context {
	return p.Context[concept.ContextId]
}*/
