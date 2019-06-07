package flow

// ConnItem ...
type ConnItem struct {
	From       *ProcItem
	To         *ProcItem
	IsLoopBack bool
}

// NewConnItem ...
func NewConnItem(from *ProcItem, to *ProcItem) *ConnItem {
	return &ConnItem{
		From: from,
		To:   to,
	}
}

// // Type ...
// func (ci *ConnItem) Type() ItemType {
// 	return ItemTypeConn
// }

// // PreviousCount ...
// func (ci *ConnItem) PreviousCount() int {
// 	return 1
// }

// // Previous ...
// func (ci *ConnItem) Previous(int) Item {
// 	return ci.from
// }

// // NextCount ...
// func (ci *ConnItem) NextCount() int {
// 	return 1
// }

// // Next ...
// func (ci *ConnItem) Next(int) Item {
// 	return ci.to
// }

// // Visit ...
// func (ci *ConnItem) Visit(ito func(Item)) {
// 	ito(ci)
// 	if !ci.isLoopBack && ci.to != nil {
// 		ci.to.Visit(ito)
// 	}
// }
