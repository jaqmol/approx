package main

import (
	"log"
	"os"
	"os/exec"
)

func runNodeCollection(collection NodeCollection) {
	index := make(map[string]Runner)
	findOutputRunners := func(r Runner) []Runner {
		nodes := collection.Outputs(r.Node())
		runners := make([]Runner, len(nodes))
		for i, node := range nodes {
			runners[i] = index[node.ID]
		}
		return runners
	}

	errors := make(chan []byte)

	for _, id := range collection.IDs() {
		n, _ := collection.Node(id)
		switch n.Class {
		case ForkClass:
			fallthrough
		case MergeClass:
			fallthrough
		case PipeClass:
			i := InfrastructureRunner{
				node:              n,
				findOutputRunners: findOutputRunners,
				channel:           make(chan []byte),
			}
			index[n.ID] = i
		case ProcessClass:
			c := Process{
				node:              n,
				findOutputRunners: findOutputRunners,
				channel:           make(chan []byte),
				errors:            errors,
			}
			index[n.ID] = c
		}
	}

	for _, r := range index {
		switch r.Node().Class {
		case ForkClass:
			fallthrough
		case MergeClass:
			fallthrough
		case PipeClass:
			go r.Start()
		}
	}

	for _, r := range index {
		switch r.Node().Class {
		case ProcessClass:
			go r.Start()
		}
	}

	for message := range errors {
		message := append(message, '\n')
		os.Stderr.Write(message)
	}
}

// Runner ...
type Runner interface {
	Node() Node
	Input() chan<- []byte
	Start()
}

// InfrastructureRunner ...
type InfrastructureRunner struct {
	node              Node
	findOutputRunners func(Runner) []Runner
	channel           chan []byte
}

// Node ...
func (i InfrastructureRunner) Node() Node { return i.node }

// Start ...
func (i InfrastructureRunner) Start() {
	dispatchChannels := collectInputChannels(i.findOutputRunners(i))

	for message := range i.channel {
		for _, c := range dispatchChannels {
			c <- message
		}
	}
}

// Input ...
func (i InfrastructureRunner) Input() chan<- []byte { return i.channel }

// Process ...
type Process struct {
	node              Node
	findOutputRunners func(Runner) []Runner
	channel           chan []byte
	errors            chan<- []byte
}

// Node ...
func (p Process) Node() Node { return p.node }

// Start ...
func (p Process) Start() {
	cmd := exec.Command(p.node.Command)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		cmdErr, err := cmd.StderrPipe()
		if err != nil {
			log.Fatal(err)
		}
		scanner := NewHeavyDutyScanner(cmdErr, []byte{'\n'})

		for scanner.Scan() {
			message := scanner.Message()
			p.errors <- message
		}
	}()

	go func() {
		dispatchChannels := collectInputChannels(p.findOutputRunners(p))
		cmdOut, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		scanner := NewMsgScanner(cmdOut)

		for scanner.Scan() {
			for _, c := range dispatchChannels {
				c <- scanner.DelimitedMessage()
			}
		}
	}()

	cmdIn, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	for message := range p.channel {
		cmdIn.Write(message)
	}
}

// Input ...
func (p Process) Input() chan<- []byte { return p.channel }

func collectInputChannels(runners []Runner) []chan<- []byte {
	channels := make([]chan<- []byte, len(runners))
	for i, r := range runners {
		channels[i] = r.Input()
	}
	return channels
}

// func startRunnerLoop(
// 	node Node,
// 	runnersForNodeOutputNodes func(Node) []Runner,
// 	channel chan []byte,
// ) {
// 	dispatchChannels := collectInputChannels(runnersForNodeOutputNodes(node))

// 	for message := range channel {
// 		for _, c := range dispatchChannels {
// 			c <- message
// 		}
// 	}
// }

// // Fork ...
// type Fork struct {
// 	node                      Node
// 	runnersForNodeOutputNodes func(Node) []Runner
// 	channel                   chan []byte
// }

// func (f Fork) start() {
// 	startRunnerLoop(f.node, f.runnersForNodeOutputNodes, f.channel)
// }

// func (f Fork) input() chan<- []byte {
// 	return f.channel
// }

// // Merge ...
// type Merge struct {
// 	node                      Node
// 	runnersForNodeOutputNodes func(Node) []Runner
// 	channel                   chan []byte
// }

// func (m Merge) start() {
// 	startRunnerLoop(m.node, m.runnersForNodeOutputNodes, m.channel)
// }

// func (m Merge) input() chan<- []byte {
// 	return m.channel
// }

// // Pipe ...
// type Pipe struct {
// 	node                      Node
// 	runnersForNodeOutputNodes func(Node) []Runner
// 	channel                   chan []byte
// }

// func (p Pipe) start() {
// 	startRunnerLoop(p.node, p.runnersForNodeOutputNodes, p.channel)
// }

// func (p Pipe) input() chan<- []byte {
// 	return p.channel
// }
