

gcc -o mapgen mapgen.c -lpng
./mapgen $1 $2 $3
rm -r mapgen