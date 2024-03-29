FROM golang:1.21.5 as builder

WORKDIR /app

COPY . .

RUN go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o uniTranslate-linux-amd64 main.go

FROM bitnami/git:latest
FROM node:18.19.0 as console

ARG CACHEBUST=1

WORKDIR /app

RUN git clone https://github.com/xgd16/UniTranslate-web-console.git console

WORKDIR /app/console

RUN npm install && npm run build

FROM alpine:latest

WORKDIR /app

COPY --from=console /app/console/dist ./dist
COPY --from=builder /app/uniTranslate-linux-amd64 .
COPY --from=builder /app/translate.json .

CMD ["./uniTranslate-linux-amd64"]

# docker build --no-cache -t uni-translate:latest .

# docker run -d --name uniTranslate -v /Users/x/docker/uniTranslate/config.yaml:/app/config.yaml -p 9431:9431 --network baseRun uni-translate:latest