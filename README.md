# kcusers (Keycloak User Management CLI)

A simple and fast CLI to perform user management tasks on a Keycloak instance.

## Installation

Install using Go:

```bash
go install github.com/ncarlier/kcusers@latest
```

Or install the binary directly (Linux/macOS):

```bash
curl -sf https://gobinaries.com/ncarlier/kcusers | sh
```

## Getting Started

1. **Initialize Configuration:**
   Generate a default configuration file:
   ```bash
   kcusers init-config -f kcusers.toml
   ```

2. **Configure:**
   Update `kcusers.toml` with your Keycloak credentials:
   ```toml
   [log]
   level = "info"
   format = "text"

   [keycloak]
   authority = "http://localhost:8080"
   realm = "test"
   client_id = "xxx"
   client_secret = "yyy"
   cache = ".kcusers-token.json"
   timeout = "5s"
   ```

3. **Usage Examples:**
   ```bash
   # Display help and available commands
   kcusers -h

   # Count total users
   kcusers count-users

   # Get a specific user by ID
   kcusers get-user -uid ffcc46cc-f66d-4df8-a623-a6d54ff242df

   # Delete multiple users (with dry-run and concurrency)
   kcusers delete-users -f users_to_delete.txt --dry-run --concurrent 5

   # Count active sessions for a client
   kcusers count-sessions -cid 2ca2c534-59e9-4039-94fb-562072cd1c11
   ```

## Development

Build the project using Make:

```bash
make build
```

## License

MIT License. See [LICENSE](./LICENSE) for details.
