package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	args := os.Args[1:]
	argsLen := len(args)

	if argsLen == 0 {
		printHeader()
		printHelp()
		return
	}

	switch args[0] {
	case "tap":
		if argsLen < 2 {
			fmt.Println("Not enough arguments for tap:")
			printTapHelp()
		} else {
			startTap(args[1])
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

func open(filename string) (*os.File, error) {
	err := syscall.Mkfifo(filename, 0666)
	if err != nil {
		return nil, err
	}
	return os.OpenFile(filename, os.O_RDWR, os.ModeNamedPipe)
}

func printHeader() {
	fmt.Println("hub")
	fmt.Println("Utility to build messaging systems by composing command line processes")
}

func printHelp() {
	printTapHelp()
	printForkHelp()
	printMergeHelp()
	printCleanupHelp()
}
func printTapHelp() {
	fmt.Println("tap <name>")
	fmt.Println("  Tap into message stream from <name>.wr to <name>.rd")
}
func printForkHelp() {
	fmt.Println("fork <wr-name> <rd-name-1> <rd-name-2> <...>")
	fmt.Println("  Fork message stream from wr-fifo into all provided rd-fifos")
}
func printMergeHelp() {
	fmt.Println("merge <wr-name-1> <wr-name-2> <...< <rd-name>")
	fmt.Println("  Merge message stream from all provided wr-fifos into rd-fifo")
}
func printCleanupHelp() {
	fmt.Println("cleanup <directory>")
	fmt.Println("  Cleanup directory from fifos (wr & rd)")
}
