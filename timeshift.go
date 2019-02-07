// Example: ./timeshift 2018-01-24

package main

import (
    "bufio"
	"flag"
	"fmt"
    "os"
    "time"

	"github.com/knz/strtime"
)

var version string
var startPosition int

type timeDiff struct {
    Days int
    Hours int
    Minutes int
    Seconds int
}


/*
    originTime, _ := strtime.Strptime(args[0], *argsFormat)
    shiftedTime := originTime.Add( time.Hour * 24 * time.Duration(*argsDays) +
        time.Hour * time.Duration(*argsHours) + time.Minute * time.Duration(*argsMinutes) + time.Second * time.Duration(*argsSeconds))
    formattedShiftedTime, _ := strtime.Strftime(shiftedTime, *argsFormat)
    fmt.Println(shiftedTime)
    fmt.Println(formattedShiftedTime)

*/

func replaceLine(origLine string, startPos int, newTime string) string {
    return origLine[:startPos] + newTime + origLine[startPos+len(newTime):]
}

func scanLine(line string, format string, shifted timeDiff) string {
    //fmt.Println(line)
    var originTime time.Time
    var i int
    origLine := line
    if(startPosition > 0) {
        line = line[startPosition:]
    }
    //fmt.Println("line:", line)
    for i,_ = range line {
        originTime,_ = strtime.Strptime(line[i:], format)
        if (originTime.String()[0] != 48) { // invalid time of "0001-01-01 00:00:00 +0000 UTC"
            startPosition = i
            break
        }
    }
    //fmt.Println("ot:", originTime)
    shiftedTime := originTime.Add( time.Hour * 24 * time.Duration(shifted.Days) + time.Hour * time.Duration(shifted.Hours) +
                    time.Minute * time.Duration(shifted.Minutes) + time.Second * time.Duration(shifted.Seconds))
    formattedShiftedTime, _ := strtime.Strftime(shiftedTime, format)
    //fmt.Println("fst:", formattedShiftedTime)

    return replaceLine(origLine, i+(len(origLine)-len(line)), formattedShiftedTime)
}

func ReadInput(input *bufio.Scanner, format string, shifted timeDiff) {
    var newLine string
    for input.Scan() {
        newLine = scanLine(input.Text(), format, shifted)
        fmt.Println(newLine)
    }
}


func main() {
	argsFormat := flag.String("f", "", "use strftime format, see http://strftime.org/")
    msg := "use a positive number to shift forwards, negative to shift backwards in time"
	argsHours := flag.Int("h", 0, fmt.Sprintf("hours, %s", msg))
	argsMinutes := flag.Int("m", 0, fmt.Sprintf("minutes, %s", msg))
	argsSeconds := flag.Int("s", 0, fmt.Sprintf("seconds, %s", msg))
    argsDays := flag.Int("d", 0, fmt.Sprintf("days, %s", msg))
    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "\n%s %s, Shift date/time from log files or from STDIN.\n\n", os.Args[0], version)
        fmt.Fprintf(os.Stderr, "Usage for %s:\n", os.Args[0])
        flag.PrintDefaults()
    }
	flag.Parse()
	args := flag.Args()

    var input *bufio.Scanner
    if 0 == len(args) { // read from STDIN
        input = bufio.NewScanner(os.Stdin)
    } else { // read from filename
        fname := args[0]
        file, err := os.Open(fname)
        if err != nil {
            fmt.Fprintf(os.Stderr, "%s\n", err)
            return
        }
        defer file.Close()
        input = bufio.NewScanner(file)
    }

    diff := timeDiff{*argsDays, *argsHours, *argsMinutes, *argsSeconds}
    ReadInput(input, *argsFormat, diff )
}
