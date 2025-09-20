package main

import (
	"github.com/abcfe-op/abcfe-node/app"
)

func main() {
	if n, err := app.New(); err != nil {
		panic(err)
	} else {
		n.Wait()
	}
}
