FROM scratch

USER 1000

COPY --chown=1000 --chmod=755 ./build/snooze-otlp /snooze-otlp

CMD ["/snooze-otlp", "run"]
