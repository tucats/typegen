#!/bin/sh

./typegen -f tests/test1.json -o tests/test1.go_benchmark
./typegen -f tests/test2.json -o tests/test2.go_benchmark
./typegen -f tests/test3.json -o tests/test3.go_benchmark
./typegen -f tests/test4.json -o tests/test4.go_benchmark
./typegen -f tests/test5.json -o tests/test5.go_benchmark
./typegen -f tests/test6.json -o tests/test6.go_benchmark


./typegen -f tests/test1.json -o tests/test1.swift_benchmark -l swift
./typegen -f tests/test2.json -o tests/test2.swift_benchmark -l swift
./typegen -f tests/test3.json -o tests/test3.swift_benchmark -l swift
./typegen -f tests/test4.json -o tests/test4.swift_benchmark -l swift
./typegen -f tests/test5.json -o tests/test5.swift_benchmark -l swift
./typegen -f tests/test6.json -o tests/test6.swift_benchmark -l swift
