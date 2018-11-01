FROM alpine
RUN apk update && \
    apk add gcc-avr avr-libc make
ADD drone-avr /bin/
ENTRYPOINT /bin/drone-avr
