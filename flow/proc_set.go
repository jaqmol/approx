package flow

// ProcSet ...
type ProcSet struct {
	data      []*ProcItem
	contained map[string]bool
}

// NewProcSet ...
func NewProcSet(firstProc *ProcItem) *ProcSet {
	data := make([]*ProcItem, 0)
	contained := make(map[string]bool)
	if firstProc != nil {
		data = append(data, firstProc)
		contained[firstProc.Conf.Name()] = true
	}
	return &ProcSet{
		data:      data,
		contained: contained,
	}
}

// IsContained ...
func (ps *ProcSet) IsContained(proc *ProcItem) bool {
	ok1, ok2 := ps.contained[proc.Conf.Name()]
	if ok2 {
		return ok1
	}
	return false
}

// Add ...
func (ps *ProcSet) Add(proc *ProcItem) bool {
	if !ps.IsContained(proc) {
		ps.data = append(ps.data, proc)
		ps.contained[proc.Conf.Name()] = true
		return true
	}
	return false
}

// Len ...
func (ps *ProcSet) Len() int {
	return len(ps.data)
}

// Data ...
func (ps *ProcSet) Data() []*ProcItem {
	return ps.data
}
