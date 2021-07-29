FROM golang:1.15-alpine3.12 as base
RUN apk add --no-cache build-base bash
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .

FROM base as build
RUN go build -o bin/drafts
CMD ["/app/bin/drafts"]