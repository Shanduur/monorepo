package cproc

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

// CommandProcessor TODO
type CommandProcessor struct {
	Commands map[string]*Command

	Delimiter byte
	Prompt    string

	ExitHandler func() error
	exitChan    chan bool
	HelpCommand string

	Reader       *bufio.Reader
	OutputStream io.Writer
	ErrorStream  io.Writer
	Log          *log.Logger
}

// NewProcessor TODO
func NewProcessor(out, err io.Writer, in io.Reader, prompt string) *CommandProcessor {
	cp := &CommandProcessor{
		Commands:  make(map[string]*Command),
		Delimiter: '\n',
		Prompt:    prompt,
		ExitHandler: func() error {
			fmt.Fprintf(out, "Exiting...\n")
			os.Exit(0)

			return nil
		},
		exitChan: make(chan bool),

		Reader:       bufio.NewReader(in),
		OutputStream: out,
		ErrorStream:  err,

		Log: log.New(err, "", log.LstdFlags),
	}

	/*
	 */
	cp.AddCommand("exit", cp.exitFunc, "built-in", "exits currently running program", `long`).SetArgs(ArgIgnore)

	/*
	 */
	cp.AddCommand("help", cp.helpFunc, "built-in", "displays basic help prompt", `long`).SetArgs(ArgString)

	/*
	 */
	cp.AddCommand("set-env", cp.setEnvFunc, "built-in", "sets environmental variables", `long`).SetArgs(ArgNString)

	/*
	 */
	cp.AddCommand("get-env", cp.getEnvFunc, "built-in", "returns list of values of requested environmental variables", `long`).SetArgs(ArgNStringOpt)

	/*
	 */
	cp.AddCommand("get-pid", cp.getPIDFunc, "built-in", "returns PID of current program", `long`).SetArgs(ArgIgnore)

	return cp
}

// AddCommand TODO
func (cp *CommandProcessor) AddCommand(name string, action func(...string), category, short, long string) *Command {
	cmd := &Command{
		Name:             name,
		Category:         category,
		Args:             ArgNStringOpt,
		ShortDescription: short,
		LongDescription:  long,
		Action:           action,
		Processor:        cp,
	}

	cp.Commands[name] = cmd

	return cmd
}

// Run TODO
func (cp *CommandProcessor) Run() error {
	for {
		select {
		case <-cp.exitChan:
			return nil

		default:
			ins, err := cp.getNext()
			if errors.Is(err, ErrNoCommandProvided) {
				continue
			} else if errors.Is(err, io.EOF) {
				continue
			} else if err != nil {
				cp.Log.Printf("unable to get next command: %s", err.Error())
				continue
			}

			cmd := cp.Commands[ins[0]]
			if cmd != nil {
				cmd.Run(ins[1:]...)
			} else {
				cp.Printf("command was not recognized\n")
			}
		}
	}
}

func (cp CommandProcessor) getNext() ([]string, error) {
	cp.Printf(cp.Prompt)
	in, err := cp.Reader.ReadString(cp.Delimiter)
	if err != nil {
		return nil, fmt.Errorf("unable to read: %s", err.Error())
	}

	in = strings.ReplaceAll(in, "\n", "")
	in = strings.ReplaceAll(in, "\r", "")

	ins := strings.Split(in, " ")
	if len(ins[0]) < 1 {
		return nil, ErrNoCommandProvided
	}

	return ins, nil
}

func (cp CommandProcessor) commandError() {
	cp.Log.Printf("Command unsupported!")
}

// Printf TODO
func (cp CommandProcessor) Printf(format string, a ...interface{}) {
	fmt.Fprint(cp.OutputStream, fmt.Sprintf(format, a...))
}

// Logf TODO
func (cp CommandProcessor) Logf(format string, a ...interface{}) {
	cp.Log.Printf(format, a...)
}

// Fatalf TODO
func (cp CommandProcessor) Fatalf(format string, a ...interface{}) {
	cp.Log.Fatalf(format, a...)
}

func (cp *CommandProcessor) exitFunc(s ...string) {
	if err := cp.ExitHandler(); err != nil {
		cp.Log.Fatalf("error ocurred during exit: %s", err.Error())
	}
	go func() {
		cp.exitChan <- true
	}()
	time.Sleep(5 * time.Millisecond)
}

func (cp CommandProcessor) helpFunc(s ...string) {
	if len(s) == 0 {
		cnf := make(map[string][]string)
		var c []string

		for k, v := range cp.Commands {
			cnf[v.Category] = append(cnf[v.Category], k)
		}

		for k := range cnf {
			c = append(c, k)
		}

		sort.Strings(c)
		for _, k := range c {
			if k == "built-in" {
				continue
			}

			if len(k) > 0 {
				cp.Printf("%s\n", k)
			}

			for _, v := range cnf[k] {
				cp.Printf("\t%s %s\t%s\n", v, cp.Commands[v].Args, cp.Commands[v].Brief())
			}

			cp.Printf("\n")
		}

		cp.Printf("built-in\n")

		for _, v := range cnf["built-in"] {
			cp.Printf("\t%s %s\t%s\n", v, cp.Commands[v].Args, cp.Commands[v].Brief())
		}
	} else {
		cmd := cp.Commands[s[0]]
		cp.Printf("%s %s\t%s\n\n\t%s\n", cmd.Name, cmd.Args, cmd.Brief(), cmd.Help())
	}
}

func (cp CommandProcessor) setEnvFunc(s ...string) {
	if len(s) == 0 {
		cp.Printf("no Key=Value pair provided\n")
	} else {
		for _, kv := range s {
			skv := strings.SplitN(kv, "=", 2)
			if len(skv) < 2 {
				cp.Printf("Key=Value pair improper: %s\n", kv)
				continue
			} else {
				os.Setenv(skv[0], skv[1])
			}
		}
	}
}

func (cp CommandProcessor) getEnvFunc(s ...string) {
	for i, k := range s {
		v := os.Getenv(k)
		cp.Printf("%d: %s=%s\n", i, k, v)
	}
}

func (cp CommandProcessor) getPIDFunc(s ...string) {
	cp.Printf("PID: %d\n", os.Getpid())
}
