#!/bin/bash
formula_dir="$(brew --prefix)/Library/Formula/"
pushd $formula_dir
 ls -l  | awk '{print $9}' | cut -d. -f1 | sort| while read cask; do echo -n $cask:; brew deps $cask | awk '{printf(" %s ", $0)}'; echo ""; done 
popd
