### ListTransaction RPC - Sequence Diagram

```mermaid
sequenceDiagram
	autonumber
	participant RPC as ListTransaction RPC
	participant UC as ListTransaction UC
	participant BTCR as BTCRepo

	RPC->>+UC: Call
	UC->>+BTCR: Call `ListTransaction`
	BTCR-->>-UC: return
	UC-->>-RPC: return
```

