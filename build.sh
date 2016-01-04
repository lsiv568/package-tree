#!/bin/bash

artifact_dir=package
candidate_package_dir=do-package-tree
test_suite_dir=test-suite

rm -rf $artifact_dir
mkdir $artifact_dir

rm -rf $candidate_package_dir
mkdir $candidate_package_dir

#compile test-suite
pushd $test_suite_dir
docker run -v=$PWD:$PWD -w=$PWD google/golang make
popd


# compile ruby solution
pushd $test_suite_dir
docker run -v=$PWD:$PWD -w=$PWD ruby make
popd


#build package
cp INSTRUCTIONS.md $candidate_package_dir/
cp $test_suite_dir/do-package-tree  $candidate_package_dir/
tar -cvzf $candidate_package_dir/source.tar.gz $test_suite_dir/*go

tar -cvzf $artifact_dir/candidate.tar.gz $candidate_package_dir 

rm -rf $candidate_package_dir
