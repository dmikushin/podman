####> This option file is used in:
####>   podman create, run
####> If file is edited, make sure the changes
####> are applicable to all of those.
#### **--shared-base-layers**

Skip copying base layers and use them directly from shared storage.

When used with container images stored on shared storage (such as NFS), this option
mounts base layers directly without copying them to local storage. This reduces
storage usage and improves container startup time, especially for large images.

Podman automatically detects when base layers are available on shared storage
and enables the optimization. If shared storage is not detected or base layers
are not available on shared storage, Podman falls back to the standard behavior
of copying layers to local storage.

**Requirements:**
- Base layers must be stored on shared storage (NFS is automatically detected)
- The shared storage must be accessible from the host system

**Example:**

    $ podman <<subcommand>> --shared-base-layers ubuntu:latest echo "Hello World"

**Note:** This option only affects base layers; writable layers are always created
in local storage regardless of this setting.