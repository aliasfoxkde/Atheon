# Atheon Architecture

## Overview

Atheon is designed as a minimal, efficient pattern matching engine with a pluggable pattern system. The architecture prioritizes performance, simplicity, and extensibility.

## Core Components

### Pattern Engine (core/)

**bundle.go**
- Pattern loading and registration
- Category-based filtering
- Active scanner management
- Bundle compilation and decompression

**pattern.go**
- Pattern interface definition
- Pattern registry
- Pattern discovery and sorting

**runner.go**
- File system scanning
- Pattern matching engine
- Gitignore compliance
- Environment scanning

**finding.go**
- Finding structure definition
- Statistics tracking

### Pattern Bundler (bundler/)

**main.go**
- YAML pattern loading
- Bundle compilation
- JSON generation
- Gzip compression

### CLI Interface (main.go)

- Command-line parsing
- User interface
- Pattern management commands
- Output formatting

### MCP Server (cmd/mcp/)

**main.go**
- Model Context Protocol implementation
- Tool registration
- Request/response handling

## Data Flow

### Pattern Loading
1. Load embedded bundle or local bundle
2. Parse JSON and decompress gzip
3. Register all patterns
4. Build category scanners
5. Enable active patterns

### Scanning Process
1. Parse ignore files (gitignore, atheonignore)
2. Walk directory tree
3. Read file contents
4. Apply active patterns
5. Generate findings

### Pattern Matching
- Category-based pre-filtering
- Combined regex optimization
- Line-by-line scanning
- Context-aware matching

## Pattern Format

Patterns are defined as YAML files in the community/ directory:

```yaml
name: pattern-name
category: category-name
match: 'regex-pattern'
# enabled: false (optional)
```

## Category System

Patterns are organized into categories:
- **secrets**: Credentials and sensitive data
- **pii**: Personal identifiable information
- **code-quality**: Code smell detection
- **healthcare**: Medical identifiers
- **finance**: Financial patterns

## Performance Optimizations

1. **Regex Combination**: Patterns within categories are combined into single regex
2. **Early Filtering**: Gitignore and binary file filtering
3. **Concurrent Scanning**: Parallel file processing
4. **Lazy Loading**: Patterns loaded on-demand
5. **Caching**: Compiled regex patterns cached

## Extension Points

### Adding Patterns
1. Create YAML file in community/<category>/
2. Run bundler to compile bundle
3. Tests automatically discover new patterns

### Adding Categories
1. Create new directory in community/
2. Add patterns to category
3. Category automatically discovered

### Custom Scanners
1. Implement Pattern interface
2. Register with Register()
3. Available in All() results

## Build Process

1. **Pattern Compilation**: YAML → JSON → Gzip
2. **Binary Embedding**: Bundle embedded in binary
3. **Multi-platform Builds**: Cross-compilation for all platforms
4. **Package Generation**: Homebrew, Scoop formulas

## Testing Strategy

- Unit tests for each component
- Integration tests for scanning
- Pattern validation tests
- Platform compatibility tests
- Performance benchmarks

## Dependencies

**Go Standard Library:**
- regexp: Pattern matching
- compress/gzip: Bundle compression
- encoding/json: Pattern serialization
- os/filepath: File system operations

**External Dependencies:**
- github.com/sabhiram/go-gitignore: Gitignore compliance
- gopkg.in/yaml.v3: YAML parsing

## Security Considerations

- Pattern validation prevents regex attacks
- Binary file detection prevents resource exhaustion
- Gitignore compliance prevents scanning sensitive files
- No network access during scanning
- Sandboxed pattern evaluation

## Performance Characteristics

- **Memory**: < 50MB for typical usage
- **Speed**: < 100ms for 1000-line files
- **Scalability**: Linear scaling for file size
- **Concurrency**: Configurable worker pool (default: 256)