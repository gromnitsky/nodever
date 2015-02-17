config=zX0ysE.json
cwd=`pwd`

assert_match_exec()
{
	local r
	r=`$1 2>&1`
	local exit_code=$?

	echo "$r" | egrep "$2" > /dev/null
	assertEquals 0 $?
	[ $? -ne 0 ] && {
		printf '`%s` does NOT match `%s`\n' "$2" "$r"
	}

	return $exit_code
}

config_default()
{
	echo "{\"dir\":\"$cwd/data/node\",\"def\":\"node-0.12.0\"}" > $config
}
