FROM busybox
WORKDIR /
USER root:root

COPY --from=bifrost-output /workspace/bifrost-ipam .
COPY install.sh .

RUN chmod +x install.sh
ENTRYPOINT ["/install.sh"]