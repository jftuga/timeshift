// Example: ./timeshift 2018-01-24

package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/knz/strtime"
)

func main() {
	argsFormat := flag.String("f", "", "use strftime format, see http://strftime.org/")
	//argsHours := flag.Int("h", 0, "use a positive number to shift forwards, negactive to shift backwards in time")
	flag.Parse()
	args := flag.Args()

	datestr := strings.Join(args, " ")
	fmt.Println()
	fmt.Printf("input: %s => %s\n", datestr, *argsFormat)

	var shifted_t time.Time
	var shifted_s string
	var err error

	if argsFormat != nil {
		shifted_t, err = strtime.Strptime(datestr, *argsFormat)
		if err != nil {
			fmt.Println("error #1:", err)
		}
		//shifted_s, err = strtime.Strftime(shifted_t, "%b %d %H:%M:%S")
		shifted_s, err = strtime.Strftime(shifted_t, "%d/%b/%Y %H:%M:%S %z")
		if err != nil {
			fmt.Println("error #2:", err)
		}
	}

	fmt.Println("shifted_t:", shifted_t)
	fmt.Println("formatted:", shifted_s)

	example := "64.242.88.10 - - [07/Mar/2004:16:47:46 -0800] \"GET /twiki/bin/rdiff/Know/ReadmeFirst?rev1=1.5&rev2=1.4 HTTP/1.1\" 200 5724"
	fmt.Println()
	fmt.Println("=================================================================")
	fmt.Println(example)

	dt_length := 26
	last := len(example) - dt_length + 1
	dtFormat := "%d/%b/%Y:%H:%M:%S %z"
	for i := 0; i < last; i++ {
		//fmt.Println(i, example[i:i+dt_length])
		shifted_t, err = strtime.Strptime(example[i:i+dt_length], dtFormat)
		if shifted_t.Year() > 1 {
			shifted_s, err = strtime.Strftime(shifted_t, "%d/%b/%Y %H:%M:%S %z")
			fmt.Println("Found:", shifted_s)
			break
		}
	}
	fmt.Println("done")
}
