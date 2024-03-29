set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -o bootstrap main.go
wsl ./build-lambda-zip bootstrap