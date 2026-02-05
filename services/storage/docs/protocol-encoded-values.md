# Encoded Values

Here you can see how values of common types are encoded in the protocol.


## bool (1)

| Field | Size | Description                      |
|-------|------|----------------------------------|
| Type  | 1 B  | Type of the value                |
| Value | 1 B  | Value of the boolean (0x01=true) |

## byte / uint8 (2)

| Field | Size | Description                      |
|-------|------|----------------------------------|
| Type  | 1 B  | Type of the value                |
| Value | 1 B  | Value of the byte / uint8        |

## uint16 (3)

| Field | Size | Description                      |
|-------|------|----------------------------------|
| Type  | 1 B  | Type of the value                |
| Value | 2 B  | Value of the uint16              |

## uint / uint32 (4)

| Field | Size | Description                      |
|-------|------|----------------------------------|
| Type  | 1 B  | Type of the value                |
| Value | 4 B  | Value of the uint / uint32       |

## uint64 (5)

| Field | Size | Description                      |
|-------|------|----------------------------------|
| Type  | 1 B  | Type of the value                |
| Value | 8 B  | Value of the uint64              |

## int16 (6)

| Field | Size | Description                      |
|-------|------|----------------------------------|
| Type  | 1 B  | Type of the value                |
| Value | 2 B  | Value of the int16               |

## int / int32 (7)

| Field | Size | Description                      |
|-------|------|----------------------------------|
| Type  | 1 B  | Type of the value                |
| Value | 4 B  | Value of the int / int32         |

## int64 (8)

| Field | Size | Description                      |
|-------|------|----------------------------------|
| Type  | 1 B  | Type of the value                |
| Value | 8 B  | Value of the int64               |

## float32 (9)

| Field | Size | Description                      |
|-------|------|----------------------------------|
| Type  | 1 B  | Type of the value                |
| Value | 4 B  | Value of the float32             |

## float64 (10)

| Field | Size | Description                      |
|-------|------|----------------------------------|
| Type  | 1 B  | Type of the value                |
| Value | 8 B  | Value of the float64             |

## binary (11)

| Field  | Size | Description               |
|--------|------|---------------------------|
| Type   | 1 B  | Type of the value         |
| Length | 4 B  | Length of the binary data |
| Value  | N B  | Binary data               |


## string (12)

| Field  | Size | Description               |
|--------|------|---------------------------|
| Type   | 1 B  | Type of the value         |
| Length | 4 B  | Length of the string data |
| Value  | N B  | String data               |

## array / list (13)

| Field          | Size | Description                    |
|----------------|------|--------------------------------|
| Type           | 1 B  | Type of the value              |
| Item Type      | 1 B  | Type of the items in the array |
| Count          | 2 B  | Number of items in the array   |
| Payload length | 4 B  | Length of the payload          |
| Items          | N B  | Encoded items in the array     |

(Note: every item in the array has the same type)

## map (14)

| Field           | Size | Description                          |
|-----------------|------|--------------------------------------|
| Type            | 1 B  | Type of the value                    |
| Value Type      | 1 B  | Type of the values in the map        |
| Count           | 2 B  | Number of key-value pairs in the map |
| Payload length  | 4 B  | Length of the payload                |
| Key-Value Pairs | N B  | Encoded key-value pairs in the map   |

(Note: Keys in the map are always strings, so there is no need to specify their type)