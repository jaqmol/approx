package visualize

import (
	"log"
	"strings"

	"github.com/jaqmol/approx/flow"
)

// Flow ...
func Flow(flo *flow.Flow) {
	log.SetFlags(0)
	log.Println("Running flow:")
	log.Println("")

	rows := injectArrowRows(makeBoxRows(flo))
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

func makeBoxRows(flo *flow.Flow) (rowsAcc []*row, maxWidth int) {
	rowsAcc = make([]*row, 0)
	flo.IterateProcs(func(procRow []*flow.ProcItem) {
		rowsAcc = append(rowsAcc, newBoxRow(procRow))
	})
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
