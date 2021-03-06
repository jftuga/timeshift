
/*
timeshift.go
-John Taylor
Feb 2019

Shift date/time from log files or from STDIN
https://github.com/jftuga/timeshift
*/

package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "sort"
    "strings"
    "time"

    "github.com/jftuga/strtime"
    "github.com/olekukonko/tablewriter"
)

const version = "1.2.0"
const dbg = false

type timeDiff struct {
    Days int
    Hours int
    Minutes int
    Seconds int
}
var shifted *timeDiff
var aliasList map[string]string

/*
CreateAliases will build a map between an easy to remember name and a date time format specifier
*/
func CreateAliases() {
    aliasList = make(map[string]string)
    aliasList["apache_access"] = "%d/%b/%Y:%H:%M:%S"
    aliasList["apache_error"] = "%a %b %d %H:%M:%S.%f"
    aliasList["mysql_error"] = "%Y-%m-%dT%H:%M:%S.%fZ"
    aliasList["o365_exchange_trace"] = "%d/%m/%Y %I:%M:%S %p"
    aliasList["debian_log"] = "%b %d %H:%M:%S"
}

/*
ReplaceLine will substitute the date/time from the original format fo the new format

Args:
    origLine: the input line to examine

    startPos: the character position in the original line to start the replacement

    newTime: the new date/time (using the output format specifier)

Returns:
    a string containing the date/time substitution
*/
func ReplaceLine(origLine string, startPos int, oldTimeLen int, newTime string) string {
    if dbg  {
        fmt.Println("------------------------------------")
        fmt.Println("origLine  :", origLine)
        fmt.Println("startPos  :", startPos)
        fmt.Println("oldTimeLen:", oldTimeLen)
        fmt.Println("newTime   :", newTime)
    }

    //return origLine[:startPos] + newTime + origLine[startPos+len(newTime):]
    return origLine[:startPos] + newTime + origLine[startPos+oldTimeLen:]
}

/* appendYearFormat will append the current year to the beginning of a date/time
   when the year is not given in the format.  The current year is assumed

Args:
    format: a string containing the date/time format specifier

Returns:
    the same format if it already contains the year; "%Y " prepended otherwise
    empty string if format contains the year; the current year otherwise
*/
func appendYearFormat(format string) (string, string) {
	if strings.Contains(format, "%Y") || strings.Contains(format, "%y") {
		return format, ""
	}

	// assume the input refers to the current year
	year, _, _ := time.Now().Date()
	return "%Y " + format, fmt.Sprintf("%d", year)
}

/* ScanLine will convert a line from its specified input format to the specified output format
   It also will perform any time conversions

Args:
    line: one line of input

    inputFormat, outputFormat: a date/time format specifier

Returns:
    the transformed line including any shifted times
    -1 if no input format is found, the start position of the input format otherwise
*/
func ScanLine(line string, inputFormat string, outputFormat *string) (string,int) {
    if len(line) <= 2 {
        return line, -1
    }

    startPosition := 0
    var originTime time.Time

    inputFormatWithYear, year := appendYearFormat(inputFormat)

    // search for a valid timestamp with in the line and get the timestamp's startPosition
    for i,_ := range line {
		if year != "" {
			originTime, _ = strtime.Strptime(year+" "+line[i:], inputFormatWithYear)
		} else {
			originTime, _ = strtime.Strptime(line[i:], inputFormatWithYear)
		}
        if (originTime.String()[0] != 48) { // invalid time of "0001-01-01 00:00:00 +0000 UTC"
            startPosition = i
            break
        }
    }
    if (originTime.String()[0] == 48) { // failed to find a formatted time within the current line
        return line, -1
    }

    // shift the time according to the -d, -h, -m, -s cmd line options
    shiftedTime := originTime.Add( time.Hour * 24 * time.Duration(shifted.Days) + time.Hour * time.Duration(shifted.Hours) +
                    time.Minute * time.Duration(shifted.Minutes) + time.Second * time.Duration(shifted.Seconds))
    formattedShiftedTime, err := strtime.Strftime(shiftedTime, *outputFormat)
    if err != nil {
        fmt.Println("error:", err)
        os.Exit(1)
    }

    // get the length of the time in the original format
    convertedOriginTime, _ := strtime.Strftime(originTime, inputFormat)
    if dbg {
        fmt.Println("originTime          :", originTime)
        fmt.Println("convertedOriginTime :", convertedOriginTime)
        fmt.Println("startPosition       :", startPosition)
        fmt.Println("shiftedTime         :", shiftedTime)
        fmt.Println("outputFormat        :", *outputFormat)
        fmt.Println("formattedShiftedTime:", formattedShiftedTime)
    }
/*
    var j int
    if strings.HasSuffix(strings.ToUpper(formattedShiftedTime), " PM") || strings.HasSuffix(strings.ToUpper(formattedShiftedTime), " AM"){
        j = -3
    }
    currentPos := i+(len(origLine)-len(line))+j
    if currentPos < 0 {
        currentPos = 0
    }
    return ReplaceLine(origLine, currentPos, len(convertedOriginTime), formattedShiftedTime), startPosition
*/
    if line[startPosition] == ' ' {
        startPosition += 1
    }
    return ReplaceLine(line, startPosition, len(convertedOriginTime), formattedShiftedTime), startPosition
}

