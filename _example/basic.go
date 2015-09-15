package main

import (
	"fmt"
	"io"
	"os"

	"github.com/andrew-d/afterarch"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <read|write> [args...]\n", os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	if os.Args[1] == "read" {
		doRead()
	} else if os.Args[1] == "write" {
		doWrite()
	} else {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		usage()
	}
}

func doRead() {
	f, err := os.Open(os.Args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening ourselves: %s\n", err)
		return
	}
	defer f.Close()

	r, err := afterarch.NewReader(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening archive: %s\n", err)
		return
	}

	for _, f := range r.File {
		fmt.Printf("--------------------\nContents of %s:\n--------------------\n", f.Name)

		rc, err := f.Open()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening %s: %s\n", f.Name, err)
			continue
		}

		_, err = io.CopyN(os.Stdout, rc, 68)
		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "Error copying %s: %s\n", f.Name, err)
			continue
		}

		rc.Close()
		fmt.Println()
	}
}

func doWrite() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s write <output file>\n", os.Args[0])
		os.Exit(1)
	}

	// Create output archive.
	arch, err := afterarch.NewWriterAfterThis(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating archive: %s\n", err)
		return
	}
	defer arch.Close()

	// Add some files to the archive.
	var files = []struct {
		Name, Body string
	}{
		{"readme.txt", "This archive contains some text files."},
		{"gopher.txt", "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
		{"todo.txt", "Get animal handling licence.\nWrite more examples."},
	}
	for _, file := range files {
		f, err := arch.Create(file.Name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating file %s: %s\n", file.Name, err)
			continue
		}

		_, err = f.Write([]byte(file.Body))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file %s: %s\n", file.Name, err)
			continue
		}
	}

	fmt.Println("Finished")
}
