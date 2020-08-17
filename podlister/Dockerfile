FROM golang:1.14-alpine AS build

WORKDIR /src/
COPY main.go go.* /src/
RUN CGO_ENABLED=0 go build -o /bin/podlister

FROM alpine
WORKDIR /go
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=build /bin/podlister /go/podlister
COPY index.template /go/
ENTRYPOINT ["/go/podlister"]
