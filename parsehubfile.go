package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

// NodeClass ...
type NodeClass int

// NodeClass ...
const (
	ProcessClass NodeClass = iota
	ForkClass
	MergeClass
	PipeClass
)

// Node ...
type Node struct {
	Class   NodeClass
	ID      string
	Command string
	OutKeys []string
}

// NodeCollection ...
type NodeCollection struct {
	index     map[string]Node
	indexKeys []string
}

// NewNodeCollection ...
func NewNodeCollection(filePath string) NodeCollection {
	index := parseHubFile(filePath)
	return NodeCollection{
		index: index,
	}
}

// IDs ...
func (c NodeCollection) IDs() []string {
	if c.indexKeys == nil {
		c.indexKeys = make([]string, len(c.index))
		i := 0
		for k := range c.index {
			c.indexKeys[i] = k
			i++
		}
	}
	return c.indexKeys
}

// Node ...
func (c NodeCollection) Node(id string) (Node, bool) {
	n, ok := c.index[id]
	return n, ok
}

// Outputs ...
func (c NodeCollection) Outputs(n Node) []Node {
	nodes := make([]Node, len(n.OutKeys))
	for i, key := range n.OutKeys {
		nodes[i] = c.index[key]
	}
	return nodes
}

func parseHubFile(filePath string) map[string]Node {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("Error reading hub file: ", filePath, ": ", err)
	}

	lines := parseLines(fileContent)
	connectedPairs := makeConnectedPairs(lines)

	mainIndex := make(map[string]Node)
	connectionIndex := make(map[string]Node)

	// 1st setup communication nodes

	for source, destinations := range findForkConnections(connectedPairs) {
		n := Node{
			Class:   ForkClass,
			ID:      makeForkID(source, destinations),
			OutKeys: mapToCmdIDs(destinations...),
		}
		mainIndex[n.ID] = n
		connectionIndex[forkFromKey(source)] = n
		for _, destination := range destinations {
			connectionIndex[forkToKey(destination)] = n
		}
	}

	for destination, sources := range findMergeConnections(connectedPairs) {
		n := Node{
			Class:   MergeClass,
			ID:      makeMergeID(sources, destination),
			OutKeys: mapToCmdIDs(destination),
		}
		mainIndex[n.ID] = n
		connectionIndex[mergeToKey(destination)] = n
		for _, source := range sources {
			connectionIndex[mergeFromKey(source)] = n
		}
	}

	for source, destination := range findPipeConnections(connectedPairs) {
		n := Node{
			Class:   PipeClass,
			ID:      makePipeID(source, destination),
			OutKeys: mapToCmdIDs(destination),
		}
		mainIndex[n.ID] = n
		connectionIndex[pipeFromKey(source)] = n
		connectionIndex[pipeToKey(destination)] = n
	}

	// 2nd setup command process nodes

	for _, cmd := range uniqueProcessNames(lines) {
		n := Node{
			Class:   ProcessClass,
			ID:      makeCmdID(cmd),
			Command: cmd,
		}

		next, ok := connectionIndex[pipeFromKey(cmd)]
		if ok {
			n.OutKeys = append(n.OutKeys, next.ID)
		}

		next, ok = connectionIndex[forkFromKey(cmd)]
		if ok {
			n.OutKeys = append(n.OutKeys, next.ID)
		}

		next, ok = connectionIndex[mergeFromKey(cmd)]
		if ok {
			n.OutKeys = append(n.OutKeys, next.ID)
		}

		mainIndex[n.ID] = n
	}

	return mainIndex
}

func makeForkID(fromCmd string, toCmds []string) string {
	return fmt.Sprintf("fork:%s->%s", fromCmd, join(toCmds, ","))
}

func makeMergeID(fromCmds []string, toCmd string) string {
	return fmt.Sprintf("merge:%s->%s", join(fromCmds, ","), toCmd)
}

func makePipeID(fromCmd string, toCmd string) string {
	return fmt.Sprintf("pipe:%s->%s", fromCmd, toCmd)
}

func makeCmdID(cmd string) string {
	return fmt.Sprintf("cmd:%s", cmd)
}

func join(elements []string, sep string) string {
	return strings.Join(elements, sep)
}

func mapToCmdIDs(cmds ...string) []string {
	cmdIDs := make([]string, len(cmds))
	for i, c := range cmds {
		cmdIDs[i] = makeCmdID(c)
	}
	return cmdIDs
}

func parseLines(fileContent []byte) [][]string {
	lines := make([][]string, 0)
	tokenLines := bytes.Split(fileContent, []byte{'\n'})
	for _, line := range tokenLines {
		tokensPerLine := bytes.Split(line, []byte{'-', '>'})
		commandsPerLine := make([]string, 0)
		for _, token := range tokensPerLine {
			command := string(bytes.TrimSpace(token))
			if len(command) > 0 {
				commandsPerLine = append(commandsPerLine, command)
			}
		}
		if len(commandsPerLine) > 0 {
			lines = append(lines, commandsPerLine)
		}
	}
	return lines
}

func makeConnectedPairs(lines [][]string) [][]string {
	pairs := make([][]string, 0)
	contained := make(map[string]bool)
	addPair := func(fromCmd string, toCmd string) {
		key := fmt.Sprintf("%s:%s", fromCmd, toCmd)
		_, ok := contained[key]
		if !ok {
			pairs = append(pairs, []string{fromCmd, toCmd})
			contained[key] = true
		}
	}
	for _, line := range lines {
		for i, command := range line {
			if i > 0 {
				if i < (len(line) - 1) {
					// both previous and next tokens
					previousCommand := line[i-1]
					nextCommand := line[i+1]
					addPair(previousCommand, command)
					addPair(command, nextCommand)
				} else {
					// only previous token
					previousCommand := line[i-1]
					addPair(previousCommand, command)
				}
			} else {
				if len(line) > 1 {
					// only next token
					nextCommand := line[i+1]
					addPair(command, nextCommand)
				}
			}
		}
	}
	return pairs
}

func uniqueProcessNames(lines [][]string) []string {
	processes := make([]string, 0)
	contained := make(map[string]bool)
	addProcess := func(cmd string) {
		_, ok := contained[cmd]
		if !ok {
			processes = append(processes, cmd)
			contained[cmd] = true
		}
	}
	for _, line := range lines {
		for _, command := range line {
			addProcess(command)
		}
	}
	return processes
}

func pipeFromKey(fromCmd string) string {
	return makeFromKey("pipe", fromCmd)
}
func pipeToKey(toCmd string) string {
	return makeToKey("pipe", toCmd)
}

func forkFromKey(fromCmd string) string {
	return makeFromKey("fork", fromCmd)
}
func forkToKey(toCmd string) string {
	return makeToKey("fork", toCmd)
}

func mergeFromKey(fromCmd string) string {
	return makeFromKey("merge", fromCmd)
}
func mergeToKey(toCmd string) string {
	return makeToKey("merge", toCmd)
}

func makeFromKey(kind string, fromCmd string) string {
	return fmt.Sprintf("%s:%s->", kind, fromCmd)
}
func makeToKey(kind string, toCmd string) string {
	return fmt.Sprintf("%s:->%s", kind, toCmd)
}

func connectionsFrom(cmd string, connectedPairs [][]string) []string {
	acc := make([]string, 0)
	for _, pair := range connectedPairs {
		if pair[0] == cmd {
			acc = append(acc, pair[1])
		}
	}
	return acc
}
func connectionsTo(cmd string, connectedPairs [][]string) []string {
	acc := make([]string, 0)
	for _, pair := range connectedPairs {
		if pair[1] == cmd {
			acc = append(acc, pair[0])
		}
	}
	return acc
}

func findForkConnections(connectedPairs [][]string) map[string][]string {
	forkConnections := make(map[string][]string)
	for _, pair := range connectedPairs {
		fromCmd := pair[0]
		connectedTo := connectionsFrom(fromCmd, connectedPairs)
		if len(connectedTo) > 1 {
			forkConnections[fromCmd] = connectedTo
		}
	}
	return forkConnections
}

func findMergeConnections(connectedPairs [][]string) map[string][]string {
	mergeConnections := make(map[string][]string)
	for _, pair := range connectedPairs {
		toCmd := pair[1]
		connectedFrom := connectionsTo(toCmd, connectedPairs)
		if len(connectedFrom) > 1 {
			mergeConnections[toCmd] = connectedFrom
		}
	}
	return mergeConnections
}

func findPipeConnections(connectedPairs [][]string) map[string]string {
	pipeConnections := make(map[string]string)
	for _, pair := range connectedPairs {
		fromCmd := pair[0]
		toCmd := pair[1]
		connectedTo := connectionsFrom(fromCmd, connectedPairs)
		connectedFrom := connectionsTo(toCmd, connectedPairs)
		if len(connectedTo) == 1 && len(connectedFrom) == 1 {
			pipeConnections[fromCmd] = connectedTo[0]
		}
	}
	return pipeConnections
}
