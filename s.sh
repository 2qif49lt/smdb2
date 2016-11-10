#!/bin/sh
 

for ((beg=10; beg>0; --beg))  
do  
    echo "begin after $beg"
    sleep 1
done 

cat README.md

for ((end=0; end<10; ++end))  
do  
    echo "end after $end"
    sleep 1
done 


echo 'end'