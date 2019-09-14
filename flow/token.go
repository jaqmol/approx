package flow

// import (
// 	"regexp"
// )

// type token interface {
// 	kind() kind
// }

// type processor struct {
// 	name      string
// 	tokenKind kind
// }

// func (p *processor) kind() kind {
// 	return p.tokenKind
// }

// type pipe struct {
// 	tokenKind kind
// }

// func (p *pipe) kind() kind {
// 	return p.tokenKind
// }

// type tappedPipe struct {
// 	name      string
// 	tokenKind kind
// }

// func (tp *tappedPipe) kind() kind {
// 	return tp.tokenKind
// }

// type kind int

// const (
// 	kindUnknown kind = iota
// 	kindProcessor
// 	kindPipe
// 	kindTappedPipe
// )

// func parseTokens(flowDefinition []string) [][]token {
// 	splitter := regexp.MustCompile("-(\\w+)->|->")
// 	acc := make([][]token, 0)
// 	for _, lineDefinition := range flowDefinition {
// 		procNames := splitter.Split(lineDefinition, -1)
// 		pipeNameSubs := splitter.FindAllStringSubmatch(lineDefinition, -1)
// 		tokens := collectTokens(procNames, pipeNameSubs)
// 		acc = append(acc, tokens)
// 	}
// 	return acc
// }
