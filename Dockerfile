#force rebuild docker build --no-cache -t tdewin/mysql-employees:latest .
FROM golang AS compiler
ENV CGO_ENABLED=0
RUN go get -d github.com/tdewin/mysql-employees && go install github.com/tdewin/mysql-employees@latest && chmod 755 /go/bin/mysql-employees

FROM alpine
LABEL maintainer="@tdewin"
WORKDIR /usr/sbin/
COPY --from=compiler /go/bin/mysql-employees /usr/sbin/mysql-employees
EXPOSE 8080
CMD /usr/sbin/mysql-employees