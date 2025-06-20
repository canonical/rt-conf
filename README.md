[![rt-conf](https://snapcraft.io/rt-conf/badge.svg)](https://snapcraft.io/rt-conf)

# rt-conf - Real-Time Configuration Tool

This is a tool to help with configuration and tuning of real-time Ubuntu.

## Installation

To install the rt-conf snap:

```shell
sudo snap install rt-conf
```

This snap is confined, which means that it can access a limited set of resources on the system.
Additional access is granted via [snap interfaces](https://snapcraft.io/docs/interfaces).

After the installation it's necessary to connect the interfaces:

- [cpu-control](https://snapcraft.io/docs/cpu-control-interface)
- `etc-default-grub` plug into the [system-files](https://snapcraft.io/docs/system-files-interface) interface;
- [hardware-observe](https://snapcraft.io/docs/hardware-observe-interface)
- [home](https://snapcraft.io/docs/home-interface) - auto connected on classic distributions

These can be done by running the following commands:

```shell
sudo snap connect rt-conf:cpu-control
sudo snap connect rt-conf:etc-default-grub
sudo snap connect rt-conf:hardware-observe
sudo snap connect rt-conf:home
```

### Default configuration file

Upon installation the default rt-conf configuration file is added at: `/var/snap/rt-conf/common/config.yaml`.

## Usage

Edit the [default configuration file](#default-configuration-file) or a copy of it.
In case of a copy, it must be placed it in a directory accessible to the snap, such as the user home directory.
The copy must be owned by and writable to the root user only.

Run rt-conf to apply the configurations:

```shell
sudo rt-conf --file=/var/snap/rt-conf/common/config.yaml
```

Set `--help` for more details.

The rt-conf app can be set to run as a oneshot service on system startup.
This is useful for re-applying non-persistent IRQ tuning and power management settings on boot.

By default, the service reads the [default configuration file](#default-configuration-file).

To change the config file path, use the `config-file` snap configuration. Example:

```shell
sudo snap set rt-conf config-file=/home/ubuntu/rt-conf.yaml
```

Then, start and enable the service:

```shell
sudo snap start --enable rt-conf
```

Verify that it ran successfully by looking into the logs:

```shell
sudo snap logs -n 100 rt-conf
```

### Verbose logging

To enable verbose logging, set:

- `--verbose` flag on the CLI
- `verbose=true` snap configuration option for the service

## Hacking

Firstly, clone the repository:

```shell
git clone https://github.com/canonical/rt-conf.git
```

It's possible to run the `rt-conf` application from source by having Go installed and running:

```shell
go run cmd/rt-conf/main.go
```

> [!TIP]
> For local hacking on GRUB systems, it's recommended to use the local grub file included at `test/grub`.
> Also, you may want to use the local `config.yaml` file provided on the root of the repository:
>
> ```shell
> go run cmd/rt-conf/main.go --file=./config.yaml -ui --grub-file=./test/grub
> ```

Run tests:

```shell
go test ./...
```

### Local Build

Firstly, build it using [Snapcraft](https://snapcraft.io/snapcraft):

```shell
snapcraft -v
```

Then, install it in [dangerous mode](https://snapcraft.io/docs/install-modes#heading--dangerous):

```shell
sudo snap install --dangerous *.snap
```
