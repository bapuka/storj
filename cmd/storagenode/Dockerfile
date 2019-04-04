ARG TAG
ARG DOCKER_ARCH amd64
FROM alpine:${DOCKER_ARCH}
ENV CONF_PATH=/root/.local/share/storj/storagenode \
    IDENT_PATH=/root/.local/share/storj/identity/storagenode \
    CERT_DIFFICULTY=30 \
    SATELLITE_ADDR=
EXPOSE 28967
WORKDIR /app
VOLUME /root/.local/share/storj/storagenode
COPY release/${TAG}/storagenode_linux_${GOARCH} /app/storagenode
COPY --from=storagenode-build-env /app /app/storagenode
COPY cmd/storagenode/entrypoint /entrypoint
ENTRYPOINT ["/entrypoint"]
