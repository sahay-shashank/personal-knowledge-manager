ARG VARIANT=2-1.25-trixie
FROM mcr.microsoft.com/devcontainers/go:${VARIANT}

## Copy all the scripts
#
COPY scripts/*.sh /tmp/scripts/

## Install taskfile
#
RUN bash /tmp/scripts/install_taskfile.sh
