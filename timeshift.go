
// Example: ./timeshift 2018-01-24

package main

import (
    "fmt"
    "os"
    "time"

    "github.com/araddon/dateparse"
)

func main() {
    datestr := os.Args[1]
    fmt.Println("input:", datestr)
    t, err := dateparse.ParseLocal(datestr)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal Error: %s\n", err)
        return
    }

    fmt.Println( t.In(time.UTC).String() )
}
