package main

import (
	"flag"
	"fmt"
)

var (
	Suffix = flag.String("suffix", "", "Unit suffix to restrict polling to (i.e. 63 for STA63 units)")
)

func main() {
	flag.Parse()

	fmt.Println("Test")
}
