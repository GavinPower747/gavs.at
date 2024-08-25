FROM golang:1.23-bookworm as builder
WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN make compile ENVIROMENT=production

FROM debian:bookworm-slim
ARG GIT_COMMIT
ARG GIT_BRANCH
ARG BUILD_DATE

RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

ENV GIT_COMMIT=$GIT_COMMIT
ENV GIT_BRANCH=$GIT_BRANCH
ENV BUILD_DATE=$BUILD_DATE

COPY --from=builder /app/bin/server /app/bin/server

CMD ["/app/bin/server"]
