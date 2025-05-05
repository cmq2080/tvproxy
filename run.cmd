set GOOS=windows
set GOARCH=amd64
go build -o .\deploy\tvproxy.exe
.\deploy\tvproxy.exe