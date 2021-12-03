# profaneword
profane password generator (probably insecure), as suggested by [u/gatestone](https://www.reddit.com/r/golang/comments/r5hn12/comment/hmnyk9k/?utm_source=share&utm_medium=web2x&context=3). This is still missing some requirements: special characters etc. maybe a `--safe` flag that inserts special characters or something. but it's a start!

## how it works

install:

```bash
go install github.com/MikkelHJuul/profaneword/profaneword@v0.1.1
```
play
```bash 
❯ profaneword 
misbehaved stoner

~ 
❯ profaneword --extend
shagging hell terminator

~ 
❯ profaneword --EXTEND
bad-breathed bellend rapping muff-diving molester

~ 
❯ profaneword 1337
cr1k3y 51r W4nk-4-107

~ 
❯ profaneword /s
milKINg SUCKEr

~ 
❯ profaneword --EXTEND /s
satAnic oFfeNDeR gOLLY rUBbiSH doNKeY

❯ profaneword --help
profaneword is a CLI library and tool for generating obscene/profane passwords. 
It's probably not particularly safe to use, as these passwords will be easy to brute force; 
if an attacker knows you use this generator. But hey, it's just for fun.

Usage:
  profaneword [flags] [..args]
  profaneword [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  help        Help about any command
  obscure     apply formatters on std in
  version     print the version and exit

Flags:
      --EXTEND                lengthen the output further (extensiveness+3)
      --extend                lengthen the output (extensiveness+1)
  -e, --extensiveness int16   how long (number of words) the password should be (default 2)
  -h, --help                  help for profaneword

Args:
	1337  output formatted as 1337-speak
	/s    sARcaSTiC OUtpUt


Use "profaneword [command] --help" for more information about a command.

~ 


❯ echo DASDAD |profaneword obscure /s 
DasDad

~ 
❯ echo DASDAD |profaneword obscure 1337
D45D4D

```
