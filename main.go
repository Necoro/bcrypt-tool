// Command bcrypt-tool is command line tool for common bcrypt functions
// including the ability to generate hashes, determine if a password
// matches a hash, and compute cost from a hash.
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/alecthomas/kong"
	"golang.org/x/crypto/bcrypt"
)

type HashCmd struct {
	Cost     int    `short:"c" help:"Cost parameter for bcrypt." default:"${DEFAULT_COST}"`
	Password string `arg:"" optional:"" help:"read from stdin if omitted or '-'"`
}

func (cmd *HashCmd) Run() error {
	fmt.Println(hash(cmd.Password, cmd.Cost))
	return nil
}

type MatchCmd struct {
	Hash     string `arg:""`
	Password string `arg:"" optional:"" help:"read from stdin if omitted or '-'"`
}

func (cmd *MatchCmd) Run() error {
	ok := match(cmd.Hash, cmd.Password)
	if ok {
		fmt.Println("yes")
	} else {
		fmt.Println("no")
		os.Exit(1)
	}
	return nil
}

type CostCmd struct {
	Hash string `arg:"" help:"read from stdin if omitted or '-'"`
}

func (cmd *CostCmd) Run() error {
	c := cost(cmd.Hash)
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
		kong.Vars{
			"DEFAULT_COST": strconv.Itoa(bcrypt.DefaultCost),
		},
	}
}

func main() {
	ctx := kong.Parse(&CLI, kongParserOptions()...)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}

func cost(hash string) int {
	h := []byte(hash)
	c, e := bcrypt.Cost(h)
	if e != nil {
		_, _ = fmt.Fprintln(os.Stderr, e)
		os.Exit(2)
	}
	return c
}

func match(password, hash string) bool {
	p := []byte(password)
	h := []byte(hash)
	e := bcrypt.CompareHashAndPassword(h, p)
	if e != nil {
		return false
	}
	return true
}

func hash(password string, cost int) string {
	p := []byte(password)
	h, e := bcrypt.GenerateFromPassword(p, cost)
	if e != nil {
		_, _ = fmt.Fprintln(os.Stderr, e)
		os.Exit(2)
	}
	return string(h)
}
