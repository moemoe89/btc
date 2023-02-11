### CreateTransaction RPC - Sequence Diagram

```mermaid
sequenceDiagram
	autonumber
	participant RPC as CreateTransaction RPC
	participant UC as CreateTransaction UC
	participant BTCR as BTCRepo

	RPC->>+UC: Call
	UC->>+BTCR: Call `CreateTransaction`
	BTCR-->>-UC: return
	UC-->>-RPC: return
```

