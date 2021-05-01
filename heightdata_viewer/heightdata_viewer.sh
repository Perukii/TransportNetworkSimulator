#gcc -o main main.c `pkg-config --cflags --libs gdk-3.0`
go build main.go
./main $1 $2 $3 $4 $5 $6 $7 $8 $9 ${10}
rm -r ./main