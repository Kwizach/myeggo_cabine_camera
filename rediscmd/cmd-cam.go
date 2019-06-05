package rediscmd

import (
	"fmt"
)

func init() {
	allCommands["SHOOT"] = func(args []string) error { return shoot(args) }
}

// shoot photo with camera and store it localy first
// args:
//		[0] when
//		[1] how (LR/HR)
func shoot(args []string) error {
	fmt.Println(args[0])
	return nil
}