/*
ReadInput calls ScanLine() for each line of input and optionally displays debugging info

Args:

    input: a pointer to either an open file or to STDIN

    debugOutput: when true, display a table including start position and the line

    inputFormat, outputFormat: a date/time specifer
*/
func ReadInput(input *bufio.Scanner, debugOutput bool, inputFormat *string, outputFormat *string) {
    var newLine string
    var startPos int
    var allRows [][]string

    if(len(*outputFormat) == 0) {
        outputFormat = inputFormat
    }
    for input.Scan() {
        newLine,startPos = ScanLine(input.Text(), *inputFormat, outputFormat)
        if debugOutput {
            allRows = append(allRows, []string{fmt.Sprintf("%d",startPos),newLine})
        } else {
            fmt.Println(newLine)
        }
    }

    if debugOutput {
        table := tablewriter.NewWriter(os.Stderr)
        table.SetHeader([]string{"Start", "Input"})
        table.SetAutoWrapText(false)
        table.AppendBulk(allRows)
        table.Render()
    }
}

/*
HelpSpecifiers is called when the -F cmd line option is given
*/
func HelpSpecifiers() {
    var allRows [][]string
    allRows = append(allRows, []string{"%a", "Short week day ('mon', 'tue', etc)", ""} )
    allRows = append(allRows, []string{"%A", "Long week day ('monday', 'tuesday', etc)", ""} )
    allRows = append(allRows, []string{"%b", "Short month name ('jan', 'feb' etc)", ""} )
    allRows = append(allRows, []string{"%B", "Long month name ('january', 'february' etc)", ""} )
    allRows = append(allRows, []string{"%c", "Equivalent to `%a %b %e %H:%M:%S %Y`", ""} )
    allRows = append(allRows, []string{"%C", "Century", "Only reliable for years -9999 to 9999"} )
    allRows = append(allRows, []string{"%d", "Day of month 01-31", ""} )
    allRows = append(allRows, []string{"%D", "Equivalent to `%m/%d/%y`", ""} )
    allRows = append(allRows, []string{"%e", "Like `%d` but leading zeros are replaced by a space.", ""} )
    allRows = append(allRows, []string{"%f", "Fractional part of a second with nanosecond precision, e.g. '`123`' is 123ms; '`123456`' is 123456µs, etc.", "`Strftime` always formats using 9 digits. This can be buggy and is not recommended for use with `timeshift`."} )
    allRows = append(allRows, []string{"%1, %2, ..., %8", "Similar to `%f`.", "Use `%1` for one decimal of precision, `%2` for two decimals, etc."} )
    allRows = append(allRows, []string{"%F", "Equivalent to `%Y-%m-%d`", ""} )
    allRows = append(allRows, []string{"%h", "Equivalent to `%b`", ""} )
    allRows = append(allRows, []string{"%H", "Hours 00-23", "See also `%k`"} )
    allRows = append(allRows, []string{"%I", "Hours 01-12", "See also `%p`, `%l`"} )
    allRows = append(allRows, []string{"%j", "Day of year 000-366", ""} )
    allRows = append(allRows, []string{"%k", "Hours 0-23 (padded with spaces)", "See also `%H`"} )
    allRows = append(allRows, []string{"%l", "Hours 1-12 (padded with spaces)", "See also `%I`"} )
    allRows = append(allRows, []string{"%m", "Month 01-12", ""} )
    allRows = append(allRows, []string{"%M", "Minutes 00-59", ""} )
    allRows = append(allRows, []string{"%n", "A newline character", ""} )
    allRows = append(allRows, []string{"%p", "AM/PM", "Only valid when placed after hour-related specifiers. See also `%I`, `%l`"} )
    allRows = append(allRows, []string{"%r", "Equivalent to `%I:%M:%S %p`", ""} )
    allRows = append(allRows, []string{"%R", "Equivalent to `%H:%M`", "See also `%T`"} )
    allRows = append(allRows, []string{"%s", "Number of seconds since 1970-01-01 00:00:00 +0000 (UTC)", ""} )
    allRows = append(allRows, []string{"%S", "Seconds 00-59", ""} )
    allRows = append(allRows, []string{"%t", "A tab character", ""} )
    allRows = append(allRows, []string{"%T", "Equivalent to `%H:%M:%S`", "See also `%R`"} )
    allRows = append(allRows, []string{"%u", "The day of the week as a decimal, range 1 to 7, Monday being 1", "See also `%w`"} )
    allRows = append(allRows, []string{"%U", "The week number of the current year as a decimal number, range 00 to 53, starting with the first Sunday as the first day of week 01", "See also `%W`"} )
    allRows = append(allRows, []string{"%w", "The day of the week as a decimal, range 0 to 6, Sunday being 1", "See also `%u`"} )
    allRows = append(allRows, []string{"%W", "The week number of the current year as a decimal number, range 00 to 53, starting with the first Monday as the first day of week 01", "See also `%U`"} )
    allRows = append(allRows, []string{"%x", "Equivalent to `%D`", ""} )
    allRows = append(allRows, []string{"%X", "Equivalent to `%T`", ""} )
    allRows = append(allRows, []string{"%y", "Year without a century 00-99", "Years 00-68 are 2000-2068"} )
    allRows = append(allRows, []string{"%Y", "Year including the century", ""} )
    allRows = append(allRows, []string{"%z", "Time zone offset +/-NNNN", "`Strftime` always prints `+0000`"} )
    allRows = append(allRows, []string{"%Z", "`UTC` or `GMT`", "`Strftime` always prints `UTC`"} )

    table := tablewriter.NewWriter(os.Stderr)
    table.SetHeader([]string{"Format", "Description", "Notes"})
    table.AppendBulk(allRows)
    table.Render()
}

