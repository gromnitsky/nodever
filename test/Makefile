# Don't run make w/ -j flag

run := /usr/bin/time -f 'user=%U system=%S elapsed=%E'

# get test files in random order
test := $(shell shuf -e test_*.sh)

.PHONY: all
all: $(test)

.PHONY: $(test)
$(test):
	@printf '\n* %s\n' $@
	@$(run) ./$@
