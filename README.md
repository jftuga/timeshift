# timeshift
Shift date/time in log files or from STDIN

```
timeshift: Shift date/time from log files or from STDIN.
usage: timeshift [options] [filename|or blank for STDIN]

  -A	show all formatting aliases and then exit
  -D	Output the format's start position
  -F	show all formatting specifiers and then exit
  -I string
    	input alias format, see -A
  -O string
    	output alias format, see -A
  -d int
    	days, use a positive number to shift forwards, negative to shift backwards in time
  -h int
    	hours, use a positive number to shift forwards, negative to shift backwards in time
  -i string
    	input format, see -F
  -m int
    	minutes, use a positive number to shift forwards, negative to shift backwards in time
  -o string
    	output format, see -F
  -s int
    	seconds, use a positive number to shift forwards, negative to shift backwards in time
  -v	show program version and then exit
```

## Examples

```bash
# shift the time forward by 4m51s
timeshift -i "%Y-%m-%dT%H:%M:%S.%fZ" -m 5 -s -9 mysql_error.log

# same as above, but use an alias
# (use -A to see all of the aliases)
timeshift -I mysql_error -m 5 -s -9 mysql_error.log

# don't shift time, just convert to another format
timeshift -I mysql_error -o "%m/%d/%y %H:%M:%S" mysql_error.log

# convert the format and also subtract 5h
timeshift -I mysql_error -o "%m/%d/%y %H:%M:%S" -h -5 mysql_error.log
```

## Windows compliation

* Even though this is a Go project, it uses the [strtime](https://github.com/knz/strtime) package which relies on some C code.
* GCC will be needed in order to compile the [strtime](https://github.com/knz/strtime) package.
* GCC for Windows can be downloaded from here: [tdm-gcc](http://tdm-gcc.tdragon.net/)
* I used `tdm64-gcc-5.1.0-2.exe` during development.
* `goreleaser release -f .goreleaser-windows.yml --skip-publish` and then manually upload to releases

## Supported format specifiers

| Format | Description | Notes
|--------|-------------|---------
| `%a` | Short week day ("mon", "tue", etc) |
| `%A` | Long week day ("monday", "tuesday", etc) |
| `%b` | Short month name ("jan", "feb" etc) |
| `%B` | Long month name ("january", "february" etc) |
| `%c` | Equivalent to `%a %b %e %H:%M:%S %Y` |
| `%C` | Century | Only reliable for years -9999 to 9999
| `%d` | Day of month 01-31 |
| `%D` | Equivalent to `%m/%d/%y` |
| `%e` | Like `%d` but leading zeros are replaced by a space. |
| `%f` | Fractional part of a second with nanosecond precision, e.g. "`123`" is 123ms; "`123456`" is 123456µs, etc. | `Strftime` always formats using 9 digits.
| `%F` | Equivalent to `%Y-%m-%d` |
| `%h` | Equivalent to `%b` |
| `%H` | Hours 00-23  | See also `%k`
| `%I` | Hours 01-12  | See also `%p`, `%l`
| `%j` | Day of year 000-366 |
| `%k` | Hours 0-23 (padded with spaces) | See also `%H`
| `%l` | Hours 1-12 (padded with spaces) | See also `%I`
| `%m` | Month 01-12 |
| `%M` | Minutes 00-59 |
| `%n` | A newline character |
| `%p` | AM/PM | Only valid when placed after hour-related specifiers. See also `%I`, `%l`
| `%r` | Equivalent to `%I:%M:%S %p` |
| `%R` | Equivalent to `%H:%M` | See also `%T`
| `%s` | Number of seconds since 1970-01-01 00:00:00 +0000 (UTC) |
| `%S` | Seconds 00-59 |
| `%t` | A tab character |
| `%T` | Equivalent to `%H:%M:%S` | See also `%R`
| `%u` | The day of the week as a decimal, range 1 to 7, Monday being 1 | See also `%w`
| `%U` | The week number of the current year as a decimal number, range 00 to 53, starting with the first Sunday as the first day of week 01 | See also `%W`
| `%w` | The day of the week as a decimal, range 0 to 6, Sunday being 1 | See also `%u`
| `%W` | The week number of the current year as a decimal number, range 00 to 53, starting with the first Monday as the first day of week 01 | See also `%U`
| `%x` | Equivalent to `%D` |
| `%X` | Equivalent to `%T` |
| `%y` | Year without a century 00-99 | Years 00-68 are 2000-2068
| `%Y` | Year including the century |
| `%z` | Time zone offset +/-NNNN | `Strftime` always prints `+0000`
| `%Z` | `UTC` or `GMT` | `Strftime` always prints `UTC`

