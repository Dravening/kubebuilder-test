FROM katanomi/distroless-static:nonroot
WORKDIR /
COPY ./manager .
USER 65532:65532

ENTRYPOINT ["/manager"]