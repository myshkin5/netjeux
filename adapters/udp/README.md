# User Datagram Protocol

## Configuration

 Dot path | Type | Required/Default | Description
 ---|---|---|---
 `udp.port` | `int` | No, `57955` | The port on which the remote reader process listens. Used by the reader to setup a listener and used by the writer to write messages to.
 `udp.remote-reader-addr` | `string` | No, `localhost` | The IP address of the remote reader process. Used by the writer process only.

### Example JSON Configuration

```
{
    "additional": {
        "udp": {
            "port": 57955,
            "remote-reader-addr": 127.0.0.1
        }
    }
}
```

### Example CLI

```
netspel ... \
    --config-int    .udp.port=57955 \
    --config-string .udp.remote-reader-addr=127.0.0.1
```
