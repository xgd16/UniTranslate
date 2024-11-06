FROM golang:1.23.2 as builder
WORKDIR /workspace
COPY . .
RUN go mod tidy && CGO_ENABLED=0 go build -o app main.go

FROM --platform=amd64 node:18.20.4 as console
WORKDIR /workspace
ARG CACHEBUST=1
RUN git clone https://github.com/xgd16/UniTranslate-web-console.git console
WORKDIR /workspace/console
RUN npm install && npm run build

FROM alpine:latest
WORKDIR /workspace
COPY --from=console /workspace/console/dist ./dist
COPY --from=builder /workspace/app .
COPY --from=builder /workspace/translate.json .
ENTRYPOINT [ "./app" ]