#!/usr/bin/env bash

mkdir -p workspace/n0/binary
mkdir -p workspace/n1/binary
mkdir -p workspace/n2/binary
mkdir -p workspace/n3/binary
mkdir -p workspace/n4/binary

toolbox build workspace/n0/binary --git-branch wtx-test
toolbox build workspace/n1/binary --git-branch wtx-test
toolbox build workspace/n2/binary --git-branch wtx-test
toolbox build workspace/n3/binary --git-branch wtx-test-1.1.4.3-log
toolbox build workspace/n4/binary --git-branch wtx-test