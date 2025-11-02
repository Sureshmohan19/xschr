# Use a specific, stable version of Ubuntu as the base
FROM ubuntu:22.04

# Avoid interactive prompts during package installation
ENV DEBIAN_FRONTEND=noninteractive

# Install all necessary development tools in a single layer
RUN apt-get update && apt-get install -y \
    build-essential \
    gdb \
    git \
    man-db \
    vim \
    && rm -rf /var/lib/apt/lists/*

# Set the default working directory for when the container starts
WORKDIR /xschr

