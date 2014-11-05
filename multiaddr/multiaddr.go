package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	ma "github.com/jbenet/go-multiaddr"
)

var usage = `multiaddr conversion
usage: multiaddr [fmt] <>`

var formats = []string{"string", "bytes", "hex", "raw"}
var format string

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s <multiaddr>\n\nFlags:\n", os.Args[0])
		flag.PrintDefaults()
	}

	usage := fmt.Sprintf("output format, one of: %v", formats)
	flag.StringVar(&format, "format", "string", usage)
	flag.StringVar(&format, "f", "string", usage+" (shorthand)")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		die("error: can only take one argument")
	}

	m, err := ma.NewMultiaddr(args[0])
	if err != nil {
		die(err)
	}

	fmt.Println(outfmt(m))
}

func outfmt(m ma.Multiaddr) string {
	switch format {
	case "string":
		return m.String()
	case "bytes":
		return fmt.Sprintf("%v", m.Bytes())
	case "raw":
		return string(m.Bytes())
	case "hex":
		return "0x" + hex.EncodeToString(m.Bytes())
	}

	die("error: invalid format", format)
	return ""
}

func die(v ...interface{}) {
	fmt.Fprint(os.Stderr, v...)
	fmt.Fprint(os.Stderr, "\n")
	flag.Usage()
	os.Exit(-1)
}
