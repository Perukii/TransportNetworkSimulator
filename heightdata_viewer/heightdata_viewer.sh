#gcc -o main main.c `pkg-config --cflags --libs gdk-3.0`
go build main.go
./main $1 $2 $3 $4
rm -r ./main