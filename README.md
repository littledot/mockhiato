# Mockhiato

[![Build Status](https://travis-ci.org/littledot/mockhiato.svg?branch=master)](https://travis-ci.org/littledot/mockhiato)
[![Go Report Card](https://goreportcard.com/badge/github.com/littledot/mockhiato)](https://goreportcard.com/report/github.com/littledot/mockhiato)
[![Go Doc](https://godoc.org/github.com/littledot/mockhiato?status.svg)](http://godoc.org/github.com/littledot/mockhiato)

Mockhiato is a mock generation CLI tool for the Go programming language. It is designed to be fast and configurable.


- [x] Mockhiato generates mocks for entire packages recursively with a single command.
- [x] Mockhiato generates mocks for 3rd party interfaces used by your packages. (eg: If your code uses `os.FileInfo`, Mockhiato will generate mocks for it even though it is not part of your package.)
- [x] Mockhiato uses `go/loader` to speed up AST parsing.
- [x] Mockhiato is highly configurable to suit your organization's coding standards. (eg: mock name format, directory name, etc.)
- [x] Mockhiato provides easy-to-use and well-documented command line interface.
- [x] Mockhiato supports `stretchr/testify`.

![asciicast](https://cloud.githubusercontent.com/assets/14984066/25729923/7cd45c64-30ed-11e7-8e29-9145085e4748.gif)

## Installation

Install with `go get`.

```
go get -u github.com/littledot/mockhiato
```

## Usage

Mockhiato's main feature is to manage mocks for your Go project.

```
mockhiato generate
```

`generate` creates mocks for the project located at the current working directory by default. Many options are available for customization. Mockhiato currently supports `github.com/stretchr/testify/mock`. Open an issue if you would like another mocking framework to be supported.

```
mockhiato clean
```

`clean` deletes generated mocks.

```
mockhiato generate -h
```

Mockhiato is highly configurable. Append `-h` to commands for more details regarding various options and usages.
