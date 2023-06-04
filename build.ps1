# Remove old main.exe
if (Test-Path .\main.exe) { Remove-Item .\main.exe }

# Build new main.exe
go build -o main.exe -ldflags="-s -w" main.go