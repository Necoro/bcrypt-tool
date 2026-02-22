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

type PwdInfo struct {
	TruncInput bool   `name:"truncate" help:"accept password input longer than ${maxPasswordLength} and truncate it"`
	Password   string `arg:"" optional:"" help:"read from stdin if omitted or '-'"`
}

type HashCmd struct {
	Cost    int `short:"c" help:"cost parameter for bcrypt" default:"${DEFAULT_COST}"`
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

type MatchCmd struct {
	Hash    string `arg:""`
	PwdInfo `embed:""`
}

func (cmd *MatchCmd) Run() error {
	pwd, err := cmd.GetPwd()
	if err != nil {
		return err
	}

	ok := match(pwd, cmd.Hash)
	if ok {
		fmt.Println("yes")
	} else {
		fmt.Println("no")
		os.Exit(1)
	}
	return nil
}

type CostCmd struct {
	Hash string `arg:""`
}

func (cmd *CostCmd) Run() error {
	c, e := cost(cmd.Hash)
	if e != nil {
		return e
	}
	fmt.Printf("%d\n", c)
	return nil
}

var CLI struct {
	Hash  HashCmd  `cmd:"" help:"generate a bcrypt hash"`
	Match MatchCmd `cmd:"" help:"check if password matches hash"`
	Cost  CostCmd  `cmd:"" help:"print the cost of the given hash"`
}

// needed to factor out for tests
func kongParserOptions() []kong.Option {
	return []kong.Option{
		kong.Description("A dandy CLI tool for generating and matching bcrypt hashes."),
		kong.UsageOnError(),
		kong.Vars{
			"DEFAULT_COST":      strconv.Itoa(bcrypt.DefaultCost),
			"maxPasswordLength": strconv.Itoa(maxPasswordLength),
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
		return nil, fmt.Errorf("no password provided")
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
