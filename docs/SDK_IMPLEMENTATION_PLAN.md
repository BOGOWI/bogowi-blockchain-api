# NFT Tickets SDK Implementation Plan

## Overview
Build a comprehensive Go SDK for interacting with the BOGOWITickets smart contract on Camino Network.

## Architecture

```
internal/
├── sdk/
│   ├── nft/
│   │   ├── client.go           # Main SDK client
│   │   ├── tickets.go          # Ticket-specific methods
│   │   ├── signatures.go       # EIP-712 signature handling
│   │   ├── events.go           # Event listening and filtering
│   │   ├── queries.go          # Read-only query methods
│   │   └── types.go            # Type definitions
│   └── contracts/
│       ├── BOGOWITickets.go    # Generated bindings
│       └── RoleManager.go      # Generated bindings
```

## Phase 1: Contract Bindings Generation

### 1.1 Generate Go Bindings
```bash
# Install abigen if not present
go install github.com/ethereum/go-ethereum/cmd/abigen@latest

# Generate BOGOWITickets bindings
abigen --abi=contracts/v1/artifacts/contracts/nft/tickets/BOGOWITickets.sol/BOGOWITickets.json \
       --pkg=contracts \
       --type=BOGOWITickets \
       --out=internal/sdk/contracts/BOGOWITickets.go

# Generate RoleManager bindings
abigen --abi=contracts/v1/artifacts/contracts/utils/RoleManager.sol/RoleManager.json \
       --pkg=contracts \
       --type=RoleManager \
       --out=internal/sdk/contracts/RoleManager.go
```

### 1.2 Type Definitions
```go
// internal/sdk/nft/types.go
type MintParams struct {
    To                address.Address
    BookingID         [32]byte
    EventID           [32]byte
    UtilityFlags      uint32
    TransferUnlockAt  uint64
    ExpiresAt         uint64
    MetadataURI       string
    RewardBasisPoints uint16
    DatakyteNftID     string // For linking
}

type TicketData struct {
    BookingID                  [32]byte
    EventID                    [32]byte
    TransferUnlockAt           uint64
    ExpiresAt                  uint64
    UtilityFlags               uint32
    State                      uint8
    NonTransferableAfterRedeem bool
    BurnOnRedeem              bool
}

type RedemptionData struct {
    TokenID   *big.Int
    Redeemer  common.Address
    Nonce     *big.Int
    Deadline  *big.Int
    ChainID   *big.Int
    Signature []byte
}
```

## Phase 2: Core SDK Client

### 2.1 Client Structure
```go
// internal/sdk/nft/client.go
type Client struct {
    ethClient         *ethclient.Client
    ticketsContract   *contracts.BOGOWITickets
    roleManager       *contracts.RoleManager
    signer           *bind.TransactOpts
    chainID          *big.Int
    networkName      string
    datakyteService  *datakyte.TicketMetadataService
}

type ClientConfig struct {
    RPCURL           string
    PrivateKey       string
    TicketsAddress   string
    RoleManagerAddr  string
    DatakyteAPIKey   string
    NetworkName      string // "testnet" or "mainnet"
}

func NewClient(config ClientConfig) (*Client, error) {
    // Connect to Camino
    // Load contracts
    // Setup signer
    // Initialize Datakyte
}
```

## Phase 3: Minting Methods

### 3.1 Single Mint
```go
// internal/sdk/nft/tickets.go
func (c *Client) MintTicket(params MintParams) (*types.Transaction, uint64, error) {
    // 1. Validate parameters
    // 2. Convert to contract format
    // 3. Estimate gas
    // 4. Send transaction
    // 5. Wait for confirmation
    // 6. Extract token ID from events
    // 7. Create Datakyte metadata
    // 8. Return tx and token ID
}
```

### 3.2 Batch Mint
```go
func (c *Client) BatchMint(params []MintParams) (*types.Transaction, []uint64, error) {
    // 1. Validate batch size (max 100)
    // 2. Estimate total gas
    // 3. Send batch transaction
    // 4. Extract all token IDs
    // 5. Create Datakyte metadata for each
}
```

## Phase 4: Redemption with EIP-712

### 4.1 Signature Generation
```go
// internal/sdk/nft/signatures.go
func (c *Client) GenerateRedemptionSignature(
    tokenID uint64,
    redeemer common.Address,
    nonce uint64,
    deadline time.Time,
) ([]byte, error) {
    // 1. Create EIP-712 domain
    domain := apitypes.TypedDataDomain{
        Name:              "BOGOWITickets",
        Version:           "1",
        ChainId:           c.chainID,
        VerifyingContract: c.ticketsContract.Address().Hex(),
    }
    
    // 2. Create typed data
    message := map[string]interface{}{
        "tokenId":  tokenID,
        "redeemer": redeemer,
        "nonce":    nonce,
        "deadline": deadline.Unix(),
        "chainId":  c.chainID,
    }
    
    // 3. Sign with private key
    // 4. Return signature
}
```

### 4.2 Redeem Ticket
```go
func (c *Client) RedeemTicket(
    tokenID uint64,
    signature []byte,
) (*types.Transaction, error) {
    // 1. Create redemption data
    // 2. Call contract method
    // 3. Wait for confirmation
    // 4. Update Datakyte status
}
```

## Phase 5: Query Methods

### 5.1 Ticket Data Queries
```go
// internal/sdk/nft/queries.go
func (c *Client) GetTicketData(tokenID uint64) (*TicketData, error)
func (c *Client) IsTransferable(tokenID uint64) (bool, error)
func (c *Client) IsRedeemed(tokenID uint64) (bool, error)
func (c *Client) IsExpired(tokenID uint64) (bool, error)
func (c *Client) GetOwner(tokenID uint64) (common.Address, error)
func (c *Client) GetMetadataURI(tokenID uint64) (string, error)
func (c *Client) GetConservationDAOFee(tokenID uint64) (uint16, error)
```

### 5.2 Balance and Enumeration
```go
func (c *Client) GetBalance(owner common.Address) (*big.Int, error)
func (c *Client) GetTokensOfOwner(owner common.Address) ([]uint64, error)
func (c *Client) GetTotalSupply() (*big.Int, error)
```

## Phase 6: Transfer Methods

### 6.1 Standard Transfers
```go
func (c *Client) Transfer(
    from common.Address,
    to common.Address,
    tokenID uint64,
) (*types.Transaction, error)

func (c *Client) Approve(
    spender common.Address,
    tokenID uint64,
) (*types.Transaction, error)
```

## Phase 7: Event Handling

### 7.1 Event Listeners
```go
// internal/sdk/nft/events.go
func (c *Client) WatchTicketMinted(
    sink chan<- *contracts.BOGOWITicketsTicketMinted,
) (event.Subscription, error)

func (c *Client) WatchTicketRedeemed(
    sink chan<- *contracts.BOGOWITicketsTicketRedeemed,
) (event.Subscription, error)

func (c *Client) FilterTicketsByBookingID(
    bookingID [32]byte,
) ([]*contracts.BOGOWITicketsTicketMinted, error)
```

## Phase 8: Integration Helpers

### 8.1 QR Code Generation
```go
func (c *Client) GenerateRedemptionQR(
    tokenID uint64,
    redeemer common.Address,
) ([]byte, error) {
    // 1. Generate signature
    // 2. Create QR data
    // 3. Generate QR image
    // 4. Return PNG bytes
}
```

### 8.2 Metadata Sync
```go
func (c *Client) SyncWithDatakyte(tokenID uint64) error {
    // 1. Get on-chain data
    // 2. Get Datakyte metadata
    // 3. Compare and update if needed
}
```

## Phase 9: Testing

### 9.1 Unit Tests
```go
// internal/sdk/nft/client_test.go
func TestMintTicket(t *testing.T)
func TestBatchMint(t *testing.T)
func TestRedemption(t *testing.T)
func TestTransfers(t *testing.T)
func TestQueries(t *testing.T)
```

### 9.2 Integration Tests
```go
// internal/sdk/nft/integration_test.go
func TestEndToEndTicketLifecycle(t *testing.T) {
    // 1. Mint ticket
    // 2. Query data
    // 3. Transfer
    // 4. Redeem
    // 5. Verify state
}
```

## Phase 10: Documentation & Examples

### 10.1 Basic Usage Example
```go
// examples/mint_ticket.go
func main() {
    // Initialize client
    client, err := nft.NewClient(nft.ClientConfig{
        RPCURL:         "https://columbus.camino.network/ext/bc/C/rpc",
        PrivateKey:     os.Getenv("PRIVATE_KEY"),
        TicketsAddress: os.Getenv("TICKETS_CONTRACT"),
        DatakyteAPIKey: os.Getenv("DATAKYTE_API_KEY"),
        NetworkName:    "testnet",
    })
    
    // Mint a ticket
    tx, tokenID, err := client.MintTicket(nft.MintParams{
        To:                recipient,
        BookingID:         bookingHash,
        EventID:           eventHash,
        TransferUnlockAt:  time.Now().Add(24 * time.Hour).Unix(),
        ExpiresAt:         time.Now().Add(30 * 24 * time.Hour).Unix(),
        RewardBasisPoints: 500, // 5%
    })
    
    fmt.Printf("Minted ticket %d in tx %s\n", tokenID, tx.Hash())
}
```

### 10.2 Redemption Example
```go
// examples/redeem_ticket.go
func main() {
    // Generate redemption signature
    signature, err := client.GenerateRedemptionSignature(
        tokenID,
        redeemer,
        nonce,
        deadline,
    )
    
    // Redeem ticket
    tx, err := client.RedeemTicket(tokenID, signature)
    
    fmt.Printf("Redeemed ticket in tx %s\n", tx.Hash())
}
```

## Implementation Timeline

| Phase | Task | Duration | Dependencies |
|-------|------|----------|--------------|
| 1 | Generate contract bindings | 1 hour | Contract ABIs |
| 2 | Build core client structure | 2 hours | Phase 1 |
| 3 | Implement minting methods | 3 hours | Phase 2 |
| 4 | Add EIP-712 redemption | 3 hours | Phase 2 |
| 5 | Implement query methods | 2 hours | Phase 2 |
| 6 | Add transfer methods | 1 hour | Phase 2 |
| 7 | Event handling | 2 hours | Phase 2 |
| 8 | Integration helpers | 2 hours | Phases 3-7 |
| 9 | Testing | 3 hours | All phases |
| 10 | Documentation | 1 hour | All phases |

**Total Estimated Time: ~20 hours**

## Success Criteria

✅ All contract methods accessible via SDK
✅ EIP-712 signature generation working
✅ Datakyte metadata integration functional
✅ Event filtering and watching operational
✅ Comprehensive test coverage (>80%)
✅ Clear documentation with examples
✅ Gas estimation accurate
✅ Error handling robust
✅ Network switching (testnet/mainnet) seamless

## Risk Mitigation

1. **Gas Estimation**: Always estimate before sending
2. **Nonce Management**: Track and handle conflicts
3. **Event Reliability**: Implement retry logic
4. **Datakyte Sync**: Handle API failures gracefully
5. **Chain Reorgs**: Wait for confirmations
6. **Rate Limiting**: Implement request throttling