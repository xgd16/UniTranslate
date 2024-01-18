FROM golang:1.21.5 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o uniTranslate-linux-amd64 main.go

FROM bitnami/git:latest
FROM node:18.19.0 as console

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

# docker build -t uni-translate:latest .