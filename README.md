# Ruckus DPSK Manager

Ruckus DPSK Manager is a command-line utility for managing Dynamic Pre-Shared Key (DPSK) authentication on Ruckus wireless controllers. It provides functionality to perform tasks such as creating DPSK users, retrieving DPSK passphrases, and backing up controller configurations.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Commands](#commands)
- [License](#license)

## Installation

To install Ruckus DPSK Manager, follow these steps:

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/your-username/ruckus-dpsk-manager.git
   cd ruckus-dpsk-manager
   ```

2. **Build the Executable:**

   ```bash
   go build
   ```

3. **Run the Application:**

   ```bash
   ./ruckus-dpsk-manager
   ```

## Usage

Ruckus DPSK Manager is a command-line tool that can be used with various commands and options. Here is the general usage:

```bash
./ruckus-dpsk-manager [options] <command> [arguments]
```

Use the `-help` option with any command to see its specific usage instructions.

## Commands

### `backup`

Backup the Ruckus controller configuration.

```bash
./ruckus-dpsk-manager backup <output_filename>
```

- `<output_filename>`: The name of the file where the backup will be saved.

### `dpsk`

Manage DPSK users.

```bash
./ruckus-dpsk-manager dpsk <wlanID> <username>
```

- `<wlanID>`: The ID of the WLAN for which to manage DPSK users.
- `<username>`: The username of the DPSK user.

### Options

- `-server`: Ruckus controller server location (default: https://unleashed.ruckuswireless.com).
- `-username`: Username for logging in to the Ruckus controller (default: dpsk).
- `-password`: Password for logging in to the Ruckus controller (required).
- `-debug`: Enable debug output.
- `-help`: Print usage information.

## License

This project is licensed under the Apache-2.0 license. See the [LICENSE](LICENSE) file for details.
