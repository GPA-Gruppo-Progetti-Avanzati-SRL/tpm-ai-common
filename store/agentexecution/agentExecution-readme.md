### Schema agentexecution

<!-- @tpm-schematics:start-region("top-file-section") -->
<!-- @tpm-schematics:end-region("top-file-section") -->

<!-- @tpm-schematics:start-region("bottom-file-section") -->

## sample JSON

Sample JSON document conforming to the `AgentExecution` struct.

```json
{
  "_id": "665f1c2a9b3e4a1d2c8f7e01",
  "domain": "tpm",
  "site": "site-01",
  "_bid": "agent-execution-bid-001",
  "_et": "agent-execution",
  "status": "ready",
  "batch_id": "batch-2026-08-25-001",
  "weight": 10,
  "bid_ref": {
    "bid": "prompt-queue-item-bid-042",
    "et": "prompt-queue-item"
  },
  "params": {
    "model": "claude-opus-4-8",
    "temperature": "0.2",
    "tags": ["research", "draft"]
  },
  "group": "group-a",
  "count": 3
}
```

## Field reference

| JSON field | Type      | Description                                            |
|------------|-----------|--------------------------------------------------------|
| `_id`      | ObjectID  | MongoDB object id.                                     |
| `domain`   | string    | Domain the execution belongs to.                       |
| `site`     | string    | Site identifier.                                       |
| `_bid`     | string    | Business id of the document.                           |
| `_et`      | string    | Entity type.                                           |
| `status`   | string    | One of `ready`, `staged`, `workings`, `done`, `error`. |
| `batch_id` | string    | Batch grouping identifier.                             |
| `weight`   | int32     | Execution weight.                                      |
| `bid_ref`  | BidEtPair | Reference to another entity (`bid` + `et`).            |
| `params`   | object    | Free-form parameters map.                              |
| `group`    | string    | Group identifier.                                      |
| `count`    | int32     | Counter value.                                         |
```

<!-- @tpm-schematics:end-region("bottom-file-section") -->
