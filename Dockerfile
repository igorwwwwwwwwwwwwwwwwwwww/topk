FROM golang:1.25

WORKDIR /app
COPY . ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /topk

ENTRYPOINT ["/topk"]
