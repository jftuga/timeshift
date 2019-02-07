// Example: ./timeshift 2018-01-24

package main

import (
	"flag"
	"fmt"
    "os"
    "regexp"
    "strconv"
	"strings"
	"time"

//	"github.com/knz/strtime"
)

type DateTimePosition struct {
    year uint32
    month uint32
    day uint32
    hour uint32
    minute uint32
    second uint32
    timezone uint32
}


func setup(tbl map[string]string, monthShort map[string]int, monthStrToTime map[string]time.Month, userFormat string) DateTimePosition {
    tbl["%d"] = "[0-3]*?[0-9]"
    tbl["%m"] = "[0-2]*?[0-9]"
    tbl["%Y"] = "[1-2][0-9][0-9][0-9]"
    tbl["%y"] = "[0-9][0-9]"
    tbl["%H"] = "[0-2]*?[0-9]"
    tbl["%M"] = "[0-5][0-9]"
    tbl["%S"] = "[0-5][0-9]"
    tbl["%p"] = "(?i)[AP]M"
    tbl["%z"] = "[-+]?[0-9]{4}"
    tbl["%f"] = "[0-9]{1,9}"
    tbl["%a"] = "(?i)(mon|tue|wed|thu|fri|sat|sun)"
    tbl["%A"] = "(?i)(sunday|monday|tueday|wednesday|thursday|friday|saturday)"
    tbl["%b"] = "(?i)jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec"
    tbl["%B"] = "(?i)(january|february|march|april|may|june|july|august|september|october|november|december)"
    tbl["%T"] = fmt.Sprintf("%s:%s:%s", tbl["%H"], tbl["%M"],tbl["%S"])
    tbl["%D"] = fmt.Sprintf("%s/%s/%s", tbl["%m"], tbl["%d"],tbl["%y"])
    tbl["%F"] = fmt.Sprintf("%s-%s-%s", tbl["%Y"], tbl["%m"],tbl["%d"])

    monthShort["jan"] = 1
    monthShort["feb"] = 2
    monthShort["mar"] = 3
    monthShort["apr"] = 4
    monthShort["may"] = 5
    monthShort["jun"] = 6
    monthShort["jul"] = 7
    monthShort["aug"] = 8
    monthShort["sep"] = 9
    monthShort["oct"] = 10
    monthShort["nov"] = 11
    monthShort["dec"] = 12

    monthStrToTime["Jan"] = time.January
    monthStrToTime["Feb"] = time.February
    monthStrToTime["Mar"] = time.March
    monthStrToTime["Apr"] = time.April
    monthStrToTime["May"] = time.May
    monthStrToTime["Jun"] = time.June
    monthStrToTime["Jul"] = time.July
    monthStrToTime["Aug"] = time.August
    monthStrToTime["Sep"] = time.September
    monthStrToTime["Oct"] = time.October
    monthStrToTime["Nov"] = time.November
    monthStrToTime["Dec"] = time.December
    monthStrToTime["January"] = time.January
    monthStrToTime["February"] = time.February
    monthStrToTime["March"] = time.March
    monthStrToTime["April"] = time.April
    monthStrToTime["May"] = time.May
    monthStrToTime["June"] = time.June
    monthStrToTime["Julu"] = time.July
    monthStrToTime["August"] = time.August
    monthStrToTime["September"] = time.September
    monthStrToTime["October"] = time.October
    monthStrToTime["November"] = time.November
    monthStrToTime["December"] = time.December


    var year uint32
    var month uint32
    var day uint32
    var hour uint32
    var minute uint32
    var second uint32
    var timezone uint32

    //fmt.Println("userFormat2:", userFormat)

    formatRE := regexp.MustCompile("(%.)")
    result := formatRE.FindAllStringSubmatch(userFormat, -1)

    fmt.Println("result2:", result)

    for i, val := range result {
        fmt.Printf("%d %s\n", i, val[0])
        switch val[0] {
            case "%d":
                day = uint32(i+1)
            case "%m", "%b":
                month = uint32(i+1)
            case "%Y", "%y":
                year = uint32(i+1)
            case "%H":
                hour = uint32(i+1)
            case "%M":
                minute = uint32(i+1)
            case "%S":
                second = uint32(i+1)
            case "%z":
                timezone = uint32(i+1)
            default:
                fmt.Println("err #2: unknown format string:", val[0])
                os.Exit(1)
        }
    }
    fmt.Println("x2:", year, month, day, hour, minute, second)
    fmt.Println("------------------------------------")
    fmt.Println()
    return DateTimePosition{year, month, day, hour, minute, second, timezone}
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

    var monthShort map[string]int
    monthShort = make(map[string]int)

    var monthStrToTime map[string]time.Month
    monthStrToTime = make(map[string]time.Month)

	////argsFormat := flag.String("f", "", "use strftime format, see http://strftime.org/")
	//argsHours := flag.Int("h", 0, "use a positive number to shift forwards, negative to shift backwards in time")
	flag.Parse()
	////args := flag.Args()
/*
	datestr := strings.Join(args, " ")
    dtPos := setup(formats,*argsFormat)
	fmt.Println()
	fmt.Printf("input: %s => %s\n", datestr, *argsFormat)
    datetimeRE := transform(formats, *argsFormat)
    fmt.Println("datetimeRE:", datetimeRE)

    result := datetimeRE.FindStringSubmatch(datestr)
    fmt.Printf("result [%s] %s\n", datestr, result)
    if len(result) == 0 {
        fmt.Fprintf(os.Stderr, "Failed to converge\n");
        os.Exit(1)
    }
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

    if year <= 68 {
        year += 2000
    } else {
        year += 1900
    }
    nativeDate := time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local)
    fmt.Println("nativeDate:", nativeDate)
*/
    fmt.Printf("\n--------------------------------------------\n\n")
	example := "64.242.88.10 - - [07/Mar/2004:16:47:46 -0800] \"GET /twiki/bin/rdiff/Know/ReadmeFirst?rev1=1.5&rev2=1.4 HTTP/1.1\" 200 5724"
    logFormat := "%d/%b/%Y:%H:%M:%S %z"
    fmt.Println("logFormat: ", logFormat)

    dtPos := setup(formats,monthShort,monthStrToTime,logFormat)
    fmt.Println("dtPos:", dtPos)

    fmt.Printf("\n--------------------------------------------\n\n")
    datetimeRE := transform(formats,logFormat)
    fmt.Println("datetimeRE:", datetimeRE)

    result := datetimeRE.FindStringSubmatch(example)
    fmt.Printf("result3:%s\n\n", result)

    if len(result) == 0 {
        fmt.Fprintf(os.Stderr, "Failed to converge\n");
        os.Exit(1)
    }
    for i,r := range result {
        fmt.Printf("\t[%d] %s\n", i,r)
    }
    fmt.Printf("\n--------------------------------------------\n\n")


    year, _ := strconv.Atoi(result[dtPos.year])
    fmt.Println("yyy:", year)
    month, err := strconv.Atoi(result[dtPos.month])
    useNativeMonth := false
    var monthTime time.Month
    if err != nil {
        fmt.Println("err #1 detected!")
        // do something for named Months, such as Mar, July, etc
        m := strings.ToLower(result[dtPos.month])
        fmt.Println("m,res:",m,result[dtPos.month])
        month, _ = monthShort[m]  //FIXME
        useNativeMonth = true
        monthTime = monthStrToTime[result[dtPos.month]]
    }
    fmt.Println("month    :",month)
    fmt.Println("monthTime:",monthTime)
    day, _ := strconv.Atoi(result[dtPos.day])
    hour, _ := strconv.Atoi(result[dtPos.hour])
    minute, _ := strconv.Atoi(result[dtPos.minute])
    second, _ := strconv.Atoi(result[dtPos.second])

    if year <= 68 {
        year += 2000
    } else {
        if year < 100 {
            year += 1900
        }
    }
    var nativeDate time.Time
    loc := time.FixedZone(" ", -8*3600)
    if useNativeMonth {
        nativeDate = time.Date(year, monthTime, day, hour, minute, second, 0, loc)
    } else {
        nativeDate = time.Date(year, time.Month(month), day, hour, minute, second, 0, loc)
    }
    fmt.Println("nativeDate:", nativeDate)


    return
    /*
    fmt.Printf("\n--------------------------------------------\n\n")

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
*/
}
