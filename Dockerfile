#
#  Builder
#
FROM golang:1.17-alpine as builder

WORKDIR /src/
COPY . /src/
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0  go build -o ./bin/mytool

#
# CLI
#
FROM amazon/aws-cli:2.4.1

COPY --from=builder ./src/bin/mytool /bin/mytool
ENTRYPOINT [ "/bin/mytool" ]
