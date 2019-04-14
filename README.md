# bl Brightness level controll command

Linux bridgtness controll command for my laptop


**Install**

```
$ go get github.com/tomlla/bl
```

**Usage**

```
// Increase backlight
$ sudo bl inc 1

// Decrease backlght
$ sudo bl dec 1
```

**(Optional Setting) Add `bl` command to you sudoers**


Run `$ sudo visudo` and add a line like this.

```sudo

nt ALL=(ALL) NOPASSWD: {{ bl command absolute path }}
```
