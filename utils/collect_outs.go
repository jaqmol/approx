package utils

import (
	"github.com/jaqmol/approx/flow"
)

// CollectOuts ...
func CollectOuts(ps ...*flow.ProcItem) []*flow.ProcItem {
	coll := make([]*flow.ProcItem, 0)
	checker := make(map[string]bool)
	for _, p := range ps {
		for _, nextConn := range p.Next {
			toProcName := nextConn.To.Conf.Name()
			_, contained := checker[toProcName]
			if !contained {
				coll = append(coll, nextConn.To)
				checker[toProcName] = true
			}
		}
	}
	return coll
}
