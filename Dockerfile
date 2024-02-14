FROM docker.io/library/golang:1.21 as builder

COPY . /work/
WORKDIR /work
RUN make

FROM gcr.io/distroless/static-debian12:latest

COPY --from=0 /work/build/ms1-server /

EXPOSE 50051

ENTRYPOINT ["/ms1-server"]
CMD ["run"]
