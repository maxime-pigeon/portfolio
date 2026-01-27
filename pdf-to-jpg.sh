#!/bin/bash

i=1
for pdf in $1*.pdf; do
    [ -e "$pdf" ] || continue
    echo "Converting ${pdf}"
    vips copy "${pdf}[dpi=300]" "$1image${i}.jpg"
    ((i++))
done
