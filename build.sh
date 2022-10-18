rm -r ./bin
mkdir ./bin
cd ./src || exit
go build -o ../bin/wafflesyrup ./wafflesyrup.go
cd ..