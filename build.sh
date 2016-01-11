#!/bin/bash
set -e

artifact_dir=package
package_contents=package_contents
test_suite_dir=test-suite

function compile(){
    code_dir=$1
    echo "== Building [$code_dir] =="
    pushd $code_dir
    pwd
    make
    popd
    echo "=========================="
}

function run_test(){
    solution=$1
    echo "-- Starts server [$solution] --"
    pushd $solution
    make run &
    sleep 2
    popd
    echo "-- Running tests --"
    ./test-suite/do-package-tree
    kill `pidof make`
    echo "---------------------------------------"
}

rm -rf $artifact_dir
mkdir $artifact_dir

rm -rf $package_contents
mkdir $package_contents

compile test-suite

compile ruby-solution
run_test ruby-solution

#build package
echo "****************** Packaging ******************"
cp INSTRUCTIONS.md $package_contents/
cp $test_suite_dir/do-package-tree  $package_contents/
tar -cvzf $package_contents/source.tar.gz $test_suite_dir/*go

tar -cvzf $artifact_dir/candidate.tar.gz $package_contents 

rm -rf $package_contents
echo "***********************************************"
