#! /usr/bin/env sh

cat README.md | tee PRINTOUT.md
echo "" | tee -a PRINTOUT.md
for f in *.go; do
    echo "#### $f" | tee -a PRINTOUT.md
    echo "" | tee -a PRINTOUT.md
    echo '```' | tee -a PRINTOUT.md
    cat $f | tee -a PRINTOUT.md
    echo '```' | tee -a PRINTOUT.md
done
