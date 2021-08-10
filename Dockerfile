FROM debian:latest
RUN apt-get update && apt-get upgrade -y
RUN apt-get install -y wget xz-utils
RUN mkdir -p server/static
WORKDIR /server
RUN wget -q https://ziglang.org/download/0.8.0/zig-linux-x86_64-0.8.0.tar.xz && \
    tar xf zig-linux-x86_64-0.8.0.tar.xz && \
    mv zig-linux-x86_64-0.8.0/zig /usr/local/bin && \
    mkdir -p /usr/local/bin/lib && \
    mv zig-linux-x86_64-0.8.0/lib/* /usr/local/bin/lib && \
    rm -rf zig-linux-x86_64-0.8.0*
COPY zig-play .
COPY static/ static/
ENTRYPOINT ./zig-play
