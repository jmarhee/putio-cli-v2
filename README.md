# Put.io Command Line Client v2

![go-release](https://github.com/jmarhee/putio-cli-v2/actions/workflows/go-release.yml/badge.svg) ![go test](https://github.com/jmarhee/putio-cli-v2/actions/workflows/go.yml/badge.svg)

This package installs the v2 `putio` command-line client. This can be also used for any generic basic auth hosted file source and local filesystem target (as show in the go tests), but the use-case is optimized for connecting put.io zip archives to media server files.

Changes from `pyputio-cli`
---

`pyputio-cli` Python package has been [archived](https://github.com/jmarhee/pyputio-cli/tree/master). The following features do not (yet) exist in v2:

* `plex-scan` to initiate scans of the Plex server after unarchiving media (and the `PUTIO_PLEX_SCAN` configuration argument)
* `plex-clean` to clean up library (empty collections, etc.)
* `PUTIO_PLEX_METADATA_REFRESH_HARD` has been removed.
  
These will likely be a separate project to avoid overburdening the genralized benefits outside of the immediate use-case.

The following are unlikely to be implemented in `v2`:

* `PUTIO_MANUAL_DL` has been removed (v2 will not fallback to a manual client download).
* `PUTIO_REPORT_TIME` is now reported by default for each action, you will no longer receive the download/extract stats as a JSON blob (default for `PUTIO_OUTPUT_MODE`, which will not be implemented as well).

Setup
--

Clone this repo, and install:

```bash
go build -v .
```

or download the OS and architecture specific binary from the [releases page](https://github.com/jmarhee/putio-cli-v2/releases); i.e for Linux AMD64 systems:

```
sudo install putio-cli-linux-amd64 /usr/local/bin/putio
```
to install it to your `$PATH`.

Either set environment variables (**recommended**):

```bash
export PUTIO_USER=your_user
export PUTIO_PASS="your_password"
export PUTIO_LIBRARY_PATH=/mnt/Plex
export PUTIO_LIBRARY_SUBPATH=Movies #i.e. (the name of the subdir: Movies, TV, etc.)
```

for example, or set `PUTIO_CONFIG_PATH` to the configuration in the following format:

```toml
[putio_config]
username = "your_username"
password = "your_password"
library_path = "/mnt/Plex"
```
and loading this configuration:

```bash
putio --config config.ini --url ""
```

or, using the command line flags to set these options credentials:
`
```bash
putio --username "" --password "" --library_path "/mnt/target" --library_subpath "Music" --url ""
```

and set `PUTIO_LIBRARY_SUBPATH`, or reply when prompted. If required values are not set using the above, you will be prompted for them. 

You can set `PUTIO_CLEAN` to any value to have it clean up the zip archives after the download attempt.

Usage
---

Once installed, run:

```bash
putio "URL"
```
using the "Zip and Download" option on the Put.io UI. 

To remove archives after download and extraction, set `PUTIO_CLEAN` to 1. 

Environment Reference
---

The following environment variables can be set to assist usage as well. Details on use of options follow in other sections

| Option            | Description                                                                                      | Value                          |
|-------------------|--------------------------------------------------------------------------------------------------|--------------------------------|
| PUTIO_CLEAN       | Removes archive after extraction                                                                 | IfPresent                      |
| PUTIO_DIR_CREATE  | Creates subdirectories if they do not exist in library path                                      | IfPresent                      |
| PUTIO_NOTIFY      | Push notifications (Requires PUSHOVER_TOKEN and PUSHOVER_USER)                                   | IfPresent                      |

Integrations
---

The following are integrations available for this tool.

### Pushover

If you use [pushover](pushover.net), set `PUSHOVER_TOKEN` and `PUSHOVER_USER` in your environment, and to have completed jobs send a notification, set `--notify` to any value, or `PUTIO_NOTIFY` to any value. 

Contributing
---

Please feel free to tackle any open issues in a new PR, propose enhancements or desired features, etc. Tests will run on commit, but please test locally (a simple `go test` will suffice for the lack of complexity in this project) before opening a pull request. 
