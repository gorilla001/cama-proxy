FROM alpine:3.5
COPY cama-proxy /root/cama-proxy
WORKDIR /root/
ENTRYPOINT ["./cama-proxy"]
