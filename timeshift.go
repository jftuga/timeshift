// Example: ./timeshift 2018-01-24

package main

import (
	"flag"
	"fmt"
    "regexp"
	"strings"
	"time"

	"github.com/knz/strtime"
)

type DateTimePosition struct {
    year uint32
    month uint32
    day uint32
    hour uint32
    minute uint32
    second uint32
}


func setup(tbl map[string]string) {
    tbl["%d"] = "[0-3]*?[0-9]"
    tbl["%m"] = "[0-2]*?[0-9]"
    tbl["%Y"] = "[1-2][0-9][0-9][0-9]"
    tbl["%H"] = "[0-2]*?[0-9]"
    tbl["%M"] = "[0-5][0-9]"
    tbl["%S"] = "[0-5][0-9]"
    tbl["%p"] = "(?i)[AP]M"
}

func transform(formats map[string]string, userFormat string) (*regexp.Regexp, DateTimePosition) {
    var year uint32
    var month uint32
    var day uint32
    var hour uint32
    var minute uint32
    var second uint32
    var previous string
    var i uint32
    for key, val := range formats {
        previous = userFormat
        userFormat = strings.Replace(userFormat, key, fmt.Sprintf("(%s)", val), -1)
        if previous == userFormat {
            continue
        }
        i += 1

        switch key {
        case "%d":
            day = i-1
        case "%m":
            month = i-1
        case "%Y":
            year = i-1
        case "%H":
            hour = i-1
        case "%M":
            minute = i-1
        case "%S":
            second = i-1
        }
    }

    fmt.Println("position:", year, month, day, hour, minute, second)
    return regexp.MustCompile(userFormat), DateTimePosition{year, month, day, hour, minute, second}
}

func main() {
    //var formats map[string]*regexp.Regexp
    //formats = make(map[string]*regexp.Regexp)
    var formats map[string]string
    formats = make(map[string]string)
    setup(formats)

	argsFormat := flag.String("f", "", "use strftime format, see http://strftime.org/")
	//argsHours := flag.Int("h", 0, "use a positive number to shift forwards, negative to shift backwards in time")
	flag.Parse()
	args := flag.Args()

	datestr := strings.Join(args, " ")
	fmt.Println()
	fmt.Printf("input: %s => %s\n", datestr, *argsFormat)
    datetimeRE, dtPosition := transform(formats, *argsFormat)

    result := datetimeRE.FindStringSubmatch(datestr)
    fmt.Printf("result [%s]: %v\n", datestr, dtPosition)
    for i,r := range result {
        fmt.Printf("\t[%d] %s\n", i,r)
    }
    return

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
