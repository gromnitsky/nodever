#!/usr/bin/env node

if (process.env.NODEVER) {
	console.log('NODEVER')
} else {
	console.log('config')
}

if (process.argv[2] === 'stop') process.exit(0)

var child_process = require('child_process')

var child = child_process.spawn('node', [__filename, 'stop'], {stdio: 'inherit'})

child.on('close', function (code) {
	process.exit(code)
})
