# syntax=docker/dockerfile:2
FROM golang:latest AS build_base

#setting working directory to /app
WORKDIR /app

#copying all source to working dir
COPY . .

#run go mod download and go mod verify
RUN go mod download && go mod verify

#run go build
RUN CGO_ENABLED=0 GOOS=linux go build -o disburser .

#stage two
FROM alpine:latest

#set working dir
WORKDIR /app

COPY . .
COPY .env .
ENV ENV=ENV_DEV

#copy bin from build stage
COPY --from=build_base /app/disburser .

#set the timezone and install CA certificates
RUN apk --no-cache add ca-certificates tzdata


#expose required ports
EXPOSE 80
EXPOSE 443

ENTRYPOINT ["/app/disburser"]