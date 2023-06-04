rm .\main.exe
go build -o main.exe -ldflags="-s -w" main.go