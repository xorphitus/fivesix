package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/xorphitus/fivesix/pkg/lock"
)

func stderr(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
}

func main() {
	args := os.Args
	if len(args) < 3 {
		stderr("Two arguments are required: </path/to/binary> <PID>\n")
		os.Exit(1)
	}

	binPath := args[1]
	pid := args[2]
	pidInt, err := strconv.Atoi(pid)
	if err != nil {
		stderr("PID %s must be integer: %v\n", pid, err)
		os.Exit(1)
	}

	table, done, err := lock.Run(binPath, pidInt)
	if err != nil {
		stderr("Runtime error: %v\n", pid, err)
		os.Exit(1)
	}
	defer done()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ticker.C:
			for it := table.Iter(); it.Next(); {
				v := binary.LittleEndian.Uint64(it.Leaf())
				fmt.Printf("%d\n", v)
			}
		case <-sig:
			os.Exit(1)
		}
	}

}
