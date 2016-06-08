package main

import (
    "fmt"
    "os/exec"
)

var code = '''
package main
import "fmt"

func main() {
  fmt.Println("hello, world")
}

'''

// add types
// add vars
// add functions

// TODO: exec test. execute some code, and then change source code, and then re-build code,
// finally replace executable file with new one.

func main() {
    fmt.Println("hello, world")
    execAfter()
    fmt.Println("program exited.")
}

func execAfter() {
    exec.Command("sleep 5; echo 'xxx' | tee /tmp/temptxt").Run()
}
