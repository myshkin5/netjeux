# User Datagram Protocol

## Configuration

 Dot path | Type | Description
 ---|---|---
 `udp.remote-reader-addr` | `string` | The IP address of the remote reader process. Used by the writer process only.
 `udp.port` | `int` | The port on which the remote reader process listens. Used by the reader to setup a listener and used by the writer to write messages to.

### Example JSON Configuration

```
{
    "additional": {
        "udp": {
            "remote-reader-addr": 127.0.0.1,
            "port": 36644
        }
    }
}
```

### Example CLI

```
netspel ... \
    --config-string udp.remote-reader-addr=127.0.0.1 \
    --config-int    udp.port=36644
```
