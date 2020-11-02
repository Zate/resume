FROM golang:alpine AS builder

RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

RUN apk add --no-cache ca-certificates git
WORKDIR /src
COPY main.go main.go
COPY go.* ./
RUN go mod tidy

RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /app .

FROM scratch AS final

WORKDIR /

COPY --from=builder /user/group /user/passwd /etc/

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app /app

# COPY site/ site/

VOLUME ["/site"]

EXPOSE 3000

# Perform any further action as an unprivileged user.
USER nobody:nobody

# Run the compiled binary.
ENTRYPOINT ["/app"]