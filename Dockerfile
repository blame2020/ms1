FROM gcr.io/distroless/static-debian12:latest

COPY ./build/ms1-server /

EXPOSE 50051

ENTRYPOINT ["/ms1-server"]
CMD ["run"]
