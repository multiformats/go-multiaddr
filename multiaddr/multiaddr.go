package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	ma "github.com/jbenet/go-multiaddr"
	manet "github.com/jbenet/go-multiaddr/net"
)

var formats = []string{"string", "bytes", "hex", "slice"}
var format string

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [<multiaddr>]\n\nFlags:\n", os.Args[0])
		flag.PrintDefaults()
	}

	usage := fmt.Sprintf("output format, one of: %v", formats)
	flag.StringVar(&format, "format", "string", usage)
	flag.StringVar(&format, "f", "string", usage+" (shorthand)")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		maddrs, err := manet.InterfaceMultiaddrs()
		if err != nil {
			die(err)
		}

		output(maddrs...)
		return
	}

	m, err := ma.NewMultiaddr(args[0])
	if err != nil {
		die(err)
	}

	output(m)
}

func output(ms ...ma.Multiaddr) {
	for _, m := range ms {
		fmt.Println(outfmt(m))
	}
}

func outfmt(m ma.Multiaddr) string {
	switch format {
	case "string":
		return m.String()
	case "slice":
		return fmt.Sprintf("%v", m.Bytes())
	case "bytes":
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
