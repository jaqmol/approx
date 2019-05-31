package visualize

import (
	"strings"
)

func newArrow(width int, at arrowType) *arrow {
	return &arrow{
		_width:     width,
		_arrowType: at,
	}
}

type arrowType int

const (
	arrowTypeStraight arrowType = iota
	arrowTypeTopLeftToBottomRight
	arrowTypeTopRightToBottomLeft
)

type arrow struct {
	_width     int
	_arrowType arrowType
}

func (a *arrow) height() int {
	if a._arrowType == arrowTypeStraight {
		return 3
	}
	return 4
}

func (a *arrow) width() int {
	return a._width
}

func (a *arrow) printLine(sb *strings.Builder, lineIndex int) {
	if a._arrowType == arrowTypeStraight {
		a.printStraightArrow(sb, lineIndex)
	} else if a._arrowType == arrowTypeTopLeftToBottomRight {
		a.printTopLeftToBottomRight(sb, lineIndex)
	} else if a._arrowType == arrowTypeTopRightToBottomLeft {
		a.printTopRightToBottomLeft(sb, lineIndex)
	}
}

func (a *arrow) printStraightArrow(sb *strings.Builder, lineIndex int) {
	if lineIndex == 2 {
		sb.WriteRune('V')
	} else {
		sb.WriteRune('│')
	}
}

func (a *arrow) printTopLeftToBottomRight(sb *strings.Builder, lineIndex int) {
	if lineIndex == 0 {
		sb.WriteRune('│')
		fillWithChar(sb, ' ', a._width)
	} else if lineIndex == 1 {
		sb.WriteRune('└')
		fillWithChar(sb, '─', a._width-1)
		sb.WriteRune('┐')
	} else if lineIndex == 2 {
		fillWithChar(sb, ' ', a._width)
		sb.WriteRune('│')
	} else if lineIndex == 3 {
		fillWithChar(sb, ' ', a._width)
		sb.WriteRune('V')
	}
}

func (a *arrow) printTopRightToBottomLeft(sb *strings.Builder, lineIndex int) {
	if lineIndex == 0 {
		fillWithChar(sb, ' ', a._width)
		sb.WriteRune('│')
	} else if lineIndex == 1 {
		sb.WriteRune('┌')
		fillWithChar(sb, '─', a._width-1)
		sb.WriteRune('┘')
	} else if lineIndex == 2 {
		sb.WriteRune('│')
		fillWithChar(sb, ' ', a._width)
	} else if lineIndex == 3 {
		sb.WriteRune('V')
		fillWithChar(sb, ' ', a._width)
	}
}

func fillWithChar(sb *strings.Builder, char rune, times int) {
	for i := 0; i < times; i++ {
		sb.WriteRune(char)
	}
}
