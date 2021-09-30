#!/bin/env bash

for i in $(ls *yaml)
do
    kubectl apply -f $i
done
