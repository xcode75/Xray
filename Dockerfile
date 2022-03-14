# Build go
FROM golang:1.17.1-alpine AS builder
WORKDIR /app
COPY . .
ENV CGO_ENABLED=0
RUN go mod download
RUN go build -v -o Xray -trimpath -ldflags "-s -w -buildid=" ./main

# Release
FROM  alpine
RUN  apk --update --no-cache add tzdata ca-certificates \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

COPY geoip.dat  /usr/local/bin/geoip.dat
COPY geosite.dat  /usr/local/bin/geosite.dat
 	
RUN mkdir /etc/Xray/
COPY --from=builder /app/Xray /usr/local/bin

ENTRYPOINT [ "/usr/local/bin/Xray", "--config"]
CMD ["/etc/Xray/config.yml"]
