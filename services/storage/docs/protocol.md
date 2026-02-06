# Storage TCP Protocol

Transport Layer: TCP

Port: 8091

Endianness: Big-Endian

## Message Structure

### Frame Format

| Field           | Size | Description                        |
|-----------------|------|------------------------------------|
| Length          | 4 B  | Total message length               |
| Message Payload | N B  | [Message](#message-format) content |

### Message Format

| Field            | Size | Description                    |
|------------------|------|--------------------------------|
| Magic Number     | 1 B  | Protocol magic number (0x7E)   |
| Protocol Version | 1 B  | Version of the protocol (0x01) |
| Flags            | 1 B  | Message flags                  |
| Type             | 1 B  | Message type identifier        |
| Payload Length   | 4 B  | Length of the command string   |
| Payload          | N B  | Payload data                   |

Message Types:

| Type | Description                 | Payload format               |
|------|-----------------------------|------------------------------|
| 0x01 | Command (client <-> server) | [Command](#command-format)   |
| 0x02 | Response (server -> client) | [Response](#response-format) |


### Command Format

| Field                  | Size | Description                            |
|------------------------|------|----------------------------------------|
| Request ID             | 2 B  | Request identifier (unique per client) |
| Command ID             | 2 B  | Unique command identifier              |
| Database Name Length   | 2 B  | Length of the database name            |
| Database Name          | N B  | Name of the target database            |
| Collection Name Length | 2 B  | Length of the collection name          |
| Collection Name        | N B  | Name of the target collection          |
| Payload Length         | 4 B  | Length of the command payload          |
| Payload                | N B  | Command-specific data                  |

See [protocol-commands.md](protocol-commands.md) for a list of supported commands and their payload formats.

### Response Format

| Field          | Size | Description                      |
|----------------|------|----------------------------------|
| Request ID     | 2 B  | Matches the command's request ID |
| Status Code    | 2 B  | Response status code             |
| Payload Length | 4 B  | Length of the response           |
| Payload        | N B  | Response data                    |

Status Codes:

| Code | Description             |
|------|-------------------------|
| 0xxx | Success codes           |
| 1xxx | Client-side error codes |
| 2xxx | Server-side error codes |

See `services/storage/internal/protocol/statuscodes.go` for a complete list of status codes.

If the status is 1xxx or 2xxx, the payload is an error message string.
Otherwise, the payload contains data relevant to the command executed (see command documentation for details).