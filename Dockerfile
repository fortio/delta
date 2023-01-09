FROM scratch
COPY delta /usr/bin/dnsping
ENTRYPOINT ["/usr/bin/delta"]
