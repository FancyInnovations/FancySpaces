# Encoded Values

Here you can see how values of common types are encoded in the protocol.

<!-- TOC -->
* [Encoded Values](#encoded-values)
  * [Bool (1)](#bool-1)
  * [Uint8 / Byte (2)](#uint8--byte-2)
  * [Uint16 (3)](#uint16-3)
  * [Uint32 (4)](#uint32-4)
  * [Uint64 (5)](#uint64-5)
  * [Int16 (6)](#int16-6)
  * [Int32 (7)](#int32-7)
  * [Int64 (8)](#int64-8)
  * [Float32 (9)](#float32-9)
  * [Float64 (10)](#float64-10)
  * [Binary (11)](#binary-11)
  * [String (12)](#string-12)
  * [List (13)](#list-13)
  * [Map (14)](#map-14)
<!-- TOC -->

## Bool (1)

| Field | Size | Description                      |
|-------|------|----------------------------------|
| Type  | 1 B  | Type of the value                |
| Value | 1 B  | Value of the boolean (0x01=true) |

## Uint8 / Byte (2)

| Field | Size | Description               |
|-------|------|---------------------------|
| Type  | 1 B  | Type of the value         |
| Value | 1 B  | Value of the uint8 / byte |

## Uint16 (3)

| Field | Size | Description                      |
|-------|------|----------------------------------|
| Type  | 1 B  | Type of the value                |
| Value | 2 B  | Value of the uint16              |

## Uint32 (4)

| Field | Size | Description         |
|-------|------|---------------------|
| Type  | 1 B  | Type of the value   |
| Value | 4 B  | Value of the uint32 |

## Uint64 (5)

| Field | Size | Description                      |
|-------|------|----------------------------------|
| Type  | 1 B  | Type of the value                |
| Value | 8 B  | Value of the uint64              |

## Int16 (6)

| Field | Size | Description                      |
|-------|------|----------------------------------|
| Type  | 1 B  | Type of the value                |
| Value | 2 B  | Value of the int16               |

## Int32 (7)

| Field | Size | Description        |
|-------|------|--------------------|
| Type  | 1 B  | Type of the value  |
| Value | 4 B  | Value of the int32 |

## Int64 (8)

| Field | Size | Description                      |
|-------|------|----------------------------------|
| Type  | 1 B  | Type of the value                |
| Value | 8 B  | Value of the int64               |

## Float32 (9)

| Field | Size | Description                      |
|-------|------|----------------------------------|
| Type  | 1 B  | Type of the value                |
| Value | 4 B  | Value of the float32             |

## Float64 (10)

| Field | Size | Description                      |
|-------|------|----------------------------------|
| Type  | 1 B  | Type of the value                |
| Value | 8 B  | Value of the float64             |

## Binary (11)

| Field  | Size | Description               |
|--------|------|---------------------------|
| Type   | 1 B  | Type of the value         |
| Length | 4 B  | Length of the binary data |
| Value  | N B  | Binary data               |


## String (12)

| Field  | Size | Description               |
|--------|------|---------------------------|
| Type   | 1 B  | Type of the value         |
| Length | 4 B  | Length of the string data |
| Value  | N B  | String data               |

## List (13)

List is encoded as an array of items, where each item can be of any type (including another list or map). 
The "Count" field specifies how many items are in the list, and the "Items" field contains the encoded items.

| Field          | Size | Description                    |
|----------------|------|--------------------------------|
| Type           | 1 B  | Type of the value              |
| Count          | 4 B  | Number of items in the array   |
| Items          | N B  | Encoded items in the array     |

Encoded item:

| Field       | Size | Description                                    |
|-------------|------|------------------------------------------------|
| Item length | 4 B  | Length of the encoded item                     |
| Item value  | N B  | Encoded value of the item (can be of any type) |


## Map (14)

Map is encoded as a collection of key-value pairs, where each key is a string and each value can be of any type (including another list or map). 
The "Count" field specifies how many key-value pairs are in the map, and the "Key Value Pairs" field contains the encoded key-value pairs.

This data type is used to represent JSON like objects in the client SDKs.

| Field           | Size | Description                          |
|-----------------|------|--------------------------------------|
| Type            | 1 B  | Type of the value                    |
| Count           | 2 B  | Number of key-value pairs in the map |
| Key-Value Pairs | N B  | Encoded key-value pairs in the map   |

Encoded key-value pair:

| Field        | Size | Description                                                  |
|--------------|------|--------------------------------------------------------------|
| Key length   | 2 B  | Length of the encoded key                                    |
| Key value    | N B  | Key as string (just the raw string bytes, not encoded again) |
| Value length | 4 B  | Length of the encoded value                                  |
| Value value  | N B  | Encoded value of the value (can be of any type)              |