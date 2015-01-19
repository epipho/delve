package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/derekparker/delve/client/cli"
)

const version string = "0.3.0.beta"

func init() {
	// We must ensure here that we are running on the same thread during
	// the execution of dbg. This is due to the fact that ptrace(2) expects
	// all commands after PTRACE_ATTACH to come from the same thread.
	runtime.LockOSThread()
}

func main() {
	var (
		pid    int
		run    bool
		printv bool
		prompt string
	)

	flag.IntVar(&pid, "pid", 0, "Pid of running process to attach to.")
	flag.BoolVar(&run, "run", false, "Compile program and begin debug session.")
	flag.BoolVar(&printv, "v", false, "Print version number and exit.")
	flag.StringVar(&prompt, "prompt", "(dlv)", "Set the debugging prompt. Default (dlv)")
	flag.Parse()

	if flag.NFlag() == 0 && len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(0)
	}

	if printv {
		fmt.Printf("Delve version: %s\n", version)
		os.Exit(0)
	}

	cli.Run(run, pid, prompt, flag.Args())
}
