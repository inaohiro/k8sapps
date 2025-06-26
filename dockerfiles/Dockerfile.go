FROM golang:1.24

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o app .

FROM scratch
COPY --from=0 /app/app .
CMD ["/app"]
