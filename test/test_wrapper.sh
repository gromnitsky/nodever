#!/bin/sh

. ./helper.sh

node="./node"

oneTimeSetUp()
{
	go build ../bin/node/node.go
}

oneTimeTearDown()
{
	rm node
	rm -f $config
}

setUp()
{
	unset NODEVER
	export NODEVER_WRAPPER="-config $config"
	config_default
}

test_first_run()
{
	rm $config
	assert_match_exec "$node -v" "nodever-wrapper error: cannot find node"
	assertEquals 66 $?
}

test_invalid_node_installation()
{
	echo "{\"dir\":\"$cwd/data/node\",\"def\":\"node-99.99.99-broken\"}" > $config
	assert_match_exec "$node -v" "fork/exec"
	assertEquals 65 $?
}

test_node()
{
	assert_match_exec "$node -v" "0.12.0"
	assertEquals 11 $?
}

test_node_how_env_config_overrides_another_config()
{
	NODEVER="{\"dir\":\"$cwd/data/node\",\"def\":\"iojs-1.2.3\"}" \
		   assert_match_exec "$node -v" "1.2.3"
	assertEquals 0 $?
}

. /usr/share/shunit2/shunit2
