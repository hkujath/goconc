package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"time"
)

const (
	cmdDelimiter     = "::"
	cmdArgsDelimiter = ":"
)

type cmdArgs struct {
	name string
	args []string
}

var (
	flagTimeout = flag.Duration("t", time.Second, "timeout in seconds  (e.g. 5s for 5 seconds)")
	flagOutput  = flag.Bool("o", false, "enable output of the commands")
)

func init() {

}

// main function of this program
// Function can be called from command line with "go run .\main.go cmd /c echo hallo holger :: cmd /c dir"
func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if *flagTimeout != 0 {
		ctx, cancel = context.WithTimeout(ctx, *flagTimeout)
	}

	c := make(chan os.Signal, 1)
	go func() {
		signal.Notify(c, os.Interrupt)
		<-c
		log.Println("Abbort wit Ctrl+C")
		cancel()
	}()

	log.Println("Start parsing...")
	cmds := ParseArgs(os.Args[1:])

	log.Println("Start running commands..")
	RunCmds(ctx, cmds)
}

// ParseArgs parses a given list of strings
func ParseArgs(args []string) []cmdArgs {

	if len(args) == 0 {
		log.Printf("No arguments were given.\n")
		return nil
	}

	var cmds []cmdArgs
	var cmd cmdArgs

	for _, arg := range args {
		switch {
		case arg == cmdDelimiter:
			cmds = append(cmds, cmd)
			cmd = cmdArgs{}
		case arg == cmdArgsDelimiter:
			newCmd := cmdArgs{name: cmd.name}
			cmds = append(cmds, cmd)
			cmd = newCmd
		case cmd.name == "":
			cmd.name = arg
		default:
			cmd.args = append(cmd.args, arg)
		}
	}
	cmds = append(cmds, cmd)
	return cmds
}

func RunCmds(ctx context.Context, cmds []cmdArgs) {
	wg := &sync.WaitGroup{}

	for i, args := range cmds {
		wg.Add(1)
		go func(nr int, args cmdArgs) {
			log.Printf("Starting cmd %d: %s %s\n", nr, args.name, strings.Join(args.args, " "))
			cmd := exec.CommandContext(ctx, args.name, args.args...)
			if *flagOutput {
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
			}
			err := cmd.Run()
			log.Printf("Cmd %d is ready\n", nr)
			if err != nil {
				log.Printf("Error happened during execution of %s (arguments: %s) \n%s", args.name, args.args, err)
			}

			wg.Done()
		}(i+1, args)
	}
	wg.Wait()
}
