FROM golang:latest AS builder
WORKDIR /fluxcd-addon
COPY . .
ENV GO_PACKAGE github.com/kluster-management/fluxcd-addon

# Build
RUN make build --warn-undefined-variables

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM alpine:latest

# Add the binaries
WORKDIR /addon/
COPY --from=builder /fluxcd-addon/bin/ .
CMD ["./fluxcd-addon", "manager"]