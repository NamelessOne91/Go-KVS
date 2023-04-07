##
## STEP 1 - Compile binaries
##

FROM golang:1.20.2-alpine3.17 AS build

COPY . /app

WORKDIR /app/src

RUN CGO_ENABLED=0 GOOS=linux go build -o go-kvs

##
## STEP 2 - Build image
##

FROM scratch

COPY --from=build /app/src/go-kvs .
COPY --from=build /app/*.pem .

EXPOSE 8080

ENTRYPOINT ["/go-kvs"]
