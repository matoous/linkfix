FROM golang:1.14.4@sha256:ebe7f5d1a2a6b884bc1a45b8c1ff7e26b7b95938a3e8847ea96fc6761fdc2b77 AS deps-cached

ARG PROJECT_PATH=/linkfix
WORKDIR ${PROJECT_PATH}

# Copy and install dependencies
COPY Makefile go.mod go.sum ${PROJECT_PATH}/
RUN make configure

# Start stage for with all files for building various images
FROM deps-cached as builder

COPY . .

ARG BUILD_REV
ARG COMMIT_DATE

RUN CGO_ENABLED=0 BUILD_REV="${BUILD_REV}" COMMIT_DATE="${COMMIT_DATE}" make build
RUN mkdir -p /build && cp bin/linkfix /build/linkfix

# Last stage with binary only
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/linkfix /

ENTRYPOINT ["/linkfix"]
LABEL name=linkfix
