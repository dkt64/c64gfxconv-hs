set GOOS=windows
set GOARCH=amd64
go build -o out/c64gfxconv-hs-win64.exe c64gfxconv-hs.go
set GOOS=linux
set GOARCH=amd64
go build -o out/c64gfxconv-hs-lin64 c64gfxconv-hs.go
set GOOS=linux
set GOARCH=arm64
go build -o out/c64gfxconv-hs-linux-arm64 c64gfxconv-hs.go
set GOOS=darwin
set GOARCH=amd64
go build -o out/c64gfxconv-hs-macos c64gfxconv-hs.go
