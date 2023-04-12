#!/bin/sh

for i in $(ls -1 tests/test*.json)

do 
   echo "Generating benchmarks for $i"

   BENCH=$(dirname $i)/$(basename $i .json).go_benchmark
   ./typegen -f $i -o $(dirname $i)/$(basename $i .json).go_benchmark -l go
   ./typegen -f $i -o $(dirname $i)/$(basename $i .json).swift_benchmark -l swift
done
