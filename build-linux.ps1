Set-Location $env:GOPATH\src\github.com\tdewin\mysql-employees
$env:GOOS="linux"
$env:GOARCH="amd64"
go build github.com/tdewin/mysql-employees

docker image rm tdewin/mysql-employees
docker build -t tdewin/mysql-employees .
docker push tdewin/mysql-employees