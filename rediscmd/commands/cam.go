package commands

import (
	"fmt"
)

func init() {
	AllCommands["SHOOT"] = func(args []string) error { return shoot(args) }
}

// shoot photo with camera and store it localy first
// args:
//		[0] when
//		[1] how (LR/HR)
func shoot(args []string) error {
	fmt.Println(args[0])
	return nil
}
