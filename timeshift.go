// Example: ./timeshift 2018-01-24

package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "strings"
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
    apache_access = "%d/%b/%Y:%H:%M:%S"
    apache_error = "%a %b %d %H:%M:%S.%f"
    mysql_error = "%Y-%m-%dT%H:%M:%S.%fZ"
    o365_exchange_trace = ""%d/%m/%Y %-I:%M:%S %p"
*/

func replaceLine(origLine string, startPos int, newTime string) string {
    return origLine[:startPos] + newTime + origLine[startPos+len(newTime):]
}

func scanLine(line string, format string, shifted timeDiff) string {
    //fmt.Println(line)
    var originTime time.Time
    var i int
    if len(line) <= 2 {
        return line
    }
    origLine := line
    if(startPosition > 0) {
        line = line[startPosition:]
    }
    //fmt.Println("lineTrunc:", line)
    for i,_ = range line {
        originTime,_ = strtime.Strptime(line[i:], format)
        //fmt.Println("originTime:", originTime)
        if (originTime.String()[0] != 48) { // invalid time of "0001-01-01 00:00:00 +0000 UTC"
            startPosition = i
            break
        }
    }
    if (originTime.String()[0] == 48) { // failed to find a formatted time within the current line
        return origLine
    }
    //fmt.Println("ot:", originTime)
    shiftedTime := originTime.Add( time.Hour * 24 * time.Duration(shifted.Days) + time.Hour * time.Duration(shifted.Hours) +
                    time.Minute * time.Duration(shifted.Minutes) + time.Second * time.Duration(shifted.Seconds))
    formattedShiftedTime, _ := strtime.Strftime(shiftedTime, format)
    //fmt.Println("fst:", formattedShiftedTime)

    var j int
    if strings.HasSuffix(strings.ToUpper(formattedShiftedTime), " PM") || strings.HasSuffix(strings.ToUpper(formattedShiftedTime), " AM"){
        j = -3
    }
    return replaceLine(origLine, i+(len(origLine)-len(line))+j, formattedShiftedTime)
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
