package main

import (
	"flag"
	"strings"
	"os/exec"
	"fmt"
	"os"
	"path/filepath"
	"log"
	"io"
)

const logFile = "log.txt"
const preCompiledXmlDir = "C:\\Users\\Adam\\Desktop\\FS2XPlane"

func main() {
	f, err := os.OpenFile(logFile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	log.SetOutput(f)

	_ = flag.Bool("t", true, "Terse mode")
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		log.Println("Not enough arguments.")
		os.Exit(1)
	}

	err = run(args[0], args[1])
	if err != nil {
		// run will have printed the error itself, so we'll just exit.
		os.Exit(1)
	}
}

func run(src, dst string) error {
	var err error
	ls := fmt.Sprintf("Converting %q -> %q: ", src, dst)

	if c := getPreCompiledXMLPath(preCompiledXmlDir, src); c != "" {
		Copy(c, dst)
		ls += fmt.Sprintf("Replaced with %q", c)
	} else {
		err = convert(src, dst)

		if err != nil {
			ls += fmt.Sprintf("Failed: %v", err)
		} else {
			ls += "Successful."
		}
	}

	log.Println(ls)
	return err
}

func convert(src, dst string) error {
	cmd := "bglxml-actual"
	args := []string{"-t", src, dst}
	return exec.Command(cmd, args...).Run()
}

// getPreCompiledXMLPath returns the absolute file path to a pre-compiled XML
// file in dir that matches the source file src.
// Returns an empty string if a matching pre-compiled XML could not be found.
func getPreCompiledXMLPath(dir, src string) string {
	base := filepath.Base(src)
	name := strings.TrimSuffix(base, filepath.Ext(base)) + ".xml"
	compiledPath := filepath.Join(dir, name)

	_, err := os.Stat(compiledPath)
	if err != nil {
		return ""
	}
	return compiledPath
}

func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
