FROM centurylink/ca-certs
ADD main /
ADD main.sh /
CMD ["/main.sh"]