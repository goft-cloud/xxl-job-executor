#!/bin/bash
echo "task begin"

i=1
while [ $i -le 100 ]
do
  echo "do task $i"
  sleep 1
  # shellcheck disable=SC2219
  let i++
done

echo "task done"