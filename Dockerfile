# base go image
FROM golang:1.22 as builder
RUN mkdir /app
WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download
COPY . /app
RUN CGO_ENABLED=0 go build -o app ./cmd/api
RUN chmod +x ./app

# build a tiny docker image
FROM alpine:latest
# make directory for product images
RUN mkdir -p /app /images/profilePics
COPY --from=builder /app/app /app
CMD [ "/app/app" ]