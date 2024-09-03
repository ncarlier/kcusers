# Simple Keycloak user management CLI

Simple CLI used to do some user management tasks on a Keycloak instance.

## Installation

Run the following command:

```bash
$ go install github.com/ncarlier/kcusers@latest
```

**Or** download the binary regarding your architecture:

```bash
$ curl -sf https://gobinaries.com/ncarlier/kcusers | sh
```

## Build

```bash
make
```

## Geting started

Create the configuration file:

```bash
./kcusers init-config -f config.toml
```

Customize the configuration file:

```toml
[log]
level = "info"
format = "text"

[keycloak]
authority = "http://loclahost:8080"
realm = "test"
client_id = "xxx"
client_secret = "yyy"
timeout = "5s"
```

Play with the CLI

```bash
# Get CLI usage
./kcusers -h
# Count users
./kcusers count-users
# Get user
./kcusers get-user -uid ffcc46cc-f66d-4df8-a623-a6d54ff242df
# Delete users
./kcusers delete-users -f users_to_delete.txt --dry-run --concurent 5
# Count client's sessions
./kcusers count-sessions -cid 2ca2c534-59e9-4039-94fb-562072cd1c11
```

## License

The MIT License (MIT)

See [LICENSE](./LICENSE) to see the full text.

---