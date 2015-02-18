package main

import (
	"fmt"
	"os/exec"
)

func main() {

	out, _ := exec.Command("afplay", "test.m4a", "-r", "1").CombinedOutput()
	fmt.Print(string(out))

}
