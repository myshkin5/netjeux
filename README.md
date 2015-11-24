# netspel - A network protocol playground

netspel is a playground to analyze various network protocols in situ -- running in a variety of environments. It can be used to understand protocol characteristics and estimate theoretical maximums.

## Configuration

The configuration glues everything together. Preferably the same configuration is used by both readers and writer but there may need to be minor differences. A JSON file is the base for all configuration; but as long as all of the required fields are present, the file is optional. The file is of the form:

```
{
    "scheme-type": "<scheme-type>",
    "writer-type": "<writer-type>",
    "reader-type": "<reader-type>",
    "additional": {
        ...
    }
}
```
 Field | CLI option | Description
 ---|---|---
 `scheme-type` | `--scheme` or `-s` | A type implementing the [Scheme interface](factory/scheme.go)
 `writer-type` | `--writer` or `-w` | A type implementing the [Writer interface](factory/adapter.go#L9)
 `reader-type` | `--reader` or `-r` | A type implementing the [Reader interface](factory/adapter.go#L14)
 `additional` | `--config-string` or `--config-int` | An optional section specifying arbitrary data used by the specified types. See below for CLI override mechanism.

Options specified on the command line take precedence over JSON values. Values in the additional section can be specified or override JSON values using the following format:

`--config-string <dot path>=<value>` or `--config-int <dot path>=<value>`

Dot paths are specified relative to the additional section. For example, given the following JSON config file:

```
{
    "additional": {
        "udp": {
            "port": 23456
        }
    }
}
```

The port value can be overridden to a value of `12345` using the CLI option `--config-int udp.port=12345`.

## Schemes

Schemes orchestrate a run without any coupling to a specific protocol. Schemes can exercise readers and writers all while measuring various attributes of the run.

### Existing Schemes

 Type | Description
 ---|---
 [`simple`](schemes/simple) | The simplest scheme available.
 [`streaming`](schemes/streaming) | The Streaming scheme continuously streams messages at a specific rate.

## Adapters

Adapters allow schemes to read and write using a specific network protocol.

### Existing Adapter Writers

 Type | Protocol
 ---|---
 [`udp`](adapters/udp) | [User Datagram Protocol](https://en.wikipedia.org/wiki/User_Datagram_Protocol)
 [`sse`](adapters/sse) | [Server-Sent Events](https://en.wikipedia.org/wiki/Server-sent_events)

### Existing Adapter Readers

 Type | Protocol
 ---|---
 [`udp`](adapters/udp) | [User Datagram Protocol](https://en.wikipedia.org/wiki/User_Datagram_Protocol)
 [`sse`](adapters/sse) | [Server-Sent Events](https://en.wikipedia.org/wiki/Server-sent_events)
