package main

import (
	"fmt"
	"log"
	"os"
	"rast/packdir"

	"github.com/akamensky/argparse"
)

func main() {
	// parse the command line
	parser := argparse.NewParser("snapdir", "Create a compressed snapshot of a directory.")
	in := parser.String("i", "indir",
		&argparse.Options{Required: true,
			Help: "Set the input directory from which to create the snapshot"})
	out := parser.String("o", "outfile",
		&argparse.Options{Required: true,
			Help: "Set the output file containing the compressed snapshot"})
	level := parser.Int("l", "level",
		&argparse.Options{Required: false,
			Help:    "Set the compression level: -2 only Huffman, -1 default, 0 no compression, 1 (fastest) to 9 (highest compression)",
			Default: 2})
	targetBaseDir := parser.String("d", "basedir",
		&argparse.Options{Required: false,
			Help: "Set the snapshot base directory, all files of the snapshot will be in basedir/"})
	noErrors := parser.Flag("E", "no-errors",
		&argparse.Options{Required: false,
			Help: "Do not display error messages"})
	noInfo := parser.Flag("M", "no-messages",
		&argparse.Options{Required: false,
			Help: "Do not display info messages"})
	quiet := parser.Flag("q", "quiet",
		&argparse.Options{Required: false,
			Help: "Do not display any messages (implies -E -M)"})
	verbose := parser.Flag("v", "verbose",
		&argparse.Options{Required: false,
			Help: "Display additional information"})
	pBar := parser.Flag("p", "progressbar",
		&argparse.Options{Required: false,
			Help: "Display a progress bar"})

	err := parser.Parse(os.Args)
	_ = out
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}
	var flags int
	if *verbose {
		flags = packdir.PRINT_INFO | packdir.PRINT_ERRORS | packdir.VERBOSE
	} else {
		flags = packdir.PRINT_INFO | packdir.PRINT_ERRORS | packdir.PROGRESSBAR
	}
	if *quiet {
		*noErrors = false
		*noInfo = false
		*pBar = false
		*verbose = false
	}
	if *noErrors {
		flags &^= packdir.PRINT_ERRORS
	}
	if *noInfo {
		flags &^= packdir.PRINT_INFO
	}
	if !*pBar {
		flags &^= packdir.PROGRESSBAR
	}

	// set up logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.SetPrefix("ERROR ")
	// do the packing
	_, err = packdir.Pack(*in, *out, *targetBaseDir,
		packdir.CompressionLevel(*level), flags)
	if err != nil {
		os.Exit(2)
	}
}
