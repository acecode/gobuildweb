package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/agtorre/gocolorize"
)

type ProjectConfig struct {
	Package PackageConfig
}

type PackageConfig struct {
	Name         string
	Version      string
	Authors      []string
	Dependencies []string `toml:"deps"`
	OmitTests    []string `toml:"omit_tests"`
}

func usage() {
	fmt.Println("Build a Golang web application")
	fmt.Println("Usage:")
	fmt.Println("  update    Update all your dependencies...")
	fmt.Println("  build     Build your web application")
	fmt.Println("  run       Will watch your file changes and run the application")
	os.Exit(1)
}

func main() {
	cmds := map[string]Command{
		"run":    commandRun,
		"dist":   commandDist,
		"update": commandUpdate,
	}
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		usage()
	}
	if cmd, ok := cmds[args[0]]; !ok {
		usage()
	} else {
		if fi, err := os.Stat("project.toml"); os.IsNotExist(err) {
			ERROR.Fatalf("Please provide a project.toml for web project.")
		} else if err != nil {
			ERROR.Fatalf("Accessing project.toml file error, %v.", err)
		} else if fi.IsDir() {
			ERROR.Fatalf("project.toml cannot be a directory.")
		}

		if _, err := toml.DecodeFile("project.toml", &rootConfig); err != nil {
			ERROR.Fatalf("Cannot decode the project.toml into TOML format, %v", err)
		}
		INFO.SetPrefix(fmt.Sprintf("[%s][INFO] ", rootConfig.Package.Name))
		SUCC.SetPrefix(fmt.Sprintf("[%s][SUCC] ", rootConfig.Package.Name))
		WARN.SetPrefix(fmt.Sprintf("[%s][WARN] ", rootConfig.Package.Name))
		ERROR.SetPrefix(fmt.Sprintf("[%s][ERROR] ", rootConfig.Package.Name))

		SUCC.Printf("Loaded project.toml...")
		if err := cmd(args[1:]); err != nil {
			ERROR.Fatalf("Executing command [%v] error, %v", args[0], err)
		}
	}
}

type ColoredLogger struct {
	c gocolorize.Colorize
	w io.Writer
}

func (cl *ColoredLogger) Write(p []byte) (n int, err error) {
	return cl.w.Write([]byte(cl.c.Paint(string(p))))
}

var rootConfig ProjectConfig
var (
	INFO  *log.Logger
	SUCC  *log.Logger
	WARN  *log.Logger
	ERROR *log.Logger
)

func init() {
	INFO = log.New(os.Stdout, "[INFO] ", log.LstdFlags)
	SUCC = log.New(&ColoredLogger{gocolorize.NewColor("green"), os.Stdout}, "[SUCC] ", log.LstdFlags)
	WARN = log.New(&ColoredLogger{gocolorize.NewColor("yellow"), os.Stdout}, "[WARN] ", log.LstdFlags)
	ERROR = log.New(&ColoredLogger{gocolorize.NewColor("red"), os.Stdout}, "[ERROR] ", log.LstdFlags)

	runtime.GOMAXPROCS(runtime.NumCPU())
}
