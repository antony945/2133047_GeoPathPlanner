FROM golang:1.25.1-alpine AS dev
WORKDIR /app

# Install dependencies for CGO + GEOS
RUN apk add --no-cache alpine-sdk geos geos-dev

# Install Go dependencies
COPY go.* ./
RUN go mod download

# Copy source
COPY . .

# This is the image VS Code will use
# gopls, autocompletion etc. will work here
CMD ["sleep", "infinity"]

# Build stage
FROM dev AS build
RUN go build -o main cmd/api/main.go

FROM alpine:3.20.1 AS prod

# Install runtime deps (only GEOS runtime, not dev tools)
RUN apk add --no-cache geos

WORKDIR /app
COPY --from=build /app/main /app/main

CMD ["./main"]