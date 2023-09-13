go build -o gogithub_linux_amd64 main.go

cd ./init

go build -o init_linux_amd64 main.go

go env -w CGO_ENABLED=0
go env -w GOOS=windows
go env -w GOARCH=amd64
go build -o init_windows_amd64.exe main.go

go env -w CGO_ENABLED=0
go env -w GOOS=darwin
go env -w GOARCH=amd64
go build -o init_darwin_amd64 main.go

cd ..

go env -w CGO_ENABLED=0
go env -w GOOS=windows
go env -w GOARCH=amd64
go build -o gogithub_windows_amd64.exe main.go

go env -w CGO_ENABLED=0
go env -w GOOS=darwin
go env -w GOARCH=amd64
go build -o gogithub_darwin_amd64 main.go

go env -w GOOS=linux