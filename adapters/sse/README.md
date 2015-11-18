# Server-Sent Events

## Configuration

 Dot path | Type | Description
 ---|---|---
 `sse.remote-writer-addr` | `string` | The IP address of the remote writer process. Used by the reader process only.
 `sse.port` | `int` | The port on which the remote writer process listens. Used by the writer to setup a listener and used by the reader to read messages from.

### Example JSON Configuration

```
{
    "additional": {
        "sse": {
            "remote-writer-addr": 127.0.0.1,
            "port": 36644
        }
    }
}
```

### Example CLI

```
netspel ... \
    --config-string sse.remote-writer-addr=127.0.0.1 \
    --config-int    sse.port=36644
```
