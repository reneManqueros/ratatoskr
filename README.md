# Ratatosk
## _DNS cache + upstreamer_

## Features

- Locally block DNS queries
- Add custom entries
- Fast!



## Installation

```sh
go build .
```

## Configuration

Edit config.json:

| Setting | Description |
| ------ | ------ |
| content_root | base path for the local blocklist/resolves folder |
| dns_port | Listen port (must be 53 for production usage |
| dns_address | Listen address |
| upstream_servers | Servers to use for upstreaming |
| upstream_timeout | Milliseconds before triggering an upstream timeout  |