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

| Field | Size     | Description         |
|-------|----------|---------------------|
| Type  | 1 B      | Authentication type |
| Data  | Variable | Authentication data |

If `Type` is 0x01 (Password), the `Data` field contains:

| Field           | Size     | Description        |
|-----------------|----------|--------------------|
| Username length | 2 B      | Length of username |
| Username        | Variable | User's username    |
| Password length | 2 B      | Length of password |
| Password        | Variable | User's password    |

If `Type` is 0x02 (API key), the `Data` field contains:

| Field          | Size     | Description            |
|----------------|----------|------------------------|
| API key length | 2 B      | Length of api key      |
| API key        | Variable | Authentication api key |

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