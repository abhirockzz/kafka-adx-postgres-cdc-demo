FROM golang as build-stage
WORKDIR /app

COPY go.mod /app
COPY go.sum /app
RUN go mod download

COPY main.go /app

RUN cd /app && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o orders-gen

FROM alpine
COPY --from=build-stage /app/orders-gen /
RUN ls -lrt
CMD ["/orders-gen"]