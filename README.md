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

- [hardware-observe](https://snapcraft.io/docs/hardware-observe-interface)
- [home](https://snapcraft.io/docs/home-interface) - only if on Ubuntu Core
- `etc-default-grub` plug into the [system-files](https://snapcraft.io/docs/system-files-interface) interface;
- `proc-device-tree-model` plug into the [system-files](https://snapcraft.io/docs/system-files-interface) interface;
- `proc-irq` plug into the [system-files](https://snapcraft.io/docs/system-files-interface) interface;
- `sys-kernel-irq` plug into the [system-files](https://snapcraft.io/docs/system-files-interface) interface;

These can be done by running the following commands:

```shell
sudo snap connect rt-conf:hardware-observe
sudo snap connect rt-conf:home # Only in case of Ubuntu Core
sudo snap connect rt-conf:etc-default-grub
sudo snap connect rt-conf:proc-device-tree-model
sudo snap connect rt-conf:proc-irq
sudo snap connect rt-conf:sys-kernel-irq
```

Copy the example configuration file to a working directory accessible to the snap.
For example, copy it to the home directory:
```shell
cp /snap/rt-conf/current/config.yaml ~/rt-conf.yaml
```

## Usage

Run rt-conf to apply the configurations:

```shell
sudo rt-conf --file=/home/ubuntu/rt-conf.yaml
```

Set `--help` for more details.

The rt-conf app can be set to as a oneshot service on system startup.
This is useful for re-applying unpersisted IRQ tuning and power management settings on boot.

To do so, set the configuration file path as snap configuration:
```shell
sudo snap set rt-conf config-file=/home/ubuntu/rt-conf.yaml
```

Then, start and enable the service:
```
sudo snap start --enable rt-conf
```

### Debug logging

To enable debug logging, set either:
- `DEBUG=1` environment variable or
- `debug=1` snap configuration option.


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
