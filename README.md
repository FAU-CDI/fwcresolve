# wdresolve

This repository implements the resolver for the [WissKI Distillery](https://github.com/FAU-CDI/wisski-distillery).
Documentation tbd.


# Deployment

[![Publish Docker Image](https://github.com/FAU-CDI/wdresolve/actions/workflows/docker.yml/badge.svg)](https://github.com/FAU-CDI/wdresolve/actions/workflows/docker.yml)

Available as a Docker Image on [GitHub Packages](https://github.com/FAU-CDI/wdresolve/pkgs/container/wdresolve).
Automatically built on every commit.

```bash
 docker run -ti -p 8080:8080 -e DEFAULT_DOMAIN=wisski.data.fau.de -e LEGACY_DOMAIN=wisski.agfd.fau.de -e PREFIX_FILE=/prefixes -v $(pwd)/cmd/wdresolve/prefixes.example:/prefixes:ro ghcr.io/fau-cdi/wdresolve
```
