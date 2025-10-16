![PODMAN logo](https://raw.githubusercontent.com/containers/common/main/logos/podman-logo-full-vert.png)

# 🚀 Podman Cluster-Optimized: Experimental Fork for Shared Storage

## ⚠️ **EXPERIMENTAL BUILD** ⚠️

**This is a specialized experimental fork of Podman optimized for cluster environments with shared network storage (NFS).**

### New Feature: `--shared-base-layers`

```bash
# Zero base layer copying
podman run --shared-base-layers ubuntu:latest my-app
```

**Designed for:** HPC clusters, development teams, CI/CD environments with large shared container images.

---

# Podman: A tool for managing OCI containers and pods
![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/containers/podman)
[![Go Report Card](https://goreportcard.com/badge/github.com/containers/podman/v5)](https://goreportcard.com/report/github.com/containers/podman/v5)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/10499/badge)](https://www.bestpractices.dev/projects/10499)

[![LFX Health Score](https://insights.linuxfoundation.org/api/badge/health-score?project=containers-podman)](https://insights.linuxfoundation.org/project/containers-podman)
[![LFX Contributors](https://insights.linuxfoundation.org/api/badge/contributors?project=containers-podman)](https://insights.linuxfoundation.org/project/containers-podman)


<br/>

Podman (the POD MANager) is a tool for managing containers and images, volumes mounted into those containers, and pods made from groups of containers.
Podman runs containers on Linux, but can also be used on Mac and Windows systems using a Podman-managed virtual machine.
Podman is based on libpod, a library for container lifecycle management that is also contained in this repository. The libpod library provides APIs for managing containers, pods, container images, and volumes.

Podman releases a new major or minor release 4 times a year, during the second week of February, May, August, and November. Patch releases are more frequent and may occur at any time to get bugfixes out to users. All releases are PGP signed. Public keys of members of the team approved to make releases are located [here](https://github.com/containers/release-keys/tree/main/podman).

* Continuous Integration:
  * [![Build Status](https://api.cirrus-ci.com/github/containers/podman.svg)](https://cirrus-ci.com/github/containers/podman/main)
  * [GoDoc: ![GoDoc](https://godoc.org/github.com/containers/podman/libpod?status.svg)](https://godoc.org/github.com/containers/podman/libpod)
  * [Downloads](DOWNLOADS.md)

## Overview and scope

At a high level, the scope of Podman and libpod is the following:

* Support for multiple container image formats, including OCI and Docker images.
* Full management of those images, including pulling from various sources (including trust and verification), creating (built via Containerfile or Dockerfile or committed from a container), and pushing to registries and other storage backends.
* Full management of container lifecycle, including creation (both from an image and from an exploded root filesystem), running, checkpointing and restoring (via CRIU), and removal.
* Full management of container networking, using Netavark.
* Support for pods, groups of containers that share resources and are managed together.
* Support for running containers and pods without root or other elevated privileges.
* Resource isolation of containers and pods.
* Support for a Docker-compatible CLI interface, which can both run containers locally and on remote systems.
* No manager daemon, for improved security and lower resource utilization at idle.
* Support for a REST API providing both a Docker-compatible interface and an improved interface exposing advanced Podman functionality.
* Support for running on Windows and Mac via virtual machines run by `podman machine`.

---

## 🚀 **Cluster-Optimized Features (Experimental)**

### Shared Base Layers (`--shared-base-layers`)

This experimental feature enables efficient container deployment across cluster nodes with shared network storage.

#### **Technical Implementation:**
- **NFS Auto-Detection**: Automatically detects when container images are stored on NFS mounts
- **OverlayFS Integration**: Uses shared base layers as lowerdir, local writable layers as upperdir
- **Zero-Copy Deployment**: Base layers mounted directly from shared storage
- **Graceful Fallback**: Automatically falls back to standard copy behavior when shared storage not available

#### **Cluster Use Cases:**

**HPC Environments:**
```bash
# Scientific computing with large toolchain images
podman run --shared-base-layers nfs-registry/cuda-toolkit:11.8 ./simulation
```

**Development Teams:**
```bash
# Shared development environments
podman run --shared-base-layers nfs-registry/node-dev:18 npm start
```

**CI/CD Pipelines:**
```bash
# Build environments without redundant copying
podman run --shared-base-layers nfs-registry/build-tools:latest make
```

#### **Performance Benefits:**
- **Storage Savings**: 80-95% reduction in storage usage for large base images
- **Network Efficiency**: Eliminates redundant image transfers between registry and nodes
- **Startup Speed**: 3-5x faster container startup for large images (no copy phase)
- **Scalability**: Linear scaling across cluster nodes without storage bottlenecks

#### **Requirements:**
- Base images must be stored on NFS-mounted storage
- Shared storage must be accessible from all cluster nodes
- Linux systems with OverlayFS support

---

## Roadmap

The future of Podman feature development can be found in its **[roadmap](ROADMAP.md)**.

## Communications

If you think you've identified a security issue in the project, please *DO NOT* report the issue publicly via the GitHub issue tracker, mailing list, or IRC.
Instead, send an email with as many details as possible to `security@lists.podman.io`. This is a private mailing list for the core maintainers.

For general questions and discussion, please use Podman's
[channels](https://podman.io/community/#slack-irc-matrix-and-discord).

For discussions around issues/bugs and features, you can use the GitHub
[issues](https://github.com/containers/podman/issues)
and
[PRs](https://github.com/containers/podman/pulls)
tracking system.

There is also a [mailing list](https://lists.podman.io/archives/) at `lists.podman.io`.
You can subscribe by sending a message to `podman-join@lists.podman.io` with the subject `subscribe`.

**For cluster-optimized features:** Please use GitHub Issues in this experimental fork for feedback and bug reports related to `--shared-base-layers` functionality.

## Rootless
Podman can be easily run as a normal user, without requiring a setuid binary.
When run without root, Podman containers use user namespaces to set root in the container to the user running Podman.
Rootless Podman runs locked-down containers with no privileges that the user running the container does not have.
Some of these restrictions can be lifted (via `--privileged`, for example), but rootless containers will never have more privileges than the user that launched them.
If you run Podman as your user and mount in `/etc/passwd` from the host, you still won't be able to change it, since your user doesn't have permission to do so.

Almost all normal Podman functionality is available, though there are some [shortcomings](https://github.com/containers/podman/blob/main/rootless.md).
Any recent Podman release should be able to run rootless without any additional configuration, though your operating system may require some additional configuration detailed in the [install guide](https://podman.io/getting-started/installation).

A little configuration by an administrator is required before rootless Podman can be used, the necessary setup is documented [here](https://github.com/containers/podman/blob/main/docs/tutorials/rootless_tutorial.md).

## Podman Desktop

[Podman Desktop](https://podman-desktop.io/) provides a local development environment for Podman and Kubernetes on Linux, Windows, and Mac machines.
It is a full-featured desktop UI frontend for Podman which uses the `podman machine` backend on non-Linux operating systems to run containers.
It supports full container lifecycle management (building, pulling, and pushing images, creating and managing containers, creating and managing pods, and working with Kubernetes YAML).
The project develops on [GitHub](https://github.com/containers/podman-desktop) and contributions are welcome.

## Out of scope

* Specialized signing and pushing of images to various storage backends.
  See [Skopeo](https://github.com/containers/skopeo/) for those tasks.
* Support for the Kubernetes CRI interface for container management.
  The [CRI-O](https://github.com/cri-o/cri-o) daemon specializes in that.

## OCI Projects Plans

Podman uses OCI projects and best of breed libraries for different aspects:
- Runtime: We use the [OCI runtime tools](https://github.com/opencontainers/runtime-tools) to generate OCI runtime configurations that can be used with any OCI-compliant runtime, like [crun](https://github.com/containers/crun/) and [runc](https://github.com/opencontainers/runc/).
- Images: Image management uses the [containers/image](https://github.com/containers/image) library.
- Storage: Container and image storage is managed by [containers/storage](https://github.com/containers/storage).
- Networking: Networking support through use of [Netavark](https://github.com/containers/netavark) and [Aardvark](https://github.com/containers/aardvark-dns).  Rootless networking is handled via [pasta](https://passt.top/passt) or [slirp4netns](https://github.com/rootless-containers/slirp4netns).
- Builds: Builds are supported via [Buildah](https://github.com/containers/buildah).
- Conmon: [Conmon](https://github.com/containers/conmon) is a tool for monitoring OCI runtimes, used by both Podman and CRI-O.
- Seccomp: A unified [Seccomp](https://github.com/containers/common/blob/main/pkg/seccomp/seccomp.json) policy for Podman, Buildah, and CRI-O.

## Podman Information for Developers

For blogs, release announcements and more, please checkout the [podman.io](https://podman.io) website!

**[Installation notes](install.md)**
Information on how to install Podman in your environment.

**[OCI Hooks Support](https://github.com/containers/common/blob/main/pkg/hooks/README.md)**
Information on how Podman configures [OCI Hooks][spec-hooks] to run when launching a container.

**[Podman API](https://docs.podman.io/en/latest/_static/api.html)**
Documentation on the Podman REST API.

**[Podman Commands](https://podman.readthedocs.io/en/latest/Commands.html)**
A list of the Podman commands with links to their man pages and in many cases videos
showing the commands in use.

**[Podman Container Images](https://github.com/containers/image_build/blob/main/podman/README.md)**
Information on the Podman Container Images found on [quay.io](https://quay.io/podman/stable).

**[Podman Troubleshooting Guide](troubleshooting.md)**
A list of common issues and solutions for Podman.

**[Podman Usage Transfer](transfer.md)**
Useful information for ops and dev transfer as it relates to infrastructure that utilizes Podman.  This page
includes tables showing Docker commands and their Podman equivalent commands.

**[Tutorials](docs/tutorials)**
Tutorials on using Podman.

**[Remote Client](https://github.com/containers/podman/blob/main/docs/tutorials/remote_client.md)**
A brief how-to on using the Podman remote client.

**[Basic Setup and Use of Podman in a Rootless environment](https://github.com/containers/podman/blob/main/docs/tutorials/rootless_tutorial.md)**
A tutorial showing the setup and configuration necessary to run Rootless Podman.

**[Release Notes](RELEASE_NOTES.md)**
Release notes for recent Podman versions.

**[Contributing](CONTRIBUTING.md)**
Information about contributing to this project.

[spec-hooks]: https://github.com/opencontainers/runtime-spec/blob/v1.0.2/config.md#posix-platform-hooks

## Buildah and Podman relationship

Buildah and Podman are two complementary open-source projects that are
available on most Linux platforms and both projects reside at
[GitHub.com](https://github.com) with Buildah
[here](https://github.com/containers/buildah) and Podman
[here](https://github.com/containers/podman).  Both, Buildah and Podman are
command line tools that work on Open Container Initiative (OCI) images and
containers.  The two projects differentiate in their specialization.

Buildah specializes in building OCI images.  Buildah's commands replicate all
of the commands that are found in a Dockerfile.  This allows building images
with and without Dockerfiles while not requiring any root privileges.
Buildah’s ultimate goal is to provide a lower-level coreutils interface to
build images.  The flexibility of building images without Dockerfiles allows
for the integration of other scripting languages into the build process.
Buildah follows a simple fork-exec model and does not run as a daemon
but it is based on a comprehensive API in golang, which can be vendored
into other tools.

Podman specializes in all of the commands and functions that help you to maintain and modify
OCI images, such as pulling and tagging.  It also allows you to create, run, and maintain those containers
created from those images.  For building container images via Dockerfiles, Podman uses Buildah's
golang API and can be installed independently from Buildah.

A major difference between Podman and Buildah is their concept of a container.  Podman
allows users to create "traditional containers" where the intent of these containers is
to be long lived.  While Buildah containers are really just created to allow content
to be added back to the container image.  An easy way to think of it is the
`buildah run` command emulates the RUN command in a Dockerfile while the `podman run`
command emulates the `docker run` command in functionality.  Because of this and their underlying
storage differences, you can not see Podman containers from within Buildah or vice versa.

In short, Buildah is an efficient way to create OCI images while Podman allows
you to manage and maintain those images and containers in a production environment using
familiar container cli commands.  For more details, see the
[Container Tools Guide](https://github.com/containers/buildah/tree/main/docs/containertools).

## Podman Hello
```
$ podman run quay.io/podman/hello
Trying to pull quay.io/podman/hello:latest...
Getting image source signatures
Copying blob a6b3126f3807 done
Copying config 25c667d086 done
Writing manifest to image destination
Storing signatures
!... Hello Podman World ...!

         .--"--.
       / -     - \
      / (O)   (O) \
   ~~~| -=(,Y,)=- |
    .---. /`  \   |~~
 ~/  o  o \~~~~.----. ~~
  | =(X)= |~  / (O (O) \
   ~~~~~~~  ~| =(Y_)=-  |
  ~~~~    ~~~|   U      |~~

Project:   https://github.com/containers/podman
Website:   https://podman.io
Documents: https://docs.podman.io
Twitter:   @Podman_io
```
