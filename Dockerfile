FROM scratch
COPY delta /usr/bin/delta
ENTRYPOINT ["/usr/bin/delta"]
