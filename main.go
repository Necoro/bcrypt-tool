// Command bcrypt-tool is command line tool for common bcrypt functions
// including the ability to generate hashes, determine if a password
// matches a hash, and compute cost from a hash.
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/alecthomas/kong"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

// bcrypt handles at most 72 characters of password
// unfortunately, there is no constant in the bcrypt package
const maxPasswordLength = 72

const (
	NAME    = "bcrypt-tool"
	VERSION = "2.0.0"
)

type PwdInfo struct {
	TruncInput bool   `name:"truncate" help:"Accept password input longer than ${maxPasswordLength} and truncate it."`
	Password   string `arg:"" optional:"" help:"Read from stdin if omitted or '-'."`
}

type HashCmd struct {
	Cost    int `short:"c" help:"Cost parameter for bcrypt." default:"${DEFAULT_COST}"`
	PwdInfo `embed:""`
}

func (cmd *HashCmd) Validate() error {
	if cmd.Cost < bcrypt.MinCost || cmd.Cost > bcrypt.MaxCost {
		return fmt.Errorf("cost must be between %d and %d", bcrypt.MinCost, bcrypt.MaxCost)
	}
	return nil
}

func (cmd *HashCmd) Run() error {
	pwd, err := cmd.GetPwd()
	if err != nil {
		return err
	}

	h, err := hash(pwd, cmd.Cost)
	if err != nil {
		return err
	}

	fmt.Println(h)
	return nil
}

type Hash string

func (h Hash) Validate() error {
	if len(h) == 0 {
		return errors.New("no hash given")
	}
	if h[0] != '$' {
		return errors.New("does not look like a bcrypt hash")
	}

	return nil
}

type MatchCmd struct {
	Quiet   bool `short:"q" help:"Omit printing the result."`
	Hash    Hash `arg:""`
	PwdInfo `embed:""`
}

func (cmd *MatchCmd) print(s string) {
	if !cmd.Quiet {
		fmt.Println(s)
	}
}

func (cmd *MatchCmd) Run() error {
	pwd, err := cmd.GetPwd()
	if err != nil {
		return err
	}

	ok := match(pwd, string(cmd.Hash))
	if ok {
		cmd.print("yes")
	} else {
		cmd.print("no")
		os.Exit(1)
	}
	return nil
}

type CostCmd struct {
	Hash Hash `arg:""`
}

func (cmd *CostCmd) Run() error {
	c, e := cost(string(cmd.Hash))
	if e != nil {
		return e
	}
	fmt.Printf("%d\n", c)
	return nil
}

var CLI struct {
	Hash    HashCmd          `cmd:"" help:"Generate a bcrypt hash."`
	Match   MatchCmd         `cmd:"" help:"Check if password matches hash."`
	Cost    CostCmd          `cmd:"" help:"Print the cost of the given hash."`
	Version kong.VersionFlag `help:"Print the current version (${version})."`
}

// needed to factor out for tests
func kongParserOptions() []kong.Option {
	return []kong.Option{
		kong.Name(NAME),
		kong.Description("A dandy CLI tool for generating and matching bcrypt hashes."),
		kong.UsageOnError(),
		kong.Vars{
			"DEFAULT_COST":      strconv.Itoa(bcrypt.DefaultCost),
			"maxPasswordLength": strconv.Itoa(maxPasswordLength),
			"version":           NAME + " v" + VERSION,
		},
	}
}

func main() {
	ctx := kong.Parse(&CLI, kongParserOptions()...)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}

func (p *PwdInfo) GetPwd() ([]byte, error) {
	if p.Password == "" || p.Password == "-" {
		return p.readFromStdin()
	}
	return []byte(p.Password), nil
}

func (p *PwdInfo) readFromStdin() ([]byte, error) {
	fd := int(os.Stdin.Fd())

	// we have a terminal, don't treat as a pipe
	if term.IsTerminal(fd) {
		fmt.Print("Password: ")
		password, err := term.ReadPassword(fd)
		fmt.Println("")

		if err != nil {
			return nil, fmt.Errorf("reading password: %w", err)
		}
		return password, nil
	}

	maxLen := maxPasswordLength + 1
	if p.TruncInput {
		maxLen = maxPasswordLength
	}

	buf := make([]byte, maxLen)
	n, err := io.ReadFull(os.Stdin, buf)
	if errors.Is(err, io.EOF) || buf[0] == '\n' || buf[0] == 0 {
		return nil, errors.New("no password provided")
	} else if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return nil, fmt.Errorf("reading password: %w", err)
	} else if n > maxPasswordLength {
		//goland:noinspection GoErrorStringFormat
		return nil, fmt.Errorf("bcrypt supports an input of at most %d bytes."+
			" Use --truncate if you want to truncate the input to that length.", maxPasswordLength)
	}

	// remove trailing '\n'
	if buf[n-1] == '\n' {
		n--
	}

	return buf[:n], nil
}

func cost(hash string) (int, error) {
	h := []byte(hash)
	return bcrypt.Cost(h)
}

func match(password []byte, hash string) bool {
	h := []byte(hash)
	e := bcrypt.CompareHashAndPassword(h, password)
	return e == nil
}

func hash(password []byte, cost int) (string, error) {
	h, e := bcrypt.GenerateFromPassword(password, cost)
	if e != nil {
		return "", e
	}
	return string(h), nil
}
