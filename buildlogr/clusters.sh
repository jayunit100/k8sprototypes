#!/bin/bash

for filename in *.txt*; do
	echo "\n\n CLUSTERS IN $filename \n\n"
	echo "\n\n\n"
	cat $filename |grep msg | grep -v hello | grep -v echo | grep FAILED -B 4 -A 4 \
		| sed 's/FAILED/\n********\n******** FAILED/g' > ${filename}.0
	
	cat ${filename}.0 | sed '/FAILED /a********\n'  > FAILURES
done


echo "COUNTING LOCALIZED FAILURES AND SORTING BY FREQUENCY......................"
cat FAILURES | cut -d' ' -f 3-100 | sort | uniq -c | sort
