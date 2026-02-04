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

| Field            | Size | Description                        |
|------------------|------|------------------------------------|
| Magic Number     | 1 B  | Protocol magic number (0x7E)       |
| Protocol Version | 1 B  | Version of the protocol (0x01)     |
| Flags            | 1 B  | Message flags                      |
| Type             | 1 B  | Message type identifier            |
| Command Length   | 4 B  | Length of the command string       |
| Command Payload  | N B  | [Command](#command-format) content |

### Command Format

| Field                  | Size | Description                   |
|------------------------|------|-------------------------------|
| Command ID             | 2 B  | Unique command identifier     |
| Database Name Length   | 2 B  | Length of the database name   |
| Database Name          | N B  | Name of the target database   |
| Collection Name Length | 2 B  | Length of the collection name |
| Collection Name        | N B  | Name of the target collection |
| Payload Length         | 4 B  | Length of the command payload |
| Payload                | N B  | Command-specific data         |