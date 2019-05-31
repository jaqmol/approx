package visualize

import "strings"

func newGutter(width int, height int) *gutter {
	var sb strings.Builder
	for i := 0; i < width; i++ {
		sb.WriteRune(' ')
	}
	return &gutter{
		_width:  width,
		_height: height,
		_cont:   sb.String(),
	}
}

type gutter struct {
	_width  int
	_height int
	_cont   string
}

func (g *gutter) height() int {
	return g._height
}

func (g *gutter) width() int {
	return g._width
}

func (g *gutter) printLine(sb *strings.Builder, lineIndex int) {
	sb.WriteString(g._cont)
}
