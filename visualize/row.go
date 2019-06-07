package visualize

import (
	"strings"

	"github.com/jaqmol/approx/flow"
)

func newBoxRow(procs []*flow.ProcItem) *row {
	printables := make([]printable, 0)
	for i, p := range procs {
		if i > 0 {
			printables = append(printables, newGutter(2, 3))
		}
		n := p.Conf.Name()
		b := newBox(
			n,
			len(p.Conf.Inputs()),
			len(p.Conf.Outputs()),
		)
		printables = append(printables, b)
	}

	pw, ph := maxWidthAndHeight(printables)

	return &row{
		inset:            0,
		printablesHeight: ph,
		printablesWidth:  pw,
		printables:       printables,
	}
}

func newArrowRow(fromPositions []int, toPositions []int) *row {
	printables := make([]printable, 0)
	count := len(fromPositions)
	inset := 0
	if count > 1 {
		inset = 1
	}
	for i := 0; i < count; i++ {
		from := fromPositions[i]
		to := toPositions[i]
		width := max(from, to) - min(from, to)
		var at arrowType
		var g *gutter
		if from < to {
			at = arrowTypeTopLeftToBottomRight
			g = newGutter(1, 4)
		} else if from > to {
			at = arrowTypeTopRightToBottomLeft
			g = newGutter(1, 4)
		} else {
			at = arrowTypeStraight
			g = newGutter(1, 3)
		}
		if i > 0 {
			printables = append(printables, g)
		}
		a := newArrow(width, at)
		printables = append(printables, a)
	}

	pw, ph := maxWidthAndHeight(printables)

	return &row{
		inset:            inset,
		printablesHeight: ph,
		printablesWidth:  pw,
		printables:       printables,
	}
}

func maxWidthAndHeight(printables []printable) (width int, height int) {
	height = printables[0].height()
	width = 0
	for _, p := range printables {
		h := p.height()
		if h > height {
			height = h
		}
		width += p.width()
	}
	return
}

type row struct {
	inset            int
	printablesHeight int
	printablesWidth  int
	printables       []printable
}

func (r *row) height() int {
	return r.printablesHeight
}

func (r *row) width() int {
	return r.printablesWidth
}

func (r *row) center(width int) {
	r.inset += (width / 2) - (r.printablesWidth / 2)
}

func (r *row) inPositions() []int {
	return r.collectPositions(func(b *box) []int {
		return b.inPositions
	})
}

func (r *row) outPositions() []int {
	return r.collectPositions(func(b *box) []int {
		return b.outPositions
	})
}

func (r *row) collectPositions(selectPositions func(b *box) []int) []int {
	acc := make([]int, 0)
	inset := r.inset
	for _, p := range r.printables {
		if box, ok := p.(*box); ok {
			for _, pos := range selectPositions(box) {
				acc = append(acc, inset+pos)
			}
		}
		inset += p.width()
	}
	return acc
}

func (r *row) print(sb *strings.Builder) {
	rowHeight := r.height()
	for i := 0; i < rowHeight; i++ {
		for j := 0; j < r.inset; j++ {
			sb.WriteRune(' ')
		}
		for _, p := range r.printables {
			p.printLine(sb, i)
		}
		sb.WriteRune('\n')
	}
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