/*
HelpAliases is called when the -A cmd line option is given
*/
func HelpAliases() {
    var keys []string
    for k := range aliasList {
        keys = append(keys,k)
    }
    sort.Strings(keys)

    var allRows [][]string
    for _,k := range keys {
        allRows = append(allRows, []string{k, aliasList[k]})
    }

    table := tablewriter.NewWriter(os.Stderr)
    table.SetHeader([]string{"Name", "Format Spec"})
    table.AppendBulk(allRows)
    table.Render()
}

func main() {
    argsInputFormat := flag.String("i", "", "input format, see -F")
    argsOutputFormat := flag.String("o", "", "output format, see -F")
    argsHelpSpecifiers := flag.Bool("F", false, "show all formatting specifiers and then exit")
    argsHelpAliases := flag.Bool("A", false, "show all formatting aliases and then exit")
    argsInputAlias := flag.String("I", "", "input alias format, see -A")
    argsOutputAlias := flag.String("O", "", "output alias format, see -A")
    argsVersion := flag.Bool("v", false, "show program version and then exit")
    argsDebugOutput := flag.Bool("D", false, "Output the format's start position")
    msg := "use a positive number to shift forwards, negative to shift backwards in time"
    argsHours := flag.Int("h", 0, fmt.Sprintf("hours, %s", msg))
    argsMinutes := flag.Int("m", 0, fmt.Sprintf("minutes, %s", msg))
    argsSeconds := flag.Int("s", 0, fmt.Sprintf("seconds, %s", msg))
    argsDays := flag.Int("d", 0, fmt.Sprintf("days, %s", msg))
    flag.Usage = func() {
        pgmName := os.Args[0]
        if(strings.HasPrefix(os.Args[0],"./")) {
            pgmName = os.Args[0][2:]
        }
        fmt.Fprintf(os.Stderr, "\n%s: Shift date/time from log files or from STDIN.\n", pgmName)
        fmt.Fprintf(os.Stderr, "usage: %s [options] [filename|or blank for STDIN]\n\n", pgmName)
        flag.PrintDefaults()
    }
    flag.Parse()
    args := flag.Args()

    if *argsVersion {
        fmt.Fprintf(os.Stderr,"version %s\n", version)
        os.Exit(1)
    }

    if *argsHelpSpecifiers {
        HelpSpecifiers()
        os.Exit(1)
    }

    CreateAliases()
    if *argsHelpAliases {
        HelpAliases()
        os.Exit(1)
    }

    var input *bufio.Scanner
    if 0 == len(args) { // read from STDIN
        input = bufio.NewScanner(os.Stdin)
    } else { // read from filename
        fname := args[0]
        file, err := os.Open(fname)
        if err != nil {
            fmt.Fprintf(os.Stderr, "%s\n", err)
            os.Exit(1)
        }
        defer file.Close()
        input = bufio.NewScanner(file)
    }

    shifted = new(timeDiff)
    shifted.Days = *argsDays
    shifted.Hours = *argsHours
    shifted.Minutes = *argsMinutes
    shifted.Seconds = *argsSeconds

    if( len(*argsInputAlias) > 0 ) {
        if _, ok := aliasList[*argsInputAlias]; ok {
            alias := aliasList[*argsInputAlias]
            argsInputFormat = &alias
        } else {
            fmt.Fprintf(os.Stderr,"\nUnknown input alias, `%s`. Use -A to see the list of aliases\n\n", *argsInputAlias)
            os.Exit(1)
        }
    }

    if( len(*argsOutputAlias) > 0 ) {
        if _, ok := aliasList[*argsOutputAlias]; ok {
            alias := aliasList[*argsOutputAlias]
            argsOutputFormat = &alias
        } else {
            fmt.Fprintf(os.Stderr,"\nUnknown output alias, `%s`. Use -A to see the list of aliases\n\n", *argsOutputAlias)
            os.Exit(1)
        }
    }
    ReadInput(input, *argsDebugOutput, argsInputFormat, argsOutputFormat)
}

