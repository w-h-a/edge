# edge

Embeddable local-first database with CRDT sync.

Go library. Import edge and embed in your application. Reads and writes are local SQLite operations. Sync happens in the background when a delta node is reachable. Offline is the default mode, not a degradation.

## Public API

```go
store, err := edge.Open("/path/to/data.db")

store.Put(ctx, "user/prefs/theme", []byte("dark"), edge.LWW)
val, err := store.Get(ctx, "user/prefs/theme")

events, err := store.Watch(ctx, edge.Scope{Project: "tally"})

store.Close()
```

That is the API. If it grows beyond Open, Put, Get, Watch, Close, the abstraction has failed.

## Architecture

```mermaid
graph TD
    subgraph "Public API"
        EDGE[edge.go<br/>Open, Close]
        STORE[store.go<br/>Put, Get, Watch]
    end

    subgraph "internal/domain/"
        ENT[entry.go<br/>context record + CRDT metadata]
        VER[vv.go<br/>version vectors]
        DLT[delta.go<br/>DeltaCompute pure function]
        SES[session.go<br/>sync state machine]
        SCP[scope.go<br/>namespace model]
    end

    subgraph "internal/service/"
        SYN[sync.go<br/>background sync lifecycle]
    end

    subgraph "internal/client/ (ports)"
        PER[persister/<br/>local SQLite]
        TRA[transport/<br/>http to delta]
    end

    subgraph "meld (library dependency)"
        CRDT[crdt/<br/>LWW-Register, OR-Set, Version Vectors]
    end

    EDGE --> STORE
    STORE --> SYN
    SYN --> DLT
    SYN --> SES
    SYN --> TRA
    SYN --> PER
    ENT --> CRDT
    VER --> CRDT
```

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
```

## Sync Session

```mermaid
sequenceDiagram
    participant E as edge (device)
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

| Data type     | CRDT         | Why                             |
| ------------- | ------------ | ------------------------------- |
| Settings      | LWW-Register | Most recent value wins          |
| Collections   | OR-Set       | Elements accumulate, never lost |
| Event history | Append-only  | Events are immutable facts      |

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
