#!/bin/sh

. ./helper.sh

config=zX0ysE.json
nodever="./nodever -config $config"
cwd=`pwd`

oneTimeSetUp()
{
	go build ../nodever.go
}

oneTimeTearDown()
{
	rm nodever
	rm -f $config
}

setUp()
{
	unset NODEVER
	config_default
}

test_first_run()
{
	rm -f $config
	assert_match_exec "$nodever" "nodever error: cannot find node"
	assertEquals 66 $?
}

test_env_var()
{
	NODEVER='{"dir":"/opt/s","def":"some-node-version"}' \
		   assert_match_exec "$nodever" '^some-node-version \(NODEVER env var\)$'
}

test_mode_info()
{
	assert_match_exec "$nodever" "^node-0.12.0 \(.+/zX0ysE.json\)$"
}

test_mode_init()
{
	rm -f $config
	local r=`$nodever init`
	assertEquals 0 $?

	assertEquals '{"Dir":"/opt/s","Def":"SET ME"}' "`cat $config`"
}

test_mode_list_fail()
{
	rm -f $config
	assert_match_exec "$nodever list" "nodever error: cannot read config file"
}

test_mode_list()
{
	local r=`$nodever list`
	assertEquals `echo "$r" | wc -l` 4
	assertEquals "`echo "$r" | grep 'node-0.12.0'`" '* node-0.12.0'
	assertEquals "`echo "$r" | grep 'iojs-1.2.3'`" '  iojs-1.2.3'
}

test_mode_use_fail_no_config()
{
	rm -f $config
	assert_match_exec "$nodever use" "nodever error: cannot read config file"
}

test_mode_use_fail()
{
	assert_match_exec "$nodever use" "the query must resolve in 1 entry"
	assert_match_exec "$nodever use o" "the query must resolve in 1 entry"
	assertEquals 1 $?
}

test_mode_use()
{
	$nodever use iojs
	assertEquals "{\"Dir\":\"$cwd/data/node\",\"Def\":\"iojs-1.2.3\"}" \
				 "`cat $config`"
}

. /usr/share/shunit2/shunit2
