# Encoded Values

Here you can see how values of common types are encoded in the protocol.

## Value types

| Type ID | Type Name | Description              |
|---------|-----------|--------------------------|
| 0       | empty     | Empty / null / nil value |
| 1       | bool      | Boolean value            |
| 2       | byte      | Unsigned 8-bit integer   |
| 3       | uint16    | Unsigned 16-bit integer  |
| 4       | uint32    | Unsigned 32-bit integer  |
| 5       | uint64    | Unsigned 64-bit integer  |
| 6       | int16     | Signed 16-bit integer    |
| 7       | int32     | Signed 32-bit integer    |
| 8       | int64     | Signed 64-bit integer    |
| 9       | float32   | 32-bit floating point    |
| 10      | float64   | 64-bit floating point    |
| 11      | binary    | Binary data              |
| 12      | string    | String data              |
| 13      | array     | Array / list             |
| 14      | map       | Map / dictionary         |


## bool (1)

| Field | Size | Description                      |
|-------|------|----------------------------------|
| Type  | 1 B  | Type of the value                |
| Value | 1 B  | Value of the boolean (0x01=true) |

## uint8 / byte (2)

| Field | Size | Description               |
|-------|------|---------------------------|
| Type  | 1 B  | Type of the value         |
| Value | 1 B  | Value of the uint8 / byte |

## uint16 (3)

| Field | Size | Description                      |
|-------|------|----------------------------------|
| Type  | 1 B  | Type of the value                |
| Value | 2 B  | Value of the uint16              |

## uint32 (4)

| Field | Size | Description         |
|-------|------|---------------------|
| Type  | 1 B  | Type of the value   |
| Value | 4 B  | Value of the uint32 |

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

## int32 (7)

| Field | Size | Description        |
|-------|------|--------------------|
| Type  | 1 B  | Type of the value  |
| Value | 4 B  | Value of the int32 |

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

Example layout list of strings:

```
| 13 (type) | 12 (item type) | 0x0003 (count) | 0x0000000F (payload length) |
| 12 (type) | 0x00000005 (length) | 'H' 'e' 'l' 'l' 'o' (string 1)          |
| 12 (type) | 0x00000005 (length) | 'W' 'o' 'r' 'l' 'd' (string 2)          |
| 12 (type) | 0x00000006 (length) | 'F' 'o' 'o' 'B' 'a' 'r' (string 3)      |
```

## map (14)

| Field           | Size | Description                          |
|-----------------|------|--------------------------------------|
| Type            | 1 B  | Type of the value                    |
| Value Type      | 1 B  | Type of the values in the map        |
| Count           | 2 B  | Number of key-value pairs in the map |
| Payload length  | 4 B  | Length of the payload                |
| Key-Value Pairs | N B  | Encoded key-value pairs in the map   |

Note: Keys in the map are always strings, so there is no need to specify their type. All values are of the same type, which is specified in the "Value Type" field.

Example layout of a map of string to int32:

```
| 14 (type) | 7 (value type) | 0x0002 (count) | 0x0000001A (payload length) |
| 12 (type) | 0x00000003 (length) | 'a' 'g' 'e' (key 1) | 7 (type) | 0x0000001E (value 1) |
| 12 (type) | 0x00000004 (length) | 'n' 'a' 'm' 'e' (key 2) | 7 (type) | 0x00000004 (value 2) |
```