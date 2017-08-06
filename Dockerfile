FROM golang:alpine as build

RUN apk add --no-cache git
WORKDIR /go/src/sample-docker-plugin
COPY . .

RUN go-wrapper download
RUN go-wrapper install


FROM alpine:3.6

COPY --from=build /go/bin/sample-docker-plugin /usr/local/bin/

CMD ["sample-docker-plugin"]