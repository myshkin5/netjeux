# Simple Scheme

The Simple scheme is the simplest scheme available. The Simple scheme is focused on maximizing single reader/writer throughput.

## Configuration

 Dot path | Type | Description
 ---|---|---
 `simple.messages-per-run` | `int` | The count of message sent to a Writer and expected from a Reader.
 `simple.bytes-per-message` | `int` | The count of bytes per message.
 `simple.wait-for-last-message` | [`time.Duration`](https://golang.org/pkg/time/#ParseDuration) | The time to wait after the last message is read before a run is considered complete. Used only in `read` mode. Default `5s` (5 seconds).
 `simple.warmup-messages-per-run` | `int` | The count of messages used to "warmup" the network channel. A non-zero value is required for some protocols to have accurate timings. For instance pull protocols must send warm up messages so that the `write` mode doesn't start the run before the Reader is ready to read messages.
 `simple.warmup-wait` | [`time.Duration`](https://golang.org/pkg/time/#ParseDuration) | The time to wait after warmup messages are sent before sending actual messages. Default `5s` (5 seconds).

### Example JSON Configuration

```
{
    "additional": {
        "simple": {
            "messages-per-run": 10000,
            "bytes-per-messages": 1000,
            "wait-for-last-message": "10s",
            "warmup-messages-per-run": 5,
            "warmup-wait": "2s"
        }
    }
}
```

### Example CLI

```
netspel ... \
    --config-int simple.messages-per-run=10000 \
    --config-int simple.bytes-per-message=10 \
    --config-string simple.wait-for-last-message=10s \
    --config-int simple.warmup-messages-per-run=5 \
    --config-string simple.warmup-wait=2s
