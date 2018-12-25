FROM fedora:27

# basics
RUN dnf -y update; dnf clean all

# fundamental pkg
RUN dnf -y install file gcc make man sudo tar git; dnf clean all

# golang (install)
ENV APPBIN="/go/bin" \
    APPDIR="/go/src/github.com/rvflash/safe" \
    GOPATH="/go" \
    GOZIP="go1.11.4.linux-amd64.tar.gz" \
    MINGW_PATH="/usr/x86_64-w64-mingw32/sys-root/mingw" \
    PATH="$PATH:/usr/local/go/bin"

RUN curl -O -s https://dl.google.com/go/$GOZIP && \
    tar -xzf $GOZIP -C /usr/local && rm $GOZIP && \
    mkdir -p "$APPBIN/"{darwin,linux,windows} && mkdir -p "$APPDIR"
ADD . "$APPDIR"
VOLUME "$APPDIR"

# mingw64 & dependecies
RUN dnf -y install \
    mingw64-gcc \
    mingw64-cairo mingw64-cairo-static cairo cairo-devel \
    mingw64-freetype freetype freetype-devel \
    mingw64-gtk3 gtk3 gtk3-devel \
    mingw64-glib2-static glib2 glib2-devel \
    mingw64-pango pango pango-devel

#RUN dnf -y install mingw64-poppler poppler poppler-devel
#RUN dnf -y install mingw64-winpthreads mingw64-winpthreads-static
#RUN dnf -y install mingw64-harfbuzz harfbuzz harfbuzz-devel
#RUN dnf -y install atk atk-devel

# gtk3 (useful: https://fedoraproject.org/wiki/Packaging:MinGW)
WORKDIR "/go"
RUN go get github.com/gotk3/gotk3/gtk

ENV CC=x86_64-w64-mingw32-gcc \
    CGO_ENABLED=1 \
    CGO_LDFLAGS_ALLOW="-Wl,-luuid" \
    GOOS=windows \
    GOARCH=amd64 \
    PKG_CONFIG_PATH="$MINGW_PATH/lib/pkgconfig"

RUN go install github.com/gotk3/gotk3/gtk

# safe (on windows)
# to copy only ddl outside, get the CONTAINER_ID and copy them:
# ~ docker ps -alq
# ~ docker cp <CONTAINER_ID>:<MINGW_PATH>/bin cmd/safe/bin/windows
WORKDIR "$APPDIR/cmd/safe"

RUN cp -ra "$MINGW_PATH/bin/." "$APPBIN/windows/" && \
    cp -ra "$MINGW_PATH/share/icons" "$APPBIN/windows/" && \
    GO111MODULE=on \
    go build -ldflags "-H windowsgui" -o "$APPBIN/windows/safe.windows.amd64.exe"

# releases
WORKDIR "$APPBIN"
RUN echo "#!/bin/sh -e" > build.sh && \
    echo "tar -zcvf $APPDIR/safe.windows.amd64.tar.gz windows" >> build.sh && \
    chmod +x build.sh

CMD ["./build.sh"]

# safe (on linux)
#GOOS=linux CC=clang CXX=clang++ go build -o "$APPBIN/linux/safe.linux.amd64"

# To build it:
# ~ docker build -t safe .
# And create Zip files:
# ~ docker run -v "$(pwd)":/go/src/github.com/rvflash/safe safe