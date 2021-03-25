# o5≠m4

[![Go](https://github.com/qwasa-net/o5/actions/workflows/go.yml/badge.svg)](https://github.com/qwasa-net/o5/actions/workflows/go.yml)

```
o5 -- super simple micro macro processor for text files
  -d value
        define macro variable (-d NAME=VALUE)
  -end string
        macro closer (suffix) (default "-->")
  -i string
        input file ('-' is stdin) (default "-")
  -o string
        output file ('-' is stdout) (default "-")
  -start string
        macro openner (prefix) (default "<!--#")
  -trim
        trim spaces in expanded macro (default true)
  -w string
        working directory (for file includes) (default ".")
```

## Rules

**\<!--# NAME --\>** — replace with variable NAME (defined with -d NAME=value)

**\<!--# @/path/file --\>** — replace with file raw content

**\<!--# $USER --\>** — replace with system environment variable

Macro prefix (*\<!--*) and suffix (*--\>*) can be changed with **-start** and **--end** options,
e.g.: `-start '###' -end '###'`.


## Example

### input.txt

```
Hello, ### $USER ###!
Today is ### today ###.
This is your /etc/passwd:
### @/etc/passwd ###
```

### o5 call

```bash
./o5 -start '###' -end '###' -d today=Now -i input.txt -w '/' -o -
```

### stdout

```
Hello, login!
Today is Now.
This is your /etc/passwd:
root:x:0:0:root:/root:/bin/bash
daemon:x:1:1:daemon:/usr/sbin:/usr/sbin/nologin
…
```