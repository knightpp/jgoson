# jgoson

This is a simple app that generates Golang code from JSON. It can generate inlined structs as well as
separate structs.

## Installation

```bash
go install github.com/knightpp/jgoson/cmd/jgoson@latest
```

or with Nix flakes

```bash
nix profile install github:knightpp/jgoson
```

## How to use

```bash
cat example.json | jgoson -inline=true -tag=yaml -tag-opts="a,b,c"
```

## Known issues

- [ ] Does not generate unique struct names for nested structs
