ARG VARIANT=24-bookworm
FROM mcr.microsoft.com/devcontainers/typescript-node:${VARIANT}

## Copy all the scripts
#
COPY scripts/*.sh /tmp/scripts/

## Install taskfile
#
RUN bash /tmp/scripts/install_taskfile.sh

## Change user to vscode
#
USER node