package run

// ViabilityCheck ...
type ViabilityCheck struct {
	inTypeForHash  map[uint32]ConnectionInType
	outTypeForHash map[uint32]ConnectionOutType
}

// NewViabilityCheck ...
func NewViabilityCheck() *ViabilityCheck {
	return &ViabilityCheck{
		inTypeForHash:  make(map[uint32]ConnectionInType),
		outTypeForHash: make(map[uint32]ConnectionOutType),
	}
}

// ContainsIn ...
func (p *ViabilityCheck) ContainsIn(hash uint32) bool {
	_, contained := p.inTypeForHash[hash]
	return contained
}

// AddInType ...
func (p *ViabilityCheck) AddInType(hash uint32, inType ConnectionInType) {
	p.inTypeForHash[hash] = inType
}

// ContainsOut ...
func (p *ViabilityCheck) ContainsOut(hash uint32) bool {
	_, contained := p.outTypeForHash[hash]
	return contained
}

// AddOutType ...
func (p *ViabilityCheck) AddOutType(hash uint32, outType ConnectionOutType) {
	p.outTypeForHash[hash] = outType
}

// InsAndOutsAreBalanced ...
func (p *ViabilityCheck) InsAndOutsAreBalanced() bool {
	for outHash := range p.outTypeForHash {
		if _, contained := p.inTypeForHash[outHash]; !contained {
			return false
		}
	}
	return true
}

// ConnectionInType ...
type ConnectionInType int

// ConnectionInTypes ...
const (
	ConnectionTypeStdin ConnectionInType = iota
	ConnectionTypeEnvInPipe
)

// ConnectionOutType ...
type ConnectionOutType int

// ConnectionOutTypes ...
const (
	ConnectionTypeStdout ConnectionOutType = iota
	ConnectionTypeEnvOutPipe
)
