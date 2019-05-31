package visualize

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func newBox(content string, insCount int, outsCount int) *box {
	middleLine := fmt.Sprintf("│ %v │", content)
	middleLineLen := utf8.RuneCountInString(middleLine)

	var builder strings.Builder

	horizLineLen := middleLineLen - 2
	for i := 0; i < horizLineLen; i++ {
		builder.WriteRune('─')
	}
	horizLine := builder.String()
	topLine := fmt.Sprintf("┌%v┐", horizLine)

	builder.Reset()

	outPositions := inOutPositions(outsCount, middleLineLen)
	lastMiddleLineIndex := middleLineLen - 1
	for i := 0; i < middleLineLen; i++ {
		if containsInt(outPositions, i) {
			if i == 0 {
				builder.WriteRune('├')
			} else if i == lastMiddleLineIndex {
				builder.WriteRune('┤')
			} else {
				builder.WriteRune('┬')
			}
		} else {
			if i == 0 {
				builder.WriteRune('└')
			} else if i == lastMiddleLineIndex {
				builder.WriteRune('┘')
			} else {
				builder.WriteRune('─')
			}
		}
	}
	bottomLine := builder.String()

	return &box{
		lineLen: middleLineLen,
		lines: []string{
			topLine,
			middleLine,
			bottomLine,
		},
		inPositions:  inOutPositions(insCount, middleLineLen),
		outPositions: outPositions,
	}
}

func inOutPositions(count int, boxWidth int) []int {
	bitmap := make([]bool, 0)
	for i := 0; i < count; i++ {
		if i > 0 {
			bitmap = append(bitmap, false)
		}
		bitmap = append(bitmap, true)
	}
	start := (boxWidth / 2) - (len(bitmap) / 2)
	acc := make([]int, 0)
	for i, has := range bitmap {
		if has {
			acc = append(acc, start+i)
		}
	}
	return acc
}

type box struct {
	lineLen      int
	lines        []string
	inPositions  []int
	outPositions []int
}

func (b *box) height() int {
	return 3
}

func (b *box) width() int {
	return b.lineLen
}

func (b *box) printLine(sb *strings.Builder, lineIndex int) {
	l := b.lines[lineIndex]
	sb.WriteString(l)
}

func containsInt(all []int, query int) bool {
	for _, i := range all {
		if i == query {
			return true
		}
	}
	return false
}
