FROM --platform=linux/amd64 golang:1.19 AS builder

ARG VERSION="1.0.0" 

WORKDIR /app

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64

COPY ["go.mod", "./"]

ADD . .

RUN go build .
RUN go test -c vm_go/vending_machine/tests -o vm_go_test

FROM --platform=linux/amd64 gcr.io/distroless/static AS runner

ARG RELEASE

COPY --from=builder /app/vm_go .
COPY --from=builder /app/vm_go_test ./vm_go_test

CMD ["/vm_go"]
