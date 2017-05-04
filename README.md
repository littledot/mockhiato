# Mockhiato

[![Build Status](https://travis-ci.org/littledot/mockhiato.svg?branch=master)](https://travis-ci.org/littledot/mockhiato)
[![Go Report Card](https://goreportcard.com/badge/github.com/littledot/mockhiato)](https://goreportcard.com/report/github.com/littledot/mockhiato)

Mockhiato is a mock generation CLI tool for the Go programming language. It is designed to be fast, configurable and correct.

[![asciicast](https://asciinema.org/a/5040vzaydr0g9dj1gl7bqrdnh.png)](https://asciinema.org/a/5040vzaydr0g9dj1gl7bqrdnh)

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

`generate` creates mocks for the project located at the current working directory. Mockhiato currently supports `github.com/stretchr/testify/mock`. Open an issue if you would like another mocking framework to be supported.

```
mockhiato clean
```

`clean` deletes generated mocks.

```
mockhiato generate -h
```

Append `-h` to commands for more details on other usages.
