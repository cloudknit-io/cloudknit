FROM scratch
ENTRYPOINT ["/zlifecycle-il-operator"]
COPY ./templates /templates
COPY zlifecycle-il-operator /
