package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/peteraba/fakie/lib"
	"github.com/tj/docopt"
)

var usage = `
  Usage: fakie
    [--tick d]
    [--max n]
    [--batch n]
    [--list]
    [--concurrent]

    fakie -h | --help
    fakie -v | --version

  Options:
    --tick d        generate data every d [default: 10ms]
    --max n         generate data up to n [default: -1]
    --batch n       batch size for concurrent runs [default: 0]
    --list          list all available generators
    --concurrent    skip ticks and generate fake data concurrently
    -v, --version   show version information
    -h, --help      show help information

`

func main() {
	args, err := docopt.Parse(usage, nil, true, "0.0.2", false)
	check(err)

	g := fakie.NewGenerator()

	if args["--list"].(bool) {
		all := g.List()
		sort.Strings(all)
		println()
		for _, name := range all {
			fmt.Printf("  %s\n", name)
		}
		println()
		os.Exit(0)
	}

	rand.Seed(time.Now().UnixNano())

	d := parseDuration(args["--tick"].(string))
	max := parseInt(args["--max"].(string))
	batch := parseInt(args["--batch"].(string))
	tmpl := readAll(os.Stdin)
	tick := time.Tick(d)
	f := compile(string(tmpl), g)

	if args["--concurrent"].(bool) {
		concurrentMain(max, tmpl, f, batch)
	} else {
		tickMain(max, tmpl, f, tick)
	}
}

func concurrentMain(max int, tmpl string, f func() string, batchSize int) {
	it := 0
	c := make(chan string)
	closure := func(c chan string) {
		c <- f()
	}

	if batchSize == 0 {
		batchSize = 100
	}

	for {
		j := 0
		for ; j < batchSize; j++ {
			go closure(c)
			if it++; -1 != max && it == max {
				break
			}
		}

		for ; j > 0; j-- {
			fmt.Printf(<-c)
		}

		if -1 != max && it == max {
			break
		}
	}

}

func tickMain(max int, tmpl string, f func() string, tick <-chan time.Time) {
	it := 0

	for range tick {
		fmt.Fprintf(os.Stdout, "%s", f())
		if it++; -1 != max && it == max {
			return
		}
	}
}

func compile(tmpl string, g *fakie.Generator) func() string {
	expr, err := regexp.Compile(`({{ *(([a-zA-Z0-9]+(\.[a-zA-Z0-9]+)?)+(\:([a-zA-Z0-9\.,-]+))?) *}})`)
	check(err)

	return func() string {
		var dataCache []string
		var r *rand.Rand = fakie.CreateRand()

		return expr.ReplaceAllStringFunc(tmpl, func(s string) string {
			var data string

			call := strings.Trim(s[2:len(s)-2], " ")

			parts := strings.Split(call, ":")
			var arguments []string
			if len(parts) == 2 {
				arguments = strings.Split(parts[1], ",")
			}

			i64, err := strconv.ParseInt(parts[0], 10, 64)

			if err != nil {
				data, err = g.GetWithArgs(parts[0], arguments, r)
				check(err)

				dataCache = append(dataCache, data)
			} else {
				if len(dataCache) <= int(i64) {
					check(errors.New("Given template references a non-existent value"))
				}
				return dataCache[i64]
			}

			return data
		})
	}
}

func check(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "fakie: %s\n", err.Error())
		os.Exit(1)
	}
}

func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	check(err)
	return i
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	check(err)
	return d
}

func readAll(r *os.File) string {
	b, err := ioutil.ReadAll(r)
	check(err)
	return string(b)
}
