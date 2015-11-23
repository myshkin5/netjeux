# Streaming Scheme

The Streaming scheme continuously streams messages at a specific rate.

## Configuration

 Dot path | Type | Description
 ---|---|---
 `streaming.messages-per-second` | `int` | The count of message sent to a Writer and expected from a Reader per second.
 `streaming.bytes-per-message` | `int` | The count of bytes per message.

### Example JSON Configuration

```
{
    "additional": {
        "streaming": {
            "messages-per-second": 10000,
            "bytes-per-messages": 1000
        }
    }
}
```

### Example CLI

```
netspel ... \
    --config-int streaming.messages-per-second=10000 \
    --config-int streaming.bytes-per-message=10
