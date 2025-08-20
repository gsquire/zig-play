FROM debian:latest
ARG version=0.15.1
RUN apt-get update && apt-get upgrade -y
RUN apt-get install -y wget xz-utils
RUN mkdir -p server/static
WORKDIR /server
RUN wget -q https://ziglang.org/download/${version}/zig-x86_64-linux-${version}.tar.xz && \
    tar xf zig-x86_64-linux-${version}.tar.xz && \
    mv zig-x86_64-linux-${version}/zig /usr/local/bin && \
    mkdir -p /usr/local/bin/lib && \
    mv zig-x86_64-linux-${version}/lib/* /usr/local/bin/lib && \
    rm -rf zig-x86_64-linux-${version}*
COPY zig-play .
COPY static/ static/
RUN sed -i "s/###version###/${version}/" static/index.html
RUN groupadd -r run && \
    useradd -r -g run -s /usr/sbin/nologin runner && \
    mkdir playground && \
    chown -R runner:run playground
ENV PLAYGROUND_DIR=playground
USER runner
ENTRYPOINT ./zig-play
