<div align="center">
  <img src="resources/logo.png?raw=true">
</div>

---

Node implementation for the [Metal](https://metalblockchain.org) network. At its core Metal is a fork of the [Avalanche](https://avax.network) project and aims to maintain upstream changes.

## Installation

The minimum recommended hardware specification for nodes connected to Mainnet is:

- CPU: Equivalent of 8 AWS vCPU
- RAM: 16 GiB
- Storage: 250GiB
- OS: Ubuntu 20.04/22.04 or macOS >= 12
- Network: Reliable IPv4 or IPv6 network connection, with an open public port.

If you plan to build MetalGo from source, you will also need the following software:

- [Go](https://golang.org/doc/install) version >= 1.20.10
- [gcc](https://gcc.gnu.org/)
- g++

### Building From Source

#### Clone The Repository

Clone the MetalGo repository:

```sh
git clone git@github.com:MetalBlockchain/metalgo.git
cd metalgo
```

This will clone and checkout the `master` branch.

#### Building the Metal Executable

Build Metal using the build script:

```sh
./scripts/build.sh
```

The Metal binary, named `metalgo`, is in the `build` directory.

### Docker Install

Make sure Docker is installed on the machine - so commands like `docker run` etc. are available.

Building the docker image of latest metalgo branch can be done by running:

```sh
./scripts/build_image.sh
```

To check the built image, run:

```sh
docker image ls
```

The image should be tagged as `metalblockchain/metalgo:xxxxxxxx`, where `xxxxxxxx` is the shortened commit of the Metal source it was built from. To run the Metal node, run:

```sh
docker run -ti -p 9650:9650 -p 9651:9651 metalblockchain/metalgo:xxxxxxxx /metalgo/build/metalgo
```

## Running Metal

### Connecting to Mainnet

To connect to the Metal Mainnet, run:

```sh
./build/metalgo
```

You should see some pretty ASCII art and log messages.

You can use `Ctrl+C` to kill the node.

### Connecting to Tahoe

To connect to the Tahoe Testnet, run:

```sh
./build/metalgo --network-id=tahoe
```

## Supported Platforms

MetalGo can run on different platforms, with different support tiers:

- **Tier 1**: Fully supported by the maintainers, guaranteed to pass all tests including e2e and stress tests.
- **Tier 2**: Passes all unit and integration tests but not necessarily e2e tests.
- **Tier 3**: Builds but lightly tested (or not), considered _experimental_.
- **Not supported**: May not build and not tested, considered _unsafe_. To be supported in the future.

The following table lists currently supported platforms and their corresponding
MetalGo support tiers:

| Architecture | Operating system | Support tier  |
| :----------: | :--------------: | :-----------: |
|    amd64     |      Linux       |       1       |
|    arm64     |      Linux       |       2       |
|    amd64     |      Darwin      |       2       |
|    amd64     |     Windows      |       3       |
|     arm      |      Linux       | Not supported |
|     i386     |      Linux       | Not supported |
|    arm64     |      Darwin      | Not supported |

To officially support a new platform, one must satisfy the following requirements:

| MetalGo continuous integration     | Tier 1  | Tier 2  | Tier 3  |
| ---------------------------------- | :-----: | :-----: | :-----: |
| Build passes                       | &check; | &check; | &check; |
| Unit and integration tests pass    | &check; | &check; |         |
| End-to-end and stress tests pass   | &check; |         |         |

## Security Bugs

**We and our community welcome responsible disclosures.**

We're working on a bug bountry program, in the meantime please submit any bugs or vulnerabilities to bugs@metalblockchain.org
