#! /usr/bin/env sh

cat README.md | tee PRINTOUT.md
echo "" | tee -a PRINTOUT.md
for f in *.go; do
    if [ "$f" != "conn_test.go" ]; then
        echo "#### $f" | tee -a PRINTOUT.md
        echo "" | tee -a PRINTOUT.md
        echo '```' | tee -a PRINTOUT.md
        cat $f | tee -a PRINTOUT.md
        echo '```' | tee -a PRINTOUT.md
    fi
done
