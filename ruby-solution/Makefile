all: dependencies compile test

compile:
	rspec

test:
	ruby lib/server.rb &
	sleep 2
	../test-suite/do-package-tree

dependencies:
	bundle update
