# timeshift
Shift date/time in log files or from STDIN

```
timeshift: Shift date/time from log files or from STDIN.
usage: timeshift [options] [filename|or blank for STDIN]

  -A	show all formatting aliases and then exit
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
