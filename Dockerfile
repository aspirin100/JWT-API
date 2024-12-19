FROM golang:1.23.3-alpine3.20 AS build

RUN apk --no-cache add make git

COPY . /go/src
WORKDIR /go/src

RUN make build


FROM alpine:3.20

COPY --from=build /go/src/bin/jwt-server /usr/local/bin/jwt-server

ENTRYPOINT ["jwt-server"]