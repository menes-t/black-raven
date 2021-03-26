ARG  BUILDER_IMAGE=golang:buster
ARG  DISTROLESS_IMAGE=gcr.io/distroless/base


FROM ${BUILDER_IMAGE} as builder

RUN update-ca-certificates

WORKDIR /app

COPY go.mod .

ENV GO111MODULE=on
RUN go mod download
RUN go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/black_raven .

FROM ${DISTROLESS_IMAGE}

COPY --from=builder /go/bin/black_raven /go/bin/black_raven
COPY --from=builder /app/resources /app/resources

ENTRYPOINT ["/go/bin/black_raven"]