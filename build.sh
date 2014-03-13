#!/bin/bash

set -e

d="sudo docker"

$d build --rm=false -t dogestry .
id=$($d inspect dogestry | jq -r '.[0].container')
$d cp $id:dogestry .
mv dogestry dist/dogestry_`uname -i`

