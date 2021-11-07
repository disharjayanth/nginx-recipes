FROM golang:1.17
WORKDIR /Users/jayantha/Desktop/practise/nginx-recipes
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /Users/jayantha/Desktop/practise/nginx-recipes .
CMD ["./app"]