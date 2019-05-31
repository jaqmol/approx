package visualize

import (
	"log"
	"strings"

	"github.com/jaqmol/approx/run"
)

// Hub ...
func Hub(hub *run.Hub) {
	log.Println("Running flow:")
	log.Println("")

	rows := injectArrowRows(makeBoxRows(hub))
	var sb strings.Builder

	for _, r := range rows {
		r.print(&sb)
	}
	log.Println(sb.String())
}

func findMaxWidth(rows []*row) int {
	mw := 0
	for _, r := range rows {
		w := r.width()
		if w > mw {
			mw = w
		}
	}
	return mw
}

func makeBoxRows(hub *run.Hub) (rowsAcc []*row, maxWidth int) {
	rowsAcc = make([]*row, 0)
	for _, publicSource := range hub.PublicProcs {
		sources := []run.Proc{publicSource}
		rowsAcc = append(rowsAcc, newBoxRow(sources))
		for sources != nil && len(sources) > 0 {
			destinations := collectOuts(sources...)
			rowsAcc = append(rowsAcc, newBoxRow(destinations))
			if containsProc(destinations, publicSource) {
				sources = nil
			} else {
				sources = destinations
			}
		}
	}
	maxWidth = findMaxWidth(rowsAcc)
	for _, r := range rowsAcc {
		r.center(maxWidth)
	}
	return rowsAcc, maxWidth
}

func injectArrowRows(boxRows []*row, maxWidth int) []*row {
	rowsAcc := make([]*row, 0)
	lastIndex := len(boxRows) - 1
	for i, boxRow := range boxRows {
		rowsAcc = append(rowsAcc, boxRow)
		if i < lastIndex {
			nextBoxRow := boxRows[i+1]
			froms := boxRow.outPositions()
			tos := nextBoxRow.inPositions()
			a := newArrowRow(froms, tos)
			a.center(maxWidth)
			rowsAcc = append(rowsAcc, a)
		}
	}
	return rowsAcc
}

func collectOuts(ps ...run.Proc) []run.Proc {
	coll := make([]run.Proc, 0)
	checker := make(map[string]bool)
	for _, p := range ps {
		for _, c := range p.Outs() {
			toProcName := c.ToProc.Conf().Name()
			_, contained := checker[toProcName]
			if !contained {
				coll = append(coll, c.ToProc)
				checker[toProcName] = true
			}
		}
	}
	return coll
}

func containsProc(ps []run.Proc, subject run.Proc) bool {
	for _, p := range ps {
		if p == subject {
			return true
		}
	}
	return false
}
