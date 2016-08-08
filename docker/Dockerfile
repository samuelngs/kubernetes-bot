FROM scratch

ADD docker/ca-certificates.crt /etc/ssl/certs/
ADD bin/kubebot /

ENTRYPOINT ["/kubebot"]
