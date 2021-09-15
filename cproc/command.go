package cproc

import "fmt"

var (
	// ErrCommandNotRecognized TODO
	ErrCommandNotRecognized = fmt.Errorf("provided command was not recognized")

	// ErrNoCommandProvided TODO
	ErrNoCommandProvided = fmt.Errorf("no command was provided")
)

type argsTemplate string

const (
	// ArgIgnore TODO
	ArgIgnore = argsTemplate("")
	// ArgSubcommand TODO
	ArgSubcommand = argsTemplate("cmd")
	// ArgString TODO
	ArgString = argsTemplate("str")
	// ArgStringOpt TODO
	ArgStringOpt = argsTemplate("[str]")
	// ArgNString TODO
	ArgNString = argsTemplate("str ...")
	// ArgNStringOpt TODO
	ArgNStringOpt = argsTemplate("[str] ...")
	// ArgInt TODO
	ArgInt = argsTemplate("1")
	// ArgIntOpt TODO
	ArgIntOpt = argsTemplate("[1]")
	// ArgNInt TODO
	ArgNInt = argsTemplate("1 ...")
	// ArgNIntOpt TODO
	ArgNIntOpt = argsTemplate("[1] ...")
)

// Command TODO
type Command struct {
	Name             string
	Category         string
	ShortDescription string
	LongDescription  string
	Args             argsTemplate

	Processor   *CommandProcessor
	Action      func(...string)
	SubCommands map[string]*Command
}

// AddSubcommand TODO
func (c *Command) AddSubcommand(name string, action func(...string), short, long string) *Command {
	cmd := &Command{
		Name:             name,
		Category:         c.Category,
		ShortDescription: short,
		LongDescription:  long,
		Args:             ArgNStringOpt,

		Action:    action,
		Processor: c.Processor,
	}

	c.SubCommands[name] = cmd

	return cmd
}

// SetArgs TODO
func (c *Command) SetArgs(args argsTemplate) *Command {
	c.Args = args

	return c
}

// Run TODO
func (c Command) Run(s ...string) {
	// TODO: validate arguments

	c.Action(s...)
}

// Brief TODO
func (c Command) Brief() string {
	return c.ShortDescription
}

// Help TODO
func (c Command) Help() string {
	return c.LongDescription
}
