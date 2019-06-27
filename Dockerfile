FROM silverstagtech/distroless:latest

LABEL maintainer="Randy Coburn <morfien101@gmail.com>"

COPY ./artifacts/* /
ENV METRIC_AUTH_CONFIG=/metric-auth.conf

ENTRYPOINT [ "/metrics-auth" ]