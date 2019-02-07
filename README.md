# timeshift
Shift date/time in log files or from STDIN

```
Usage for timeshift:
  -d int
        days, use a positive number to shift forwards, negative to shift backwards in time
  -f string
        use strftime format, see http://strftime.org/
  -h int
        hours, use a positive number to shift forwards, negative to shift backwards in time
  -m int
        minutes, use a positive number to shift forwards, negative to shift backwards in time
  -s int
        seconds, use a positive number to shift forwards, negative to shift backwards in time
```

## Examples

```
(to do)
```

## Windows compliation

* Even though this is a Go project, it uses the [strtime](https://github.com/knz/strtime) package which relies on some C code.
* GCC will be needed in order to compile the [strtime](https://github.com/knz/strtime) package.
* GCC for Windows can be downloaded from here: [tdm-gcc](http://tdm-gcc.tdragon.net/)
* I used `tdm64-gcc-5.1.0-2.exe` during development.
