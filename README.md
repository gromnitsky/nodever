# nodever

A minimalistic node version manager.

## Features

* Per user or per project defaults.
* Wrappers around `node` & `npm` commands.
* Doesn't poison bash or any other shell.
* Works w/ any node/iojs versions.
* Written in Go.

## Requirements

* Linux
* Go 1.4.1

## Overview

nodever _doesn't_ know how to install any of node/iojs versions. You
either grab binaries or compile node by yourself & place the result
under 1 directory, say `/opt/s`. For example:

	$ ls -1d /opt/s/{node,iojs}*
	/opt/s/iojs-v1.1.0-linux-ia32
	/opt/s/node-v0.10.35-linux-x86
	/opt/s/node-v0.12.0-linux-x86

After installing nodever, you'll get 2 _wrappers_ for node/npm
binaries. Those wrappers check 2 places for instructions:

1. `NODEVER` env variable
2. `.nodever.json` file

The search for the `.nodever.json` begins form the current directory
(_not_ the directory of .js file), then switches to its parent & all way
down to the root. This allows us to have `~/.nodever.json` file w/ a
default settings for a particular user & a completely different
configuration for `$HOME/some/project`.

## Installation

	$ go get github.com/gromnitsky/nodever/...

(Yes, those `...` are required.)

This implies that you have a working Go installation & know what
`GOPATH` is.

## Setup

In `$GOPATH/bin` you'll find 3 new executables:

	nodever
	npm
	node

Run nodever for the first time:

	$ nodever
	nodever error: cannot find node; rerun w/ '-v 1' argument; \
		for help see https://github.com/gromnitsky/nodever

This means no configuration was found. Type:

	$ nodever -u init

Which should not raise an error. Now

	$ cat ~/.nodever.json
	{"Dir":"/opt/s","Def":"SET ME"}

Means that **you must manually** set "Dir" to a directory where
different node versions are installed. (Don't touch "Def".)

	$ nodever list
	(/home/alex/.nodever.json)
	  iojs-v1.1.0-linux-ia32
	  node-v0.10.35-linux-x86
	  node-v0.12.0-linux-x86

Gives the source of the configuration & all available installations. To
select node 0.12.0 type:

	$ nodever use 0.12

Which again should not raise an error. Now

	$ nodever list
	(/home/alex/.nodever.json)
	  iojs-v1.1.0-linux-ia32
	  node-v0.10.35-linux-x86
	* node-v0.12.0-linux-x86

Shows that `node-v0.12.0-linux-x86` is indeed a default. The same info
can be viewed by

	$ nodever
	node-v0.12.0-linux-x86 (/home/alex/.nodever.json)

Finally, the most interesting part:

	$ node -v
	v0.12.0

What was that? The `node` wrapper (`$GOPATH/bin/node`) found
`$HOME/.nodever.json` file & executed a
`/opt/s/node-v0.12.0-linux-x86/bin/node` binary w/ `-v` option. The
wrapper passed all CL options directly to the 'wrapped' command. It
connected its stdin/stdout/stderr to the corresponding streams of the
wrapped binary & it returned the exit code from the wrapped binary, so
everything worked just as if you did run the real node executable
directly.

## Shell scripts

There is no need to modify or create `.nodever.json` files if you want
select a node version inside a shell script. Use `NODEVER` env variable
for that. For example:

	$ cat my-script.sh
	#!/bin/sh
	export NODEVER='{"Dir":"/opt/s","Def":"iojs-v1.1.0-linux-ia32"}'
	node -v

	$ ./my-script.sh
	v1.1.0

## Bugs

* Tested on Fedora 21 only.
* Probably won't work under Windows.

## TODO

* Add `nodever exec 0.12 node foo.js` as a shortcut for setting
  `NODEVER`.

## License

MIT.

