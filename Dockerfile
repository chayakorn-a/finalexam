FROM centurylink/ca-certs
COPY ./main /work/main
EXPOSE 2019
WORKDIR /work
CMD ["main"]