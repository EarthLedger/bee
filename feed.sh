#!/bin/bash
idx = 0
while (( idx < 10 ))
do
	./dist/bee feed ./chunks/0009-50-69.json 10000000
	(( idx = $idx + 1 ))
done
