# Message Format

## File Handshake

The file handshake payload is used to determine if the destination system will accept an inbound file transfer. The receiving system SHOULD check that it has available space, permission to write to the file path and that the signature matches and public key of the agent name.

```json
{
  "payload": {
    "id": "{uuid}",
    "type": "file_handshake",
    "agent": "{agent name}",
    "file_path": "{file path}",
    "file_size": 123456
  },
  "signature": "{signed payload}"
}
```

## File Accept

The file accept payload is used to tell the source agent that the destination agent is ready to accept the file. The `id` of the `file_handshake` message must be the same in this accept message, and the message must be placed onto the queue topic matching the `agent` within the payload.

The source agent should also verify that the receiving agent signature is valid.

```json
{
  "payload": {
    "id": "{uuid}",
    "type": "file_accept"
  },
  "signature": "{signed payload}"
}
```

## File Ready

Once a file is ready to be downloaded, the payload will include all available information used to retrieve the file in case the service has restarted in the duration of the file upload.

```json
{
  "payload": {
    "id": "{uuid}",
    "type": "file_ready",
    "signed_url": "{signed url}",
    "file_path": "{file path}",
    "file_size": 000000,
    "file_sha256": "{file sha256 checksum}"
  }
}
```

## File Reject

The file reject payload is used to tell the source system that the transfer will be rejected & can supply a reason.

```json
{
  "payload": {
    "id": "{uuid}",
    "type": "file_reject",
    "reason": "INSUFFICIENT_SPACE | INSUFFICIENT_PERMISSION | NOT_ALLOWED | PENDING_UPDATE | MAINTENNACE | OTHER_ERROR",
    "retry_in": 000 // optionally specify a retry_in if the reason is PENDING_UPDATE or MAINTENANCE
  }
}
```

## File Request

The file request payload is used to ask a source system if it can supply a file. If accepted, it will initiate the file handshake from its end using the same UUID of the request. The requesting agent will be aware of the

```json
{
  "payload": {
    "id": "{uuid}",
    "type": "file_request",
    "agent": "{requesting agent}",
    "file_path": "{file path}"
  }
}
```
