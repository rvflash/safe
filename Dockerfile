FROM fedora:37

# basics
RUN dnf -y update; dnf clean all

# fundamental pkg
RUN dnf -y install file gcc make man sudo git rsync; dnf clean all

# golang 1.19
ENV APPBIN="/go/bin" \
    APPDIR="/go/src/github.com/rvflash/safe" \
    GOPATH="/go" \
    GOZIP="go1.19.4.linux-amd64.tar.gz" \
    PATH="$PATH:/usr/local/go/bin"

RUN curl -O -s https://dl.google.com/go/$GOZIP && \
    tar -xzf $GOZIP -C /usr/local && rm $GOZIP && \
    mkdir -p "$APPBIN/"{darwin,linux,windows}/share && \
    mkdir -p "$APPDIR"
ADD . "$APPDIR"
VOLUME "$APPDIR"

# mingw64 & dependecies
RUN dnf -y install \
    mingw32-gcc mingw64-gcc \
    mingw32-cairo mingw32-cairo-static mingw64-cairo mingw64-cairo-static cairo cairo-devel \
    #mingw32-freetype mingw64-freetype freetype freetype-devel \
    mingw32-gtk3 mingw64-gtk3 gtk3 gtk3-devel \
    mingw32-glib2-static mingw64-glib2-static glib2-devel\
    #mingw32-pango mingw64-pango pango pango-devel \
    ; dnf clean all
#
# gtk3 (useful: https://fedoraproject.org/wiki/Packaging:MinGW)
WORKDIR "/go"

# -- CPU arch: 64 bytes
ENV MINGW="x86_64-w64-mingw32"
ENV MINGW_PATH="/usr/$MINGW/sys-root/mingw"

# -- linux
ENV CC=gcc \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR "$APPDIR/cmd/safe"
RUN go mod download
RUN GO111MODULE=on \
    go build -o "$APPBIN/$GOOS/safe.$GOOS.$GOARCH"

# -- windows
ENV CC="$MINGW-gcc" \
    CGO_ENABLED=1 \
    CGO_LDFLAGS_ALLOW="-Wl,-luuid" \
    GOOS=windows \
    GOARCH=amd64 \
    PKG_CONFIG_PATH="$MINGW_PATH/lib/pkgconfig"

# to copy only ddl outside, get the CONTAINER_ID and copy them:
# ~ docker ps -alq
# ~ docker cp <CONTAINER_ID>:<MINGW_PATH>/bin <OUTSIDE_PATH>
WORKDIR "$APPDIR/cmd/safe"
RUN go mod download
RUN cp -ra "$MINGW_PATH/bin/." "$APPBIN/$GOOS/" && \
    rsync -av "$MINGW_PATH/share/icons" "$APPBIN/windows/share" \
        --exclude cursors --exclude scalable --exclude scalable-up-to-32 && \
    GO111MODULE=on \
    go build -ldflags "-H windowsgui" -o "$APPBIN/$GOOS/safe.$GOOS.$GOARCH.exe"


# -- build all releases
WORKDIR "$APPBIN"
RUN echo "#!/bin/sh -e" > build.sh && \
    #echo "tar -zcvf $APPDIR/safe.linux.386.tar.gz linux" >> build.sh && \
    echo "tar -zcvf $APPDIR/safe.linux.amd64.tar.gz linux" >> build.sh && \
    #echo "tar -zcvf $APPDIR/safe.windows.386.tar.gz windows" >> build.sh && \
    echo "tar -zcvf $APPDIR/safe.windows.amd64.tar.gz windows" >> build.sh && \
    chmod +x build.sh

CMD ["./build.sh"]

# To build it:
# ~ docker build -t safe .
# And create Zip files:
# ~ docker run -v "$(pwd)":/go/src/github.com/rvflash/safe safe
# Or see details about configs:
# ~ docker run -ti safe bash