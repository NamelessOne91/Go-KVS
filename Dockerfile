##
## STEP 1 - BUILD
##

FROM golang:1.20.2-alpine3.17 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o /go-kvs

##
## STEP 2 - DEPLOY
##

FROM scratch

WORKDIR /

COPY --from=build /go-kvs /go-kvs

EXPOSE 8080

ENTRYPOINT ["/go-kvs"]
