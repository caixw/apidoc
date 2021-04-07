# blogit dockerfile

FROM scratch

MAINTAINER caixw <https://caixw.io>

COPY ./blogit /

ENTRYPOINT ["/blogit"]
