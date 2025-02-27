---
title: Configuration
description: Understanding configuring.
---

Settings can be specified in three ways (in order of precedence):

1. Command line arguments.
1. Environment variables starting with **`SYNCRA_`**
1. **`syncra.yml`** config file

## Config file location

Config file will be loaded from the following paths:

- `/etc/syncra`
- `$HOME/.syncra`
- `./config`

### Config file example

```yaml
# Syncra example configuration file
# server: false
# bootstrap-expect: 3
# data-dir: syncra.data
# log-level: debug
# tags:
#   dc: eu
# encrypt: a-valid-key-generated-with-syncra-keygen
# retry-join:
#   - 10.0.0.1
#   - 10.0.0.2
#   - 10.0.0.3
# raft-multiplier: 1
```
Yes, it's that simple

### SEE ALSO

* [syncra agent](/en/cli)	 - Start a syncra agent
* [syncra doc](/en/cli/syncra_doc)	 - Generate Markdown documentation for the Syncra CLI.
* [syncra keygen](/en/cli/syncra_keygen)	 - Generates a new encryption key
* [syncra version](/en/cli/syncra_version)	 - Show version