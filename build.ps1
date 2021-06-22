Set-Location $env:GOPATH\src\github.com\tdewin\mysql-employees
docker build -t tdewin/mysql-employees .
docker push tdewin/mysql-employees