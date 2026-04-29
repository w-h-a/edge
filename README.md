# edge

Embeddable local-first database with CRDT sync.

Go library with a thin host daemon (edged). The library provides Open, Put, Get, Watch, Close. The daemon embeds the library, watches files, polls bees, and syncs to delta in the background. edged runs as a launchd service: starts on login, runs forever, restarts on crash. One-time setup, then invisible.

Reads and writes are local SQLite operations. Sync happens in the background when a delta node is reachable. Offline is the default mode, not a degradation.

## Public API

```go
store, err := edge.Open("/path/to/data.db")

store.Put(ctx, "user/prefs/theme", []byte("dark"), edge.LWW)
val, err := store.Get(ctx, "user/prefs/theme")

events, err := store.Watch(ctx, edge.Scope{Project: "tally"})

store.Close()
```

That is the API. If it grows beyond Open, Put, Get, Watch, Close, the abstraction has failed.

## Data Flow

```mermaid
graph LR
    subgraph "Local (always available)"
        APP[Application] -->|Put/Get| STORE[Store]
        STORE --> SQL[SQLite<br/>entries + version vectors]
    end

    subgraph "Background (when available)"
        STORE -->|periodic| SM[SyncManager]
        SM -->|connect| DELTA[delta node]
        SM -->|exchange version vectors| DELTA
        SM -->|send/receive deltas| DELTA
        SM -->|CRDT merge| STORE
    end

    subgraph "Producers"
        BEES[bees<br/>task history] -->|poll| IM[IngestManager]
        FILES[~/.claude/<br/>CLAUDE.md, memory] -->|watch| IM
        IM -->|Put with merge strategy| STORE
    end

    style SQL fill:#0f3460,stroke:#e94560,color:#eee
    style STORE fill:#0f3460,stroke:#e94560,color:#eee
```

## Architecture

```mermaid
graph TD
    subgraph "cmd/edged"
        MAIN[main.go<br/>launchd service]
    end

    subgraph "Public API"
        EDGE[edge.go<br/>Open, Close]
        STORE[store.go<br/>Put, Get, Watch]
    end

    subgraph "internal/domain/"
        ENT[entry.go<br/>context record + CRDT metadata]
        CLK[clock.go<br/>version vectors]
        DLT[delta.go<br/>DeltaCompute pure function]
        SES[session.go<br/>sync state machine]
        SCP[scope.go<br/>namespace model]
    end

    subgraph "internal/service/"
        SYN[sync.go<br/>background sync lifecycle]
        ING[ingest.go<br/>context capture from producers]
    end

    subgraph "internal/client/ (ports)"
        PER[persister/<br/>local SQLite]
        TRA[transport/<br/>gRPC to delta]
    end

    subgraph "meld (library dependency)"
        CRDT[crdt/<br/>LWW-Register, OR-Set]
        VC[crdt/vclock]
    end

    MAIN --> EDGE
    EDGE --> STORE
    STORE --> SYN
    STORE --> ING
    SYN --> DLT
    SYN --> SES
    SYN --> TRA
    SYN --> PER
    ING --> ENT
    ING --> PER
    ENT --> CRDT
    CLK --> VC
```

## Sync Session

```mermaid
sequenceDiagram
    participant E as edged (device)
    participant D as delta (homelab)

    Note over E: Background timer fires
    E->>D: Connect
    E->>D: Send my version vector
    D->>E: Send your version vector
    Note over E,D: Both compute deltas
    E->>D: Entries you need (DeltaCompute)
    D->>E: Entries I need (DeltaCompute)
    Note over E: CRDT merge received entries
    Note over D: CRDT merge received entries
    E->>D: Confirm convergence
    Note over E: Back to idle
```

If sync is interrupted mid-delta, version vectors track what was confirmed. Next sync resumes from last confirmed vector.

## Merge Strategies

| Context type | CRDT         | Why                                                 |
| ------------ | ------------ | --------------------------------------------------- |
| Preferences  | LWW-Register | Most recent setting wins                            |
| Corrections  | LWW-Register | Latest correction supersedes                        |
| Patterns     | OR-Set       | Patterns accumulate, never lost                     |
| Task history | Append-only  | Events are immutable facts                          |
| Files (V1)   | LWW per file | Simplest. Concurrent file edits rare across devices |

## Seven Ideals

| Ideal                  | How edge satisfies it                                                        |
| ---------------------- | ---------------------------------------------------------------------------- |
| No spinners            | Reads and writes are local SQLite. Instant.                                  |
| Multi-device           | CRDT sync via delta relay.                                                   |
| Network optional       | Full functionality offline. Sync when available.                             |
| Seamless collaboration | CRDTs resolve concurrent edits automatically.                                |
| The Long Now           | Data outlives any server or service. No server required to access your data. |
| Privacy by default     | Data lives on the user's device.                                             |
| User ownership         | User controls and owns their data.                                           |

## Dependencies

- **meld**: CRDT types, version vectors, delta-state computation. Must be complete.
- **delta**: sync relay endpoint. Must have sync handler.

Observability via Telemetry port (OTel).
