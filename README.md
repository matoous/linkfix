# Linkfix [![Godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/matoous/linkfix) [![license](https://img.shields.io/github/license/matoous/linkfix)](https://raw.githubusercontent.com/matoous/linkfix/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/matoous/linkfix)](http:/goreportcard.com/report/matoous/linkfix) ![Build](https://github.com/matoous/linkfix/workflows/Build/badge.svg)

![Logo](/.github/logo.png)

Linkfix is a simple tool that helps you avoid the [Link rot](https://en.wikipedia.org/wiki/Link_rot)
by reporting on no-longer working links in your files and suggesting replacements with
[Wayback Machine snapshots](https://archive.org/web/) wherever possible.

Linkfix can be also used as link availability checker with support for various protocols
and complex error reporting.

## Link rot

> Link rot is the phenomenon of hyperlinks tending over time to cease to point
> to their originally targeted file, web page, or server due to that resource
> being relocated or becoming permanently unavailable. A link that no longer points to its target,
> often called a broken or dead link, is a specific form of dangling pointer. 
> The rate of link rot is a subject of study and research due to its significance 
> to the internet's ability to preserve information. Estimates of that rate vary dramatically
> between studies.

## Features

1. Automatic replacement suggestions for http links by snapshots from [Wayback Machine](https://archive.org/web/).
2. Multithreaded with configurable level of concurrency.
3. Multiple output formats.
4. Supports various protocols `http`, `https`, `ftp`, `file`.

## Installation

```shell script
go get -u github.com/matoous/linkfix
```

or run using the docker image without installation ![](https://img.shields.io/docker/image-size/matousdz/linkfix?sort=date)

```shell script
docker run --rm -it -u $(id -u):$(id -g) -v $(pwd):/app linkfix -w 4 /app
```

## Usage

```shell script
linkfix .
```

Advanced usage:

```shell script
linkfix . --ignore "internal/old_files/**" --verbose --exclude "dont.check.me" --workerks 16 --yes
```

This will run the `linkfix` from current directory, ignoring all files under `internal/old_files/`,
in a verbose mode, excluding urls on `http://dont.check.me` (e.g. won't try to fix URL that are hosted
on `http://dont.check.me`), with 16 workers and will update the links automatically without user confirmation
wherever possible. 

## Alternatives

1. [Linkchecker](https://github.com/linkchecker/linkchecker) - rich on features and usable on running websites.
2. [Linkcheck](https://github.com/filiph/linkcheck) - fast and easy alternative to _Linkchecker_
