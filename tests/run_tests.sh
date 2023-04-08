#!/bin/sh

for i in 1 2 3 4 5 6
do 
   echo "Testing $i"

   ./typegen -f tests/test$i.json -o tests/test$i.go_test
   diff tests/test$i.go_test tests/test$i.go_benchmark 

   if [ $? != 0 ]; then
      echo "Error on test for Go: $i"
      exit 1
   else
      rm tests/test$i.go_test
   fi

   ./typegen -f tests/test$i.json -o tests/test$i.swift_test -l swift
   diff tests/test$i.swift_test tests/test$i.swift_benchmark 

   if [ $? != 0 ]; then
      echo "Error on test for Swift: $i"
      exit 1
   else
      rm tests/test$i.swift_test
   fi


done

