# Protocol commands

ID Ranges:

| ID Range | Scope                            |
|----------|----------------------------------|
| 0xxx     | System commands                  |
| 1xxx     | Database and collection commands |
| 2xxx     | Key-value engine commands        |

<!-- TOC -->
* [Protocol commands](#protocol-commands)
  * [System Commands](#system-commands)
    * [Ping (1)](#ping-1)
    * [Supported protocol versions (2)](#supported-protocol-versions-2)
    * [Login (100)](#login-100)
    * [Auth status (101)](#auth-status-101)
  * [Database and collection commands](#database-and-collection-commands)
  * [Key-value engine commands](#key-value-engine-commands)
    * [KV Set (2000)](#kv-set-2000)
    * [KV Set with TTL (2001)](#kv-set-with-ttl-2001)
    * [KV Set multiple](#kv-set-multiple)
    * [KV Set multiple with TTL](#kv-set-multiple-with-ttl)
    * [KV Set if exists](#kv-set-if-exists)
    * [KV Set if exists with TTL](#kv-set-if-exists-with-ttl)
    * [KV Set if not exists](#kv-set-if-not-exists)
    * [KV Set if not exists with TTL](#kv-set-if-not-exists-with-ttl)
    * [KV Delete (2020)](#kv-delete-2020)
    * [KV Delete multiple](#kv-delete-multiple)
    * [KV Delete all](#kv-delete-all)
    * [KV Exists (2030)](#kv-exists-2030)
    * [KV Get (2031)](#kv-get-2031)
    * [KV Get multiple (2032)](#kv-get-multiple-2032)
    * [KV Get all (2033)](#kv-get-all-2033)
    * [KV Keys (2034)](#kv-keys-2034)
    * [KV Number Increment](#kv-number-increment)
    * [KV Number Decrement](#kv-number-decrement)
    * [KV Number Multiply](#kv-number-multiply)
    * [KV Number Divide](#kv-number-divide)
    * [KV Number Modulo](#kv-number-modulo)
    * [KV Number Left shift](#kv-number-left-shift)
    * [KV Number Right shift](#kv-number-right-shift)
    * [KV Number Bitwise AND](#kv-number-bitwise-and)
    * [KV Number Bitwise OR](#kv-number-bitwise-or)
    * [KV Number Bitwise XOR](#kv-number-bitwise-xor)
    * [KV Number Bitwise NOT](#kv-number-bitwise-not)
    * [KV String Append](#kv-string-append)
    * [KV String Prepend](#kv-string-prepend)
    * [KV String Length](#kv-string-length)
    * [KV String Substring](#kv-string-substring)
    * [KV List length](#kv-list-length)
    * [KV List Get](#kv-list-get)
    * [KV List Set](#kv-list-set)
    * [KV List Remove](#kv-list-remove)
    * [KV List Push left](#kv-list-push-left)
    * [KV List Push right](#kv-list-push-right)
    * [KV List Pop left](#kv-list-pop-left)
    * [KV List Pop right](#kv-list-pop-right)
    * [KV Map length](#kv-map-length)
    * [KV Map Get](#kv-map-get)
    * [KV Map Set](#kv-map-set)
    * [KV Map Delete](#kv-map-delete)
    * [KV Map Exists](#kv-map-exists)
    * [KV Map Keys](#kv-map-keys)
    * [KV Map Values](#kv-map-values)
<!-- TOC -->

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

### KV Set (2000)

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

### KV Set with TTL (2001)

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

### KV Set multiple

Not implemented yet.

### KV Set multiple with TTL

Not implemented yet.

### KV Set if exists

Not implemented yet.

### KV Set if exists with TTL

Not implemented yet.

### KV Set if not exists

Not implemented yet.

### KV Set if not exists with TTL

Not implemented yet.

### KV Delete (2020)

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

### KV Delete multiple

Not implemented yet.

### KV Delete all

Not implemented yet.

### KV Exists (2030)

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

### KV Get (2031)

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

### KV Get multiple (2032)

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

### KV Get all (2033)

The KV Get all command retrieves all key-value pairs in the store.

Payload format: None

Response:

| Status code | Description   |
|-------------|---------------|
| 0000        | Success       |

The response payload for a successful KV Get all command contains a map of all keys to their corresponding values, where each value is encoded as an encoded value (see [Encoded Values](protocol-encoded-values.md)).

### KV Keys (2034)

The KV Keys command retrieves a list of all keys in the store.

Payload format: None

Response:

| Status code | Description   |
|-------------|---------------|
| 0000        | Success       |

The response payload for a successful KV Keys command contains the list of keys (strings).

### KV Number Increment

Not implemented yet.

### KV Number Decrement

Not implemented yet.

### KV Number Multiply

Not implemented yet.

### KV Number Divide

Not implemented yet.

### KV Number Modulo

Not implemented yet.

### KV Number Left shift

Not implemented yet.

### KV Number Right shift

Not implemented yet.

### KV Number Bitwise AND

Not implemented yet.

### KV Number Bitwise OR

Not implemented yet.

### KV Number Bitwise XOR

Not implemented yet.

### KV Number Bitwise NOT

Not implemented yet.

### KV String Append

Not implemented yet.

### KV String Prepend

Not implemented yet.

### KV String Length

Not implemented yet.

### KV String Substring

Not implemented yet.

### KV List length

Not implemented yet.

### KV List Get

Not implemented yet.

### KV List Set

Not implemented yet.

### KV List Remove

Not implemented yet.

### KV List Push left

Not implemented yet.

### KV List Push right

Not implemented yet.

### KV List Pop left

Not implemented yet.

### KV List Pop right

Not implemented yet.

### KV Map length

Not implemented yet.

### KV Map Get

Not implemented yet.

### KV Map Set

Not implemented yet.

### KV Map Delete

Not implemented yet.

### KV Map Exists

Not implemented yet.

### KV Map Keys

Not implemented yet.

### KV Map Values

Not implemented yet.
