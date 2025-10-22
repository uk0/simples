package main

import (
	"fmt"

	"github.com/otiai10/opengraph/v2"
)

func main() {
	ogp, err := opengraph.Fetch("https://firsh.me/blog/0145.html")
	fmt.Println(ogp, err)
}
