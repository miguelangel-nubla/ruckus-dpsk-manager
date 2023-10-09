# Ruckus DPSK Manager

Ruckus DPSK Manager is a command-line utility for managing Dynamic Pre-Shared Key (DPSK) authentication on Ruckus wireless controllers. It provides functionality to perform tasks such as creating DPSK users, retrieving DPSK passphrases, and backing up controller configurations.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Commands](#commands)
- [License](#license)

## Installation

### Download from GitHub Releases

You can also download pre-built binaries from the [GitHub Releases page](https://github.com/miguelangel-nubla/ruckus-dpsk-manager/releases). Follow these steps:

1. Visit the [GitHub Releases page](https://github.com/miguelangel-nubla/ruckus-dpsk-manager/releases).

2. Download the appropriate binary for your operating system and architecture and rename it to `ruckus-dpsk-manager`.

3. Place the downloaded binary in a directory included in your system's PATH to make it executable.

4. Verify the installation by running:

   ```bash
   ruckus-dpsk-manager
   ```

Please note that this method provides pre-built binaries, so you don't need Go installed on your system.

### Using `go install`

To install Ruckus DPSK Manager using `go install`, you need to have Go (Golang) installed on your system.

Open a terminal and run the following command to install the tool directly from the source code:

```bash
go install github.com/miguelangel-nubla/ruckus-dpsk-manager/cmd/ruckus-dpsk-manager@latest
```

## Usage

Ruckus DPSK Manager is a command-line tool that can be used with various commands and options. Here is the general usage:

```bash
ruckus-dpsk-manager [options] <command> [arguments]
```

Use the `-help` option with any command to see its specific usage instructions.

### Main options

- `-server`: Ruckus controller server location (default: https://unleashed.ruckuswireless.com).
- `-username`: Username for logging in to the Ruckus controller (default: dpsk).
- `-password`: Password for logging in to the Ruckus controller (required).
- `-cacert`: Path to a custom CA certificate.
- `-debug`: Enable debug output.
- `-help`: Print usage information.

## Commands

### `backup`

Backup the Ruckus controller configuration.

```bash
ruckus-dpsk-manager backup <output_filename>
```

- `<output_filename>`: The name of the file where the backup will be saved.

### `dpsk`

Manage DPSK entries.

#### `create`

Create a new DPSK and print the passphrase.

```bash
ruckus-dpsk-manager dpsk create <wlanID> <username>
```

- `<wlanID>`: The ID of the WLAN.
- `<username>`: The username of the DPSK user.

#### `modify`

Finds DSPK entries matching `[filter-flags]` and modifies its properties according to `[value-flags]`.

```bash
ruckus-dpsk-manager dpsk modify [filter-flags] set [value-flags]
```

The list of available `[filter-flags]` and `[value-flags]` is the same and represent the property keys of a DPSK entry, use `--help` to list the available flags and its valid values.

You can delete entries setting the expiration date to a value in the past.

#### `list`

Finds DSPK entries matching `[filter-flags]`.

```bash
ruckus-dpsk-manager dpsk list [filter-flags]
```

The list of available `[filter-flags]` represent the property keys of a DPSK entry, use `--help` to list the available flags and its valid values.

## License

This project is licensed under the Apache-2.0 license. See the [LICENSE](LICENSE) file for details.
