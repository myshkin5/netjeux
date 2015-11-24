# Server-Sent Events

## Configuration

 Dot path | Type | Required/Default | Description
 ---|---|---|---
 `sse.port` | `int` | No, `38208` | The port on which the remote writer process listens. Used by the writer to setup a listener and used by the reader to read messages from.
 `sse.remote-writer-addr` | `string` | No, `localhost` | The IP address of the remote writer process. Used by the reader process only.

### Example JSON Configuration

```
{
    "additional": {
        "sse": {
            "port": 38208,
            "remote-writer-addr": 127.0.0.1
        }
    }
}
```

### Example CLI

```
netspel ... \
    --config-int    sse.port=38208 \
    --config-string sse.remote-writer-addr=127.0.0.1
```
