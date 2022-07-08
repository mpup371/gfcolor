package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	yellow  = color.New(color.FgYellow)
	cyan    = color.New(color.FgCyan)
	names   = color.New(color.Faint)
	config  = color.New(color.FgHiBlue)
	info    = color.New(color.FgWhite, color.BgGreen)
	fine    = color.New(color.FgHiGreen)
	finest  = color.New(color.FgGreen)
	finer   = color.New(color.FgGreen)
	warning = color.New(color.FgMagenta)
	severe  = color.New(color.FgRed)
	debug   = color.New(color.BgBlue)
)

const RFCJava = "2006-01-02T15:04:05.999-0700"

func split(data []byte, atEOF bool) (int, []byte, error) {
	start := -1
	last := false
	for i, c := range data {
		if c == '[' && start == -1 {
			if i+1 < len(data) && data[i+1] == '[' {
				last = true
				start = i + 3
			} else {
				start = i + 1
			}
		}
		if c == ']' {
			if last && i+1 < len(data) && data[i+1] == ']' {
				return i + 3, data[start:i], nil
			} else if !last && start > -1 {
				return i + 1, data[start:i], nil
			}
		}
	}
	if atEOF {
		return 0, nil, bufio.ErrFinalToken
	}
	return 0, nil, nil
}

func init() {
	color.NoColor = false
}

func main() {
	dbg := flag.Bool("d", false, "mode debug")
	flag.Parse()
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(split)
	field := 0
	for scanner.Scan() {
		line := scanner.Text()
		if _, err := time.Parse(RFCJava, line); err == nil {
			yellow.Print(line)
			fmt.Print(" ")
			field = 1
			break
		}
	}
	if field == 0 {
		panic("cannot find start of line")
	}
	for scanner.Scan() {
		line := scanner.Text()
		if *dbg {
			debug.Println(line)
		}
		if _, err := time.Parse(RFCJava, line); err == nil {
			field = 0
		}
		switch field {
		case 0:
			yellow.Print(line)
			fmt.Print(" ")
		case 2:
			cl := sortie(line)
			cl.Print(line)
		case 4:
			fmt.Print(" ")
			cyan.Println(line)
		case 8, 9, 10:
			if strings.HasPrefix(line, "CLASSNAME:") {
				names.Println(line)
			} else if strings.HasPrefix(line, "METHODNAME:") {
				names.Println(line)
			} else {
				fmt.Println(line)
				field = -1
			}
		}
		field++
	}

	if scanner.Err() != nil {
		fmt.Println(scanner.Err())
	}
}

func sortie(level string) *color.Color {
	if level == "FINER" {
		return finer
	}
	if level == "FINEST" {
		return finest
	}
	if level == "FINE" {
		return fine
	}
	if level == "INFO" {
		return info
	}
	if level == "WARNING" {
		return warning
	}
	if level == "WARN" {
		return warning
	}
	if level == "SEVERE" {
		return severe
	}
	if level == "CONFIG" {
		return config
	}
	fmt.Println("level=", level)
	return nil
}
