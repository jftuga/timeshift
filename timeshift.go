// Example: ./timeshift 2018-01-24

package main

import (
	"flag"
	"fmt"
    "regexp"
    "strconv"
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


func setup(tbl map[string]string, userFormat string) DateTimePosition {
    tbl["%d"] = "[0-3]*?[0-9]"
    tbl["%m"] = "[0-2]*?[0-9]"
    tbl["%Y"] = "[1-2][0-9][0-9][0-9]"
    tbl["%H"] = "[0-2]*?[0-9]"
    tbl["%M"] = "[0-5][0-9]"
    tbl["%S"] = "[0-5][0-9]"
    tbl["%p"] = "(?i)[AP]M"

    var year uint32
    var month uint32
    var day uint32
    var hour uint32
    var minute uint32
    var second uint32

    //fmt.Println("userFormat2:", userFormat)

    formatRE := regexp.MustCompile("(%.)")
    result := formatRE.FindAllStringSubmatch(userFormat, -1)
    //fmt.Println("result2:", result)

    for i, val := range result {
        //fmt.Printf("%d %s\n", i, val[0])
        switch val[0] {
            case "%d":
                day = uint32(i+1)
            case "%m":
                month = uint32(i+1)
            case "%Y":
                year = uint32(i+1)
            case "%H":
                hour = uint32(i+1)
            case "%M":
                minute = uint32(i+1)
            case "%S":
                second = uint32(i+1)
        }
    }
    //fmt.Println("x2:", year, month, day, hour, minute, second)
    return DateTimePosition{year, month, day, hour, minute, second}
}

func transform(formats map[string]string, userFormat string) *regexp.Regexp  {
    for key, val := range formats {
        userFormat = strings.Replace(userFormat, key, fmt.Sprintf("(%s)", val), -1)
    }

    return regexp.MustCompile(userFormat)
}

func main() {
    //var formats map[string]*regexp.Regexp
    //formats = make(map[string]*regexp.Regexp)
    var formats map[string]string
    formats = make(map[string]string)

	argsFormat := flag.String("f", "", "use strftime format, see http://strftime.org/")
	//argsHours := flag.Int("h", 0, "use a positive number to shift forwards, negative to shift backwards in time")
	flag.Parse()
	args := flag.Args()

	datestr := strings.Join(args, " ")
    dtPos := setup(formats,*argsFormat)
	fmt.Println()
	fmt.Printf("input: %s => %s\n", datestr, *argsFormat)
    datetimeRE := transform(formats, *argsFormat)

    result := datetimeRE.FindStringSubmatch(datestr)
    fmt.Printf("result [%s]\n", datestr)
    for i,r := range result {
        fmt.Printf("\t[%d] %s\n", i,r)
    }
    fmt.Println("dtPos:", dtPos)
    year, _ := strconv.Atoi(result[dtPos.year])
    month, _ := strconv.Atoi(result[dtPos.month])
    day, _ := strconv.Atoi(result[dtPos.day])
    hour, _ := strconv.Atoi(result[dtPos.hour])
    minute, _ := strconv.Atoi(result[dtPos.minute])
    second, _ := strconv.Atoi(result[dtPos.second])

    nativeDate := time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local)
    fmt.Println("nativeDate:", nativeDate)
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
