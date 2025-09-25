# Blockchain Study Project - Development Guidelines

## Project Overview

This is a blockchain study project with three main components:
- **blockchain-v1**: Basic blockchain implementation with transactions and validation
- **blockchain-v2**: Enhanced version with proof-of-work mining and blockchain forking scenarios
- **crypto**: Cryptographic utilities demonstrating SHA-256 hashing and ECDSA key generation/signing

## Build/Configuration Instructions

### Prerequisites
- Go 1.16+ installed
- No external dependencies required (uses only Go standard library)

### Building Components

**Important**: This project does not use Go modules and should be built using traditional Go build commands.

#### blockchain-v1
```bash
cd blockchain-v1
go build main.go
# Creates executable: main (or main.exe on Windows)
```

#### blockchain-v2
```bash
cd blockchain-v2
go build main.go mining.go
# Creates executable: main (or main.exe on Windows)
# Note: Both files must be specified together due to cross-file dependencies
```

#### crypto
```bash
cd crypto
go build hash.go
# Creates executable: hash (or hash.exe on Windows)
```

### Running Components
```bash
# Run blockchain-v1
cd blockchain-v1 && go run main.go

# Run blockchain-v2
cd blockchain-v2 && go run main.go mining.go

# Run crypto demo
cd crypto && go run hash.go
```

## Testing Information

### Test Configuration
- Tests use Go's built-in testing framework
- No external testing dependencies required
- Tests must be run by specifying both main and test files explicitly

### Running Tests

#### Basic Test Execution
```bash
cd blockchain-v1
go test main.go main_test.go -v
```

The `-v` flag provides verbose output showing individual test results.

#### Example Test Output
```
=== RUN   TestCreateGenesisBlock
--- PASS: TestCreateGenesisBlock (0.00s)
=== RUN   TestGenerateNewBlock
--- PASS: TestGenerateNewBlock (0.00s)
=== RUN   TestIsBlockValid
--- PASS: TestIsBlockValid (0.00s)
=== RUN   TestTransaction
--- PASS: TestTransaction (0.00s)
PASS
ok      command-line-arguments  0.002s
```

### Adding New Tests

#### Test File Structure
- Test files should be named `*_test.go`
- Place test files in the same directory as the code being tested
- Use package `main` to match the main package

#### Test Function Template
```go
func TestFunctionName(t *testing.T) {
    // Arrange
    // Set up test data
    
    // Act
    // Call the function being tested
    
    // Assert
    // Verify the results
    if actual != expected {
        t.Errorf("Expected %v, got %v", expected, actual)
    }
}
```

#### Example Test Implementation
```go
func TestCreateGenesisBlock(t *testing.T) {
    genesis := createGenesisBlock()
    
    if genesis.Index != 0 {
        t.Errorf("Expected genesis block index to be 0, got %d", genesis.Index)
    }
    
    if genesis.PrevHash != "0" {
        t.Errorf("Expected genesis block previous hash to be '0', got %s", genesis.PrevHash)
    }
    
    if genesis.Hash == "" {
        t.Error("Expected genesis block to have a hash")
    }
}
```

### Testing Guidelines
1. **Test Core Logic**: Focus on testing blockchain validation, block generation, and transaction handling
2. **Edge Cases**: Test invalid blocks, wrong hashes, incorrect indices
3. **Data Integrity**: Verify hash calculations and block linking
4. **Error Conditions**: Test validation failures and boundary conditions

## Additional Development Information

### Code Architecture

#### blockchain-v1
- **Simple Implementation**: Basic blockchain with transaction validation
- **Key Functions**:
  - `createGenesisBlock()`: Creates the first block
  - `generateNewBlock()`: Creates new blocks with transactions
  - `isBlockValid()`: Validates block integrity
  - `calculateHash()`: Computes SHA-256 hash for blocks

#### blockchain-v2
- **Enhanced Features**: Includes proof-of-work mining and forking scenarios
- **Additional Components**:
  - `mining.go`: Implements proof-of-work algorithm
  - `mineBlock()`: Mines blocks with specified difficulty
  - `isValidHash()`: Validates hash difficulty requirements
- **Mining Process**: Uses nonce iteration to find hashes with leading zeros

#### crypto
- **Cryptographic Demos**: Shows hashing and digital signatures
- **Key Features**:
  - SHA-256 hashing examples
  - ECDSA key generation and signing
  - Signature verification

### Development Best Practices

#### File Organization
- Keep related functionality in separate files (e.g., mining.go)
- Use descriptive function names following Go conventions
- Maintain consistent package structure

#### Error Handling
- Use explicit validation in `isBlockValid()` function
- Return meaningful error messages in tests
- Handle edge cases (genesis block, nil pointers)

#### Performance Considerations
- Mining operations are CPU-intensive (difficulty 4+ may take significant time)
- Hash calculations use SHA-256 from crypto/sha256
- Consider progress reporting for long-running mining operations

#### Code Style
- Follow Go formatting conventions (use `gofmt`)
- Use meaningful variable names (e.g., `prevBlock`, `newBlock`)
- Include comments for complex logic, especially in mining algorithms
- Maintain consistent indentation and spacing

### Debugging Tips

#### Common Issues
1. **Module Errors**: Use file-specific build commands instead of `go build .`
2. **Test Failures**: Ensure both main and test files are specified in test commands
3. **Hash Mismatches**: Verify timestamp and transaction data consistency
4. **Mining Performance**: Lower difficulty for testing (1-2 zeros)

#### Debugging Commands
```bash
# Build with race detection
go build -race main.go

# Run with verbose output
go run -v main.go

# Test with coverage
go test main.go main_test.go -cover
```

### Future Development Considerations

#### Potential Enhancements
- Add network communication between nodes
- Implement transaction pools and mempool management
- Add wallet functionality with key management
- Implement consensus mechanisms beyond proof-of-work
- Add persistence layer for blockchain storage

#### Scalability Notes
- Current implementation stores entire blockchain in memory
- Consider database integration for larger blockchains
- Mining difficulty should be adjustable based on network conditions
- Transaction validation could be optimized for higher throughput