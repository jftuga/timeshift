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

/*
func setup(tbl map[string]*regexp.Regexp) {
    tbl["%d"] = regexp.MustCompile("[0-3][0-9]")
    tbl["%m"] = regexp.MustCompile("[0-2][0-9]")
    tbl["%Y"] = regexp.MustCompile("[1-2][0-9][0-9][0-9]")
    tbl["%H"] = regexp.MustCompile("[0-2][0-9]")
    tbl["%M"] = regexp.MustCompile("[0-5][0-9]")
    tbl["%S"] = regexp.MustCompile("[0-5][0-9]")
    tbl["%p"] = regexp.MustCompile("(?i)[AP]M")

}
*/

func setup(tbl map[string]string) {
    tbl["%d"] = "[0-3]*?[0-9]"
    tbl["%m"] = "[0-2]*?[0-9]"
    tbl["%Y"] = "[1-2][0-9][0-9][0-9]"
    tbl["%H"] = "[0-2]*?[0-9]"
    tbl["%M"] = "[0-5][0-9]"
    tbl["%S"] = "[0-5][0-9]"
    tbl["%p"] = "(?i)(A|P)M"
}

func transform(formats map[string]string, userFormat string) *regexp.Regexp {
    for key, val := range formats {
        fmt.Printf("%T %s %T %s\n", key, key, val, val)
        userFormat = strings.Replace(userFormat, key, fmt.Sprintf("(%s)", val), -1)
    }

    return regexp.MustCompile(userFormat)
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
    datetimeRE := transform(formats, *argsFormat)
    fmt.Printf("[%s]: %s\n", datestr, datetimeRE.FindString(datestr))
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
