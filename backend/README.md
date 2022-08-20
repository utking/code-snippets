# Code Snippets Backend

To build the code, run

```bash
cd backend
go get
go build # append -ldflags "-w -s" to strip the debug symbols
```

## Install

To install it, run `install.sh`. All the required files will be copied over to the right places.

Before running the application, check and adjust `/etc/code-snippets/env`. The session token was auto-generated and can be left as-is. The SQLite DB file path can also be changed.

To start the application, run

```bash
systemctl start code-snippets
```

## Uninstall

To uninstall - run `uninstall.sh`, which will remove all the related files.
