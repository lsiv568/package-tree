all: dependencies compile test

compile:
	rspec

test:
	ruby lib/server.rb &
	../test-suite/do-package-tree

dependencies:
	bundle update
