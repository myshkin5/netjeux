# Streaming Scheme

The Streaming scheme continuously streams messages at a specific rate.

## Configuration

 Dot path | Type | Required/Default | Description
 ---|---|---|---
 `streaming.messages-per-second` | `int` | No, `1000` | The count of message written to a Writer and read from a Reader per second. When set to zero (`0`), the reader or writer will read or write as quickly as possible.
 `streaming.expected-messages-per-second` | `int` | No, `0` | The count of messages **expected** to be written or read per second. When set to zero (`0`, the default), the value matches `streaming.messages-per-second`. Used when calculating message throughput percent.
 `streaming.bytes-per-message` | `int` | No, `1024` | The count of bytes per message.
 `streaming.report-cycle` | [`time.Duration`](https://golang.org/pkg/time/#ParseDuration) | No, `1s` (1 second) | The length of time between reports.

### Example JSON Configuration

```
{
    "additional": {
        "streaming": {
            "messages-per-second": 1000,
            "expected-messages-per-second": 0,
            "bytes-per-messages": 1024,
            "report-cycle": "1s"
        }
    }
}
```

### Example CLI

```
netspel ... \
    --config-int    .streaming.messages-per-second=1000 \
    --config-int    .streaming.expected-messages-per-second=0 \
    --config-int    .streaming.bytes-per-message=1024 \
    --config-string .streaming.report-cycle=1s
