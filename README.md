# Linkfix

Linkfix is a simple tool that helps you avoid the [Link rot](https://en.wikipedia.org/wiki/Link_rot)
by reporting on no-longer working links in your files and suggesting replacements with
[Wayback Machine snapshots](https://archive.org/web/) wherever possible.

## Features

1. Automatic suggestions for replacement by snapshots from [Wayback Machine](https://archive.org/web/).
2. Multithreaded with configurable level of concurrency.
3. Can be used both as a tool or through Docker.
4. (working on) Multiple output formats.

## Link rot

> Link rot is the phenomenon of hyperlinks tending over time to cease to point
> to their originally targeted file, web page, or server due to that resource
> being relocated or becoming permanently unavailable. A link that no longer points to its target,
> often called a broken or dead link, is a specific form of dangling pointer. 
> The rate of link rot is a subject of study and research due to its significance 
> to the internet's ability to preserve information. Estimates of that rate vary dramatically
> between studies.

## Installation

```shell script
go get -u github.com/matoous/linkfix
```

or run using the docker image without installation:

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
