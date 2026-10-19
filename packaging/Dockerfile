# Registry index and linux/amd64 manifest resolved 2026-10-18; see environment-lock.json.
FROM debian@sha256:918311b7b6c4c6f68b232ba516584925f6c78ad82b6fd534b98979df6438e483
RUN rm -f /etc/apt/sources.list.d/debian.sources && \
    printf '%s\n' 'deb [check-valid-until=no] http://snapshot.debian.org/archive/debian/20261005T000000Z/ trixie main' > /etc/apt/sources.list && \
    apt-get update && apt-get install -y --no-install-recommends \
      dpkg=1.22.22 dpkg-dev=1.22.22 debhelper=13.24.2 python3 systemd adduser binutils lintian && \
    dpkg-query -W -f='${binary:Package}\t${Version}\n' > /packaging-packages.lock && \
    rm -rf /var/lib/apt/lists/*
WORKDIR /work
