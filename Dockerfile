############################
# STEP 1 build executable binary
############################

FROM golang:1.13-alpine AS builder

# Create appuser
ENV USER=machmum
ENV UID=10001
ENV GO111MODULE=on

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /usr/share/machmum/counter-queue
COPY . /usr/share/machmum/counter-queue

RUN export GO111MODULE=on
RUN export GOPROXY=direct
RUN export GOSUMDB=off

# build binary file
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -ldflags="-w -s" -o bin/counter-queue /usr/share/machmum/counter-queue/main.go

############################
# STEP 2 build a small image
############################
FROM ubuntu:18.04

# Import from builder.
COPY --from=builder /usr/share/machmum/counter-queue/bin /usr/share/machmum/counter-queue/bin/

WORKDIR /usr/share/machmum/counter-queue

# Create group and user to the group
RUN groupadd -r machmum && useradd -r -s /bin/false -g machmum machmum

# Set ownership golang directory
RUN chown -R machmum:machmum /usr/share/machmum/counter-queue

# Make docker container rootless
USER machmum