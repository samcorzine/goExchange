FROM scratch
ADD GoExchange /GoExchange
ENV GODEBUG http2debug=1

ENTRYPOINT ["/GoExchange"]
