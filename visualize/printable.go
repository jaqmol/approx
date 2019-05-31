package visualize

import "strings"

type printable interface {
	height() int
	width() int
	printLine(sb *strings.Builder, lineIndex int)
}
