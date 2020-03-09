package main

import (
	"fmt"

	"github.com/fairjungle/frontapp/cmd/frontapp/cmds"
)

func main() {
	if err := cmds.Execute(); err != nil {
		fmt.Println(err)
	}
}
