FROM scratch
COPY cyberark-to-k8s /
ENTRYPOINT ["/cyberark-to-k8s"]
