# o5≠m4


```
o5 -- super simple micro macro processor (for text files only!)
  -d value
        define macro value
  -end string
        macro closer (suffix) (default "-->")
  -i string
        input file (stdin) (default "-")
  -o string
        output file (stdin) (default "-")
  -start string
        macro openner (prefix) (default "<!--#")
  -trim
        trim spaces in expanded macro (default true)
```

## input.txt

```
Hello, ### $USER ###!
Today is ### today ###.
This is your /etc/passwd:
### @/etc/passwd ###

```

## call

```bash
./o5 -start '###' -end '###' -d today=Now -i input.txt -o -
```

## stdout

```
Hello, login!
Today is Now.
This is your /etc/passwd:
root:x:0:0:root:/root:/bin/bash
daemon:x:1:1:daemon:/usr/sbin:/usr/sbin/nologin
…
```