# Streaming Scheme

The Streaming scheme continuously streams messages at a specific rate.

## Configuration

 Dot path | Type | Required/Default | Description
 ---|---|---|---
 `streaming.messages-per-second` | `int` | No, `1000` | The count of message written to a Writer and read from a Reader per second. When set to zero (`0`), the reader or writer will read or write as quickly as possible.
 `streaming.bytes-per-message` | `int` | No, `1024` | The count of bytes per message.

### Example JSON Configuration

```
{
    "additional": {
        "streaming": {
            "messages-per-second": 1000,
            "bytes-per-messages": 1024
        }
    }
}
```

### Example CLI

```
netspel ... \
    --config-int streaming.messages-per-second=1000 \
    --config-int streaming.bytes-per-message=1024
