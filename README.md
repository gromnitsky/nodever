# nodever

A minimalistic node version manager.

## Features

* Per user or per project defaults.
* Wrappers around `node` & `npm` commands.
* Doesn't poison bash or any other shell.
* Works w/ any node/iojs versions.
* Über fast.

## Requirements

* Linux
* Go 1.17.1

## Overview

nodever _doesn't_ know how to install any of node/iojs versions. You
either grab binaries or compile node by yourself & place the result
under 1 directory, say `/opt/s`. For example:

	$ ls -1d /opt/s/{node,iojs}*
	/opt/s/iojs-v1.1.0-linux-ia32
	/opt/s/node-v0.10.35-linux-x86
	/opt/s/node-v0.12.0-linux-x86

After installing nodever, you'll get _wrappers_ for node/npm/npx/corepack
binaries. Those wrappers check 2 places for instructions:

1. `NODEVER` env variable
2. `.nodever.json` file

The search for the `.nodever.json` begins form the current directory
(_not_ the directory of a .js file), then switches to its parent & all
the way down to /. This allows us to have `~/.nodever.json` file w/ a
default settings for a particular user & a completely different
configuration for `~/some/project`.

## Installation

~~~
$ git clone https://github.com/gromnitsky/nodever
$ cd nodever
$ go install ./...
~~~

Yes, `./...` is required verbatim. You may also remove the cloned dir
afterwards.

## Setup

In `~/go/bin` you'll find 5 new executables:

~~~
corepack
node
nodever
npm
npx
~~~

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

What was that? The `node` wrapper (`~/go/node`) found
`$HOME/.nodever.json` file & executed a
`/opt/s/node-v0.12.0-linux-x86/bin/node` binary w/ `-v` option. The
wrapper passed all CL options directly to the 'wrapped' command. It
connected its stdin/stdout/stderr to the corresponding streams of the
wrapped binary & it returned the exit code from the wrapped binary, so
everything worked just as if you run the real node executable
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

## Hints

You can run an arbitrary command in a subshell with the desired node
version:

	$ nodever exec 0.10 coffee
	coffee> process.version
	'v0.10.36'

	$ nodever exec iojs npm -v
	2.5.1

## NPM

One annoyance you'll get is a weird path for globally installed
packages.

Suppose we set 0.10 as a default, chose `/opt/s` as the node umbrella
directory & we expect global packages to be installed in `/opt/lib`. But

	# npm ls --depth=0 -g
	/opt/s/node-v0.10.35-linux-x86/lib
	└── npm@1.4.28

Huh? It turns out, npm uses a prefix for modules path that depends on
`process.execPath`, which in our case is
`/opt/s/node-v0.10.35-linux-x86/bin/node`. So that

	# node
	> path.dirname(path.dirname(process.execPath))
	'/opt/s/node-v0.10.35-linux-x86'

To fix it we may create a global npm config

	# cat /etc/npmrc
	prefix = /opt

& force npm to read it by exporting `npm_config_globalconfig` env
variable.

	# export npm_config_globalconfig=/etc/npmrc
	# npm ls --depth=0 -g
	/opt/lib
	├── ...
	└── ...

## Bugs

* Tested on Fedora only.
* Probably won't work under Windows.

## License

MIT.
