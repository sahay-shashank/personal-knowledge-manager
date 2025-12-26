ARG VARIANT=3.16
FROM mcr.microsoft.com/devcontainers/go:${VARIANT}

## Copy all the scripts
#
COPY scripts/*.sh /tmp/scripts/

## Install taskfile
#
RUN bash /tmp/scripts/install_taskfile.sh