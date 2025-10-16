FROM golang:1.23 as builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o hrms ./main.go

FROM gcr.io/distroless/base-debian12
WORKDIR /app
ENV SERVER_PORT=8082
COPY --from=builder /app/hrms /app/hrms
COPY --from=builder /app/docs /app/docs
EXPOSE 8082
USER 65532:65532
ENTRYPOINT ["/app/hrms"]


