package flow

// ConnSet ...
type ConnSet struct {
	data      []*ConnItem
	contained map[uint32]bool
}

// NewConnSet ...
func NewConnSet(firstConns []*ConnItem) *ConnSet {
	set := &ConnSet{
		data:      make([]*ConnItem, 0),
		contained: make(map[uint32]bool),
	}
	if firstConns != nil && len(firstConns) > 0 {
		for _, c := range firstConns {
			set.Add(c)
		}
	}
	return set
}

// IsContained ...
func (ps *ConnSet) IsContained(conn *ConnItem) bool {
	ok1, ok2 := ps.contained[conn.Hash]
	if ok2 {
		return ok1
	}
	return false
}

// Add ...
func (ps *ConnSet) Add(conn *ConnItem) bool {
	if !ps.IsContained(conn) {
		ps.data = append(ps.data, conn)
		ps.contained[conn.Hash] = true
		return true
	}
	return false
}

// AddAll ...
func (ps *ConnSet) AddAll(conns []*ConnItem) {
	for _, c := range conns {
		ps.Add(c)
	}
}

// Len ...
func (ps *ConnSet) Len() int {
	return len(ps.data)
}

// Data ...
func (ps *ConnSet) Data() []*ConnItem {
	return ps.data
}
