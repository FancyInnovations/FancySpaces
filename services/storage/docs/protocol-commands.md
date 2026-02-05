# Protocol commands

ID Ranges:

| ID Range | Scope                            |
|----------|----------------------------------|
| 0xxx     | System commands                  |
| 1xxx     | Database and collection commands |
| 2xxx     | Key-value engine commands        |

## System Commands

### Ping (1)

The Ping command is used to check the connectivity and responsiveness of the server.

Payload format: None

Response payload: 

| Field         | Size   | Description            |
|---------------|--------|------------------------|
| Pong          | 4 B    | "pong" string          |


### Supported protocol versions (2)

Returns the list of supported protocol versions by the server.

Payload format: None

Response payload:

| Field    | Size | Description                                |
|----------|------|--------------------------------------------|
| Legnth   | 1 B  | Number of versions                         |
| Versions | N B  | List of supported versions (one byte each) |


### Login (100)

The Login command is used to authenticate a user with the server.

Payload format:

| Field | Size | Description         |
|-------|------|---------------------|
| Type  | 1 B  | Authentication type |
| Data  | N B  | Authentication data |

If `Type` is 0x01 (Password), the `Data` field contains:

| Field           | Size | Description        |
|-----------------|------|--------------------|
| Username length | 2 B  | Length of username |
| Username        | N B  | User's username    |
| Password length | 2 B  | Length of password |
| Password        | N B  | User's password    |

If `Type` is 0x02 (API key), the `Data` field contains:

| Field          | Size | Description            |
|----------------|------|------------------------|
| API key length | 2 B  | Length of api key      |
| API key        | N B  | Authentication api key |

Response: 

| Status code | Description                |
|-------------|----------------------------|
| 0000        | Successfully authenticated |
| 1002        | Invalid credentials        |

### Auth status (101)

The Auth status command checks the current authentication status of the user.

Payload format: None

Response payload:

Response:

| Status code | Description   |
|-------------|---------------|
| 0000        | Authenticated |
| 1004        | Unauthorized  |

## Database and collection commands

Not implemented yet.

## Key-value engine commands

### KV Get (2000)

The KV Get command retrieves the value associated with a given key.

Payload format:

| Field      | Size | Description         |
|------------|------|---------------------|
| Key length | 2 B  | Length of the key   |
| Key        | N B  | The key to retrieve |

Response:

| Status code | Description   |
|-------------|---------------|
| 0000        | Success       |
| 1008        | Key not found |

The response payload for a successful KV Get command contains the value associated with the key, encoded as an encoded value (see [Encoded Values](protocol-encoded-values.md)).

### KV Set (2001)

The KV Set command sets the value for a given key.

Payload format:

| Field                                       | Size | Description       |
|---------------------------------------------|------|-------------------|
| Key length                                  | 2 B  | Length of the key |
| Key                                         | N B  | The key to set    |
| [Encoded Value](protocol-encoded-values.md) | N B  | The value to set  |

Response:

| Status code | Description   |
|-------------|---------------|
| 0000        | Success       |

### KV Set with TTL (2002)

The KV Set with TTL command sets the value for a given key with a specified time-to-live (TTL).

Payload format:

| Field                                       | Size | Description         |
|---------------------------------------------|------|---------------------|
| Key length                                  | 2 B  | Length of the key   |
| Key                                         | N B  | The key to set      |
| [Encoded Value](protocol-encoded-values.md) | N B  | The value to set    |
| Expires At                                  | 8 B  | unix nano timestmap |


Response:

| Status code | Description   |
|-------------|---------------|
| 0000        | Success       |
