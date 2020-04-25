package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	argsLen := len(args)

	if argsLen == 0 {
		printHeader()
		printHelp()
		return
	} else if argsLen == 1 {
		coll := NewNodeCollection(args[0])
		for _, nodeID := range coll.IDs() {
			node, _ := coll.Node(nodeID)
			for _, nextNodeID := range node.OutKeys {
				_, nextNodeExists := coll.Node(nextNodeID)
				var found string
				if nextNodeExists {
					found = "FOUND"
				} else {
					found = "NOT FOUND"
				}
				log.Printf("%s -> %s [%s]\n", nodeID, nextNodeID, found)
			}
		}
		// runNodeCollection
		return
	}

	switch args[0] {
	case "pipe":
		if argsLen < 2 {
			fmt.Println("Not enough arguments for pipe:")
			printPipeHelp()
		} else {
			startPipe(args[1])
		}
	case "fork":
		if argsLen < 4 {
			fmt.Println("Not enough arguments for fork:")
			printForkHelp()
		} else {
			startFork(args[1], args[2:])
		}
	case "merge":
		if argsLen < 4 {
			fmt.Println("Not enough arguments for merge:")
			printMergeHelp()
		} else {
			lastIdx := len(args) - 1
			startMerge(args[1:lastIdx], args[lastIdx])
		}
	case "input":
		if argsLen < 2 {
			fmt.Println("Not enough arguments for input:")
			printInputHelp()
		} else {
			startInput(args[1])
		}
	case "cleanup":
		if argsLen < 2 {
			startCleanup(".")
		} else {
			startCleanup(args[1])
		}
	default:
		printHeader()
		printHelp()
	}
}

// func test() {
// 	src := strings.NewReader(`eyJuYW1lIjoicmVxdWVzdCIsInBheWxvYWQiOnsibWV0aG9kIjoiR0VUIn1
// 9
// ---
// eyJuYW1lIjoicmVxdWVzdCIsInBheWxvYWQiOnsibWV0aG9kIjoiUE9TVCI
// sImRhdGEiOnsiaGVsbG8iOiJ3b3JsZCJ9fX0=
// ---
// `)

// 	reader := NewMsgReader(src)
// 	for {
// 		msgB64, err := reader.ReadMessage()
// 		if err != nil {
// 			log.Println(err)
// 			break
// 		}
// 		if len(msgB64) > 0 {
// 			log.Printf("Did read message: %s\n", msgB64)
// 		}
// 	}

// }

func printHeader() {
	fmt.Println("hub")
	fmt.Println("Utility to build messaging systems by composing command line processes")
}

func printHelp() {
	printPipeHelp()
	printForkHelp()
	printMergeHelp()
	printInputHelp()
	printCleanupHelp()
}
func printPipeHelp() {
	fmt.Println("pipe <name>")
	fmt.Println("  Pipe message stream from <name>.wr to <name>.rd")
}
func printForkHelp() {
	fmt.Println("fork <wr-name> <rd-name-1> <rd-name-2> <...>")
	fmt.Println("  Fork message stream from wr-fifo into all provided rd-fifos")
}
func printMergeHelp() {
	fmt.Println("merge <wr-name-1> <wr-name-2> <...< <rd-name>")
	fmt.Println("  Merge message stream from all provided wr-fifos into rd-fifo")
}
func printInputHelp() {
	fmt.Println("input <name>")
	fmt.Println("  Input JSON messages to stream them to <name>")
}
func printCleanupHelp() {
	fmt.Println("cleanup <directory>")
	fmt.Println("  Cleanup directory from fifos (wr & rd)")
}
