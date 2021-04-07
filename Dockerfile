# apidoc dockerfile

FROM scratch

MAINTAINER caixw <https://caixw.io>

COPY ./apidoc /

ENTRYPOINT ["/apidoc"]
