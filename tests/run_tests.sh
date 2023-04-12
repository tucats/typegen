#!/bin/sh

for i in $(ls -1 tests/test*.json)

do 
   echo "Testing $i"

   RUN=$(dirname $i)/$(basename $i .json).go_test
   BENCH=$(dirname $i)/$(basename $i .json).go_benchmark

   ./typegen -f $i -o $RUN
   diff $RUN $BENCH

   if [ $? != 0 ]; then
      echo "Error on test for Go: $i"
   else
      rm $RUN
   fi

   RUN=$(dirname $i)/$(basename $i .json).swift_test
   BENCH=$(dirname $i)/$(basename $i .json).swift_benchmark

   ./typegen -f $i -o $RUN -l swift
   diff $RUN $BENCH

   if [ $? != 0 ]; then
      echo "Error on test for Swift: $i"
   else
      rm $RUN
   fi

done

