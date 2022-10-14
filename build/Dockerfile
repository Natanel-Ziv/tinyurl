FROM golang:1.19-bullseye as builder

ENV GO111MODULE=on

WORKDIR /tinyurl

COPY go* ./
RUN go mod download

COPY . ./

RUN go build -v -o app ./server

RUN chmod +x server

RUN cp /tinyurl/app /
RUN cp /tinyurl/server/server.env /tinyurl/

EXPOSE 8000

# Run the web service on container startup.
ENTRYPOINT ["/app"]