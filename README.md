# portwatch

A lightweight CLI daemon that monitors and logs port activity changes on a host in real time.

---

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git
cd portwatch && go build -o portwatch .
```

---

## Usage

Start the daemon with default settings:

```bash
portwatch start
```

Watch specific ports and log output to a file:

```bash
portwatch start --ports 80,443,8080 --log /var/log/portwatch.log
```

Run with a custom polling interval (in seconds):

```bash
portwatch start --interval 5
```

Example output:

```
2024/01/15 10:23:01 [OPEN]   port 8080 — PID 3421 (node)
2024/01/15 10:23:11 [CLOSED] port 8080
2024/01/15 10:24:05 [OPEN]   port 3000 — PID 3789 (python3)
```

---

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--ports` | all | Comma-separated list of ports to watch |
| `--interval` | `2` | Polling interval in seconds |
| `--log` | stdout | Path to log file |

---

## License

MIT © 2024 [yourusername](https://github.com/yourusername)