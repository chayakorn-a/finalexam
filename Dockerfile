FROM centurylink/ca-certs
COPY ./main /work/main
COPY ./main.sh /work/main.sh
EXPOSE 2019
WORKDIR /work
CMD ["./main.sh"]