# Protocol commands

ID Ranges:

| ID Range | Scope                            |
|----------|----------------------------------|
| 0xxx     | System commands                  |
| 1xxx     | Database and collection commands |
| 2xxx     | Key-value engine commands        |
| 3xxx     | Document engine commands         |
| 4xxx     | Object engine commands           |
| 5xxx     | Analytical engine commands       |
| 6xxx     | Broker engine commands           |

<!-- TOC -->
* [Protocol commands](#protocol-commands)
  * [System Commands (0xxx)](#system-commands-0xxx)
    * [Ping (1)](#ping-1)
    * [Supported protocol versions (2)](#supported-protocol-versions-2)
    * [Login (100)](#login-100)
    * [Auth status (101)](#auth-status-101)
  * [Database and collection commands](#database-and-collection-commands)
  * [Database and collection commands (1xxx)](#database-and-collection-commands-1xxx)
  * [Key-value engine commands (2xxx)](#key-value-engine-commands-2xxx)
    * [Set (2000)](#set-2000)
    * [Set with TTL (2001)](#set-with-ttl-2001)
    * [Delete (2020)](#delete-2020)
    * [Delete all (2022)](#delete-all-2022)
    * [Exists (2030)](#exists-2030)
    * [Get (2031)](#get-2031)
    * [Get multiple (2032)](#get-multiple-2032)
    * [Get all (2033)](#get-all-2033)
    * [Keys (2034)](#keys-2034)
    * [Count (2035)](#count-2035)
  * [Document engine commands (3xxx)](#document-engine-commands-3xxx)
  * [Object engine commands (4xxx)](#object-engine-commands-4xxx)
    * [Put (4000)](#put-4000)
    * [Get (4001)](#get-4001)
    * [Get metadata (4002)](#get-metadata-4002)
    * [Delete (4003)](#delete-4003)
  * [Analytical engine commands (5xxx)](#analytical-engine-commands-5xxx)
  * [Broker engine commands (6xxx)](#broker-engine-commands-6xxx)
    * [Subscribe (6000)](#subscribe-6000)
    * [Subscribe queue (6001)](#subscribe-queue-6001)
    * [Unsubscribe (6002)](#unsubscribe-6002)
    * [Publish (6003)](#publish-6003)
    * [Message (client bound) (6004)](#message-client-bound-6004)
<!-- TOC -->

## System Commands (0xxx)

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

## Database and collection commands (1xxx)

## Key-value engine commands (2xxx)

### Set (2000)

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

### Set with TTL (2001)

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

### Delete (2020)

The KV Delete command deletes a key-value pair from the store.

Payload format:

| Field      | Size | Description         |
|------------|------|---------------------|
| Key length | 2 B  | Length of the key   |
| Key        | N B  | The key to delete   |

Response:

| Status code | Description   |
|-------------|---------------|
| 0000        | Success       |

### Delete multiple (2021)

The KV Delete multiple command deletes multiple key-value pairs from the store.

Payload format:

| Field | Size | Description                                   |
|-------|------|-----------------------------------------------|
| Keys  | N B  | [List of strings](protocol-encoded-values.md) |

Response:

| Status code | Description   |
|-------------|---------------|
| 0000        | Success       |

### Delete all (2022)

The KV Delete all command deletes all key-value pairs from the store.

Payload format: None

Response:

| Status code | Description   |
|-------------|---------------|
| 0000        | Success       |

### Exists (2030)

The KV Exists command checks if a key exists in the store.

Payload format:

| Field      | Size | Description         |
|------------|------|---------------------|
| Key length | 2 B  | Length of the key   |
| Key        | N B  | The key to check    |

Response:

| Status code | Description   |
|-------------|---------------|
| 0000        | Key exists    |
| 1008        | Key not found |

### Get (2031)

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

### Get multiple (2032)

The KV Get multiple command retrieves the values associated with multiple keys.

Payload format:

| Field | Size | Description                                   |
|-------|------|-----------------------------------------------|
| Keys  | N B  | [List of strings](protocol-encoded-values.md) |

Response:

| Status code | Description   |
|-------------|---------------|
| 0000        | Success       |

The response payload for a successful KV Get multiple command contains a map of keys to their corresponding values, where each value is encoded as an encoded value (see [Encoded Values](protocol-encoded-values.md)).
If a key does not exist, the value will be Empty.

### Get all (2033)

The KV Get all command retrieves all key-value pairs in the store.

Payload format: None

Response:

| Status code | Description   |
|-------------|---------------|
| 0000        | Success       |

The response payload for a successful KV Get all command contains a map of all keys to their corresponding values, where each value is encoded as an encoded value (see [Encoded Values](protocol-encoded-values.md)).

### Get TTL (2034)

The KV Get TTL command retrieves the uinx nano timestamp of when a key will expire.

Payload format:

| Field      | Size | Description         |
|------------|------|---------------------|
| Key length | 2 B  | Length of the key   |
| Key        | N B  | The key to check    |

Response:

| Status code | Description   |
|-------------|---------------|
| 0000        | Success       |
| 1008        | Key not found |

Response payload for a successful KV Get TTL command contains:

| Field      | Size | Description                                                                          |
|------------|------|--------------------------------------------------------------------------------------|
| Expires At | 8 B  | unix nano timestamp of when the key will expire, or 0 if the key does not have a TTL |

### Get multiple TTL (2035)

The KV Get multiple TTL command retrieves the uinx nano timestamp of when multiple keys will expire.

Payload format:

| Field | Size | Description                                   |
|-------|------|-----------------------------------------------|
| Keys  | N B  | [List of strings](protocol-encoded-values.md) | 

Response:

| Status code | Description   |
|-------------|---------------|
| 0000        | Success       |

Response payload for a successful KV Get multiple TTL command contains a map of keys to their corresponding expiration timestamps, where each timestamp is an 8-byte unix nano timestamp of when the key will expire, or 0 if the key does not have a TTL. If a key does not exist or is expired, the timestamp will be -1.

### Get all TTL (2036)

The KV Get all TTL command retrieves the uinx nano timestamp of when all keys will expire.

Payload format: None

Response:

| Status code | Description   |
|-------------|---------------|
| 0000        | Success       |

Response payload for a successful KV Get all TTL command contains a map of all keys to their corresponding expiration timestamps, where each timestamp is an 8-byte unix nano timestamp of when the key will expire, or 0 if the key does not have a TTL. If a key does not exist or is expired, the timestamp will be -1.

### Keys (2037)

The KV Keys command retrieves a list of all keys in the store.

Payload format: None

Response:

| Status code | Description   |
|-------------|---------------|
| 0000        | Success       |

The response payload for a successful KV Keys command contains the list of keys (strings).

### Count (2038)

The KV Count command retrieves the total number of key-value pairs in the store.

Payload format: None

Response:

| Status code | Description   |
|-------------|---------------|
| 0000        | Success       |

## Document engine commands (3xxx)

## Object engine commands (4xxx)

### Put (4000)

The Obj Put command stores an object in the engine.

Payload format:

| Field                                              | Size | Description         |
|----------------------------------------------------|------|---------------------|
| Key length                                         | 2 B  | Length of the key   |
| Key                                                | N B  | The key to set      |
| [Encoded binary value](protocol-encoded-values.md) | N B  | The object to store |

Response:

| Status code | Description   |
|-------------|---------------|
| 0000        | Success       |

### Get (4001)

The Obj Get command retrieves an object from the engine.

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

The response payload for a successful Obj Get command contains the object associated with the key, encoded as an encoded binary value (see [Encoded Values](protocol-encoded-values.md)).

### Get metadata (4002)

The Obj Get metadata command retrieves the metadata of an object from the engine.

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

Response payload:

| Field       | Size | Description                  |
|-------------|------|------------------------------|
| Size        | 8 B  | Size of the object in bytes  |
| Checksum    | 4 B  | CRC32 checksum of the object |
| Created at  | 8 B  | Unix millisecond timestamp   |
| Modified at | 8 B  | Unix millisecond timestamp   |

### Delete (4003)

The Obj Delete command deletes an object from the engine.

Payload format:

| Field      | Size | Description         |
|------------|------|---------------------|
| Key length | 2 B  | Length of the key   |
| Key        | N B  | The key to delete   |

Response:

| Status code | Description   |
|-------------|---------------|
| 0000        | Success       |
| 1008        | Key not found |

## Analytical engine commands (5xxx)

## Broker engine commands (6xxx)

### Subscribe (6000)

The Broker Subscribe command subscribes the client to a specific subject.

Payload format:

| Field          | Size | Description                 |
|----------------|------|-----------------------------|
| Subject length | 2 B  | Length of the subject       |
| Subject        | N B  | The subject to subscribe to |

Response:

| Status code | Description   |
|-------------|---------------|
| 0000        | Success       |

### Subscribe queue (6001)

The Broker Subscribe command subscribes the client to a specific subject.

Payload format:

| Field          | Size | Description                     |
|----------------|------|---------------------------------|
| Subject length | 2 B  | Length of the subject           |
| Subject        | N B  | The subject to subscribe to     |
| Queue length   | 2 B  | Length of the queue group       |
| Queue          | N B  | The queue group to subscribe to |

Response:

| Status code | Description   |
|-------------|---------------|
| 0000        | Success       |

### Unsubscribe (6002)

The Broker Unsubscribe command unsubscribes the client from a specific subject.

Payload format:

| Field          | Size | Description                     |
|----------------|------|---------------------------------|
| Subject length | 2 B  | Length of the subject           |
| Subject        | N B  | The subject to unsubscribe from |

### Publish (6003)

The Broker Publish command publishes a message to a specific subject.

Payload format:

| Field                                | Size | Description               |
|--------------------------------------|------|---------------------------|
| Subject length                       | 2 B  | Length of the subject     |
| Subject                              | N B  | The subject to publish to |
| [Binary](protocol-encoded-values.md) | N B  | The message to publish    |

### Message (client bound) (6004)

The Broker Message command is sent by the server to deliver a message to the client.

Payload format:

| Field                                        | Size | Description                |
|----------------------------------------------|------|----------------------------|
| Subject length                               | 2 B  | Length of the subject      |
| Subject                                      | N B  | The subject of the message |
| [List of Binary](protocol-encoded-values.md) | N B  | The messages to deliver    |

Response:

| Status code | Description   |
|-------------|---------------|
| 0000        | Success       |