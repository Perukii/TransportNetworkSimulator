
#gcc -o mapgen mapgen.c -lpng
go build main.go
./main $1 $2 $3 $4 $5 $6
rm -r main