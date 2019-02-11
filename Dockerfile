# Mulder Docker Image
# This is a 2-steps build using "multistage builds"
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/

###########################
# 1st step: build the app #
###########################
FROM golang:1.11 AS builder

# Copy the files required to build
COPY . /workspace/
WORKDIR /workspace

# Build a "static" binary - by disabling CGO
ENV CGO_ENABLED=0
RUN go build -o mulder

####################################
# 2nd step: define the "run" image #
####################################
# we want a small base
FROM alpine:3.9

# copy the pre-built binary
COPY --from=builder /workspace/mulder /usr/bin/mulder

ENTRYPOINT ["/usr/bin/mulder"]
