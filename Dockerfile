FROM golang:1.10.3-stretch
WORKDIR /go/src/kubernetes-pod-scheduler
RUN apt-get update && apt-get install curl -y
RUN curl https://glide.sh/get | sh
COPY glide.yaml .
COPY glide.lock .
RUN glide install
COPY main.go .
# RUN go build main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/kubernetes-pod-scheduler/main .
CMD ["./main"]
