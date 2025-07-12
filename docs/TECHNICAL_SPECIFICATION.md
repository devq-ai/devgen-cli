# DevGen CLI - Technical Specification

## Overview
DevGen CLI is a customer-facing command-line interface for managing MCP (Model Context Protocol) servers with integrated knowledge management capabilities.

## Current Architecture

### Core Components
- **Dashboard**: Interactive TUI for MCP server management
- **Registry**: HTTP-based MCP server discovery and management
- **SSH Server**: Remote terminal access

### Technology Stack
- **Language**: Go 1.23+
- **UI Framework**: Charm (Bubble Tea, Lipgloss)
- **Configuration**: JSON-based registry (`mcp_status.json`)
- **Build System**: Make + Go modules

## Proposed Enhancements

### 1. Knowledge Base Integration (`devgen kb`)

#### Database Statistics Component
```bash
devgen kb stats              # Show database statistics
devgen kb stats --db surrealdb  # Specific database stats
devgen kb stats --format json   # JSON output
```

**Implementation Requirements:**
- Integration with SurrealDB (primary) and Neo4j
- Real-time metrics dashboard
- Performance monitoring
- Storage usage analytics
- Query performance statistics

**Data Sources:**
- Ptolemies knowledge base (784-page documentation)
- Context7 Redis cache
- Neo4j relationship graphs

#### Features:
- Document count by source
- Vector embedding statistics
- Query response times
- Cache hit rates
- Storage utilization

### 2. RAG Lookup System (`devgen search`)

#### Retrieval-Augmented Generation Interface
```bash
devgen search "FastAPI authentication"     # Semantic search
devgen search --type code "error handling" # Code-specific search
devgen search --graph "API relationships"  # Graph traversal
devgen search --source "anthropic docs"    # Source-filtered search
```

**Implementation Requirements:**
- Vector similarity search (embeddings)
- Hybrid search (semantic + keyword)
- Knowledge graph traversal
- Source attribution
- Relevance scoring

**Integration Points:**
- Context7 MCP server (vector search)
- Ptolemies MCP server (knowledge base)
- SurrealDB MCP server (multi-model queries)

#### Search Capabilities:
- Semantic similarity using OpenAI embeddings
- Code pattern matching
- Documentation lookup
- Relationship discovery
- Multi-source aggregation

### 3. DeHallucinator Application (`devgen dehall`)

#### AI Hallucination Detection Tool
```bash
devgen dehall check --input response.txt   # Check text for hallucinations
devgen dehall verify --claim "fact"        # Verify specific claims
devgen dehall analyze --source code.py     # Analyze code explanations
devgen dehall report --format html         # Generate hallucination report
```

**Implementation Requirements:**
- Fact verification against knowledge base
- Code analysis for accuracy
- Citation validation
- Confidence scoring
- Interactive correction suggestions

**Detection Methods:**
- Cross-reference with verified sources
- Code execution validation
- API documentation verification
- Logical consistency checks
- Temporal fact validation

#### Integration with Existing Systems:
- Uses Ptolemies knowledge base for fact-checking
- Leverages Context7 for contextual verification
- Connects to GitHub MCP for code validation

## Command Structure Enhancement

### Current Commands
```
devgen dashboard    # Interactive MCP server dashboard
devgen list        # List all MCP servers
devgen registry    # Registry management
devgen ssh         # SSH server for remote access
devgen toggle      # Toggle server status
```

### Proposed New Commands
```
devgen kb          # Knowledge base management
├── stats          # Database statistics
├── search         # Knowledge search
├── import         # Import data sources
└── export         # Export knowledge

devgen search      # RAG lookup system
├── semantic       # Semantic similarity search
├── code          # Code-specific search
├── graph         # Knowledge graph traversal
└── sources       # Source-filtered search

devgen dehall      # DeHallucinator tool
├── check         # Check for hallucinations
├── verify        # Verify specific claims
├── analyze       # Analyze content accuracy
└── report        # Generate reports
```

## Technical Implementation

### 1. Knowledge Base Integration

#### Database Connection Manager
```go
type KnowledgeBase struct {
    SurrealDB *surrealdb.Client
    Neo4j     *neo4j.Driver
    Context7  *context7.Client
}

func (kb *KnowledgeBase) GetStats() (*DBStats, error) {
    // Aggregate statistics from all databases
}
```

#### Statistics Collection
- Document counts and sizes
- Vector embedding metrics
- Query performance data
- Cache utilization
- Relationship graph statistics

### 2. RAG Search Engine

#### Search Interface
```go
type SearchEngine struct {
    VectorStore   VectorDB
    KnowledgeGraph GraphDB
    DocumentStore DocumentDB
}

type SearchRequest struct {
    Query      string
    Type       SearchType // semantic, code, graph
    Sources    []string
    MaxResults int
    Threshold  float64
}

type SearchResult struct {
    Content     string
    Source      string
    Relevance   float64
    Citations   []Citation
    Metadata    map[string]interface{}
}
```

#### Search Capabilities
- Hybrid vector + keyword search
- Source attribution and filtering
- Relevance scoring and ranking
- Result aggregation and deduplication

### 3. DeHallucinator Architecture

#### Verification Engine
```go
type DeHallucinator struct {
    KnowledgeBase *KnowledgeBase
    FactChecker   *FactChecker
    CodeValidator *CodeValidator
}

type HallucinationCheck struct {
    Input       string
    Claims      []Claim
    Verifications []Verification
    ConfidenceScore float64
    Recommendations []string
}

type Claim struct {
    Text        string
    Category    ClaimType // fact, code, citation
    Confidence  float64
    Sources     []Source
    IsVerified  bool
}
```

#### Verification Methods
- Knowledge base cross-reference
- Code execution testing
- API documentation validation
- Citation verification
- Logical consistency analysis

## Integration with Existing MCP Servers

### Required MCP Server Connections
1. **Context7 MCP** - Vector search and contextual reasoning
2. **Ptolemies MCP** - Knowledge base access
3. **SurrealDB MCP** - Multi-model database operations
4. **GitHub MCP** - Code repository access for validation

### Data Flow
```
User Query → DevGen CLI → MCP Servers → Knowledge Sources → Results
```

## Performance Requirements

### Response Times
- Knowledge base stats: < 1 second
- Search queries: < 2 seconds  
- Hallucination checks: < 5 seconds
- Dashboard refresh: < 500ms

### Scalability
- Support for 10K+ documents
- Concurrent user sessions
- Real-time statistics updates
- Efficient caching strategies

## Security Considerations

### Data Protection
- Secure MCP server connections
- API key management
- Rate limiting for search queries
- Audit logging for verification requests

### Access Control
- User authentication for sensitive features
- Role-based permissions
- Secure credential storage
- Network security for remote access

## Development Phases

### Phase 1: Knowledge Base Integration
- Database statistics dashboard
- Basic search functionality
- MCP server integration

### Phase 2: Advanced Search
- RAG implementation
- Knowledge graph traversal
- Multi-source search

### Phase 3: DeHallucinator
- Fact verification engine
- Code validation
- Report generation

### Phase 4: Production Hardening
- Performance optimization
- Security enhancements
- Documentation completion
- Customer deployment tools

## Testing Strategy

### Unit Tests
- Individual component testing
- Mock MCP server responses
- Database operation validation

### Integration Tests
- End-to-end search workflows
- MCP server communication
- Cross-component interaction

### Performance Tests
- Load testing for search queries
- Database performance benchmarks
- Concurrent user simulation

## Deployment

### Distribution Methods
- Single binary distribution
- Docker containerization
- Package manager integration
- Cloud deployment options

### Configuration Management
- Environment-based settings
- MCP server auto-discovery
- Credential management
- Feature flags

---

**Document Version**: 1.0  
**Last Updated**: 2025-07-12  
**Status**: Draft - Pending Implementation