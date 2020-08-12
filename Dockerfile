FROM golang:1 as builder
WORKDIR /app
COPY .  .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .
FROM alpine:latest  
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/app /bin/app
CMD ["/bin/app"]%