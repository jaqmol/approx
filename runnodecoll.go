package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func validateNodeCollection(collection NodeCollection) error {
	makeError := func(node Node) error {
		return fmt.Errorf(
			"Error interpreting node if type \"%s\" [%s]",
			NodeClassToString(node.Class),
			node.ID,
		)
	}
	for _, nodeID := range collection.IDs() {
		node, nodeExists := collection.Node(nodeID)
		if !nodeExists {
			return makeError(node)
		}
		for _, nextNodeID := range node.OutKeys {
			nextNode, nextNodeExists := collection.Node(nextNodeID)
			if !nextNodeExists {
				return makeError(nextNode)
			}
		}
	}
	return nil
}

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
		printLogLn(string(message))
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
	cmd := exec.Command(p.node.Process.Command, p.node.Process.Arguments...)
	cmdErr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmdOut, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmdIn, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		scanner := bufio.NewScanner(cmdErr)
		for scanner.Scan() {
			message := scanner.Bytes()
			p.errors <- message
		}
	}()

	go func() {
		dispatchChannels := collectInputChannels(p.findOutputRunners(p))
		scanner := NewHeavyDutyScanner(cmdOut, MsgDelimiter)
		// scanner.Decode = DecodeBase64Message NOT NEEDED DecodeMessage never called
		for scanner.Scan() {
			for _, c := range dispatchChannels {
				c <- scanner.DelimitedMessage()
			}
		}
	}()

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	printLogLn(fmt.Sprintf(
		"Did start: %s %s",
		p.node.Process.Command,
		strings.Join(p.node.Process.Arguments, ", "),
	))

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
