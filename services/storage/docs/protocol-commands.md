# Protocol commands

ID Ranges:

| ID Range | Scope                            |
|----------|----------------------------------|
| 0xxx     | System commands                  |
| 1xxx     | Database and collection commands |
| 2xxx     | Key-value engine commands        |

## System Commands

### Ping (0001)

The Ping command is used to check the connectivity and responsiveness of the server.

Payload format: None

Response payload: 

| Field         | Size   | Description            |
|---------------|--------|------------------------|
| Pong          | 4 B    | "pong" string          |


### Supported protocol versions (0002)

Returns the list of supported protocol versions by the server.

Payload format: None

Response payload:

| Field    | Size | Description                                |
|----------|------|--------------------------------------------|
| Legnth   | 1 B  | Number of versions                         |
| Versions | N B  | List of supported versions (one byte each) |
