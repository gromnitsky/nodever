#!/bin/sh

. ./helper.sh

node="./node"

oneTimeSetUp()
{
	go build ../node.go
}

oneTimeTearDown()
{
	rm node
	rm -f $config
}

setUp()
{
	export NODEVER_WRAPPER="-config $config"
	config_default
}

test_first_run()
{
	rm $config
	assert_match_exec "$node -v" "nodever-wrapper error: cannot find node"
}

. /usr/share/shunit2/shunit2
