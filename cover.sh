#!/bin/bash

cat coverage.out | \
awk 'BEGIN {cov=0; stat=0;} \
	$3!="" { cov+=($3==1?$2:0); stat+=$2; } \
    END {printf("%.2f%% statements\n", (cov/stat)*100);}'