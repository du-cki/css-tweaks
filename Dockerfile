FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG GIT_COMMIT
ARG GIT_BRANCH
ARG GIT_ORIGIN

RUN go build -ldflags \
    "-X 'main.CommitHash=$GIT_COMMIT' \
     -X 'main.Branch=$GIT_BRANCH' \
     -X 'main.Origin=$GIT_ORIGIN'" \
    -o server .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

COPY snippets/ ./snippets/
COPY templates/ ./templates/
COPY static/ ./static/

EXPOSE 8080

CMD ["./server"]