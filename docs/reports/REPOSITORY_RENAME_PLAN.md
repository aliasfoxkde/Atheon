# Repository Rename Plan: Atheon-Enhanced
## Maintaining Upstream Compatibility

### Overview
Rename the repository from "Atheon" to "Atheon-Enhanced" while maintaining full compatibility with the upstream HoraDomu/Atheon project for PR submissions and collaboration.

### Critical Requirements

#### 1. Go Module Compatibility
- **KEEP**: `github.com/aliasfoxkde/Atheon` as the Go module path
- **DO NOT CHANGE**: `go.mod` module declaration
- **REASON**: Go modules are tied to import paths; changing this breaks all imports

#### 2. Git Repository Structure
- **KEEP**: GitHub repository path as `github.com/aliasfoxkde/Atheon`
- **CHANGE**: Display name to "Atheon-Enhanced"
- **REASON**: Git URLs remain the same, only display name changes

#### 3. Upstream PR Compatibility
- **MAINTAIN**: Ability to submit PRs to HoraDomu/Atheon
- **PRESERVE**: Code structure and import compatibility
- **ENSURE**: Cherry-pick and merge capabilities work

### Implementation Strategy

#### Phase 1: Documentation Updates (Safe)
Update display names and descriptions while keeping technical paths:

**Files to Update:**
- `README.md` - Change title/description, keep URLs
- `*.md` files - Update mentions of repository name
- `INSTALL.md` - Update installation instructions
- `docs/**/*` - Update documentation references

**Pattern:**
```markdown
# Before
**Atheon** - Pattern matching engine

# After
**Atheon-Enhanced** - Pattern matching engine (Enhanced fork of HoraDomu/Atheon)
```

#### Phase 2: Workflow Updates (Safe)
Update GitHub Actions to work with both repositories:

**Workflow Updates:**
- Update repository references in workflows
- Add upstream synchronization checks
- Ensure actions work for both repos

**Example:**
```yaml
# Before
- uses: actions/checkout@v3
  with:
    repository: aliasfoxkde/Atheon

# After
- uses: actions/checkout@v3
  with:
    repository: aliasfoxkde/Atheon  # Technical path unchanged
```

#### Phase 3: Test Compatibility (Critical)
Ensure tests and imports work with both repositories:

**Test Updates:**
- Update test comments and messages
- Ensure Go imports still work
- Verify CI/CD compatibility

#### Phase 4: GitHub Repository Settings (Manual)
Actions that require GitHub web interface:

**GitHub Settings:**
1. Rename repository on GitHub to "Atheon-Enhanced"
2. Update repository description
3. Update topics/keywords
4. Update About section

**IMPORTANT**: The GitHub URL path will stay as `github.com/aliasfoxkde/Atheon` - only the display name changes.

### What NOT to Change

#### Critical Compatibility Items
- ❌ **DO NOT CHANGE**: `go.mod` module path (`github.com/aliasfoxkde/Atheon`)
- ❌ **DO NOT CHANGE**: Import statements in Go code
- ❌ **DO NOT CHANGE**: Git remote URLs
- ❌ **DO NOT CHANGE**: GitHub Actions repository references (technical paths)
- ❌ **DO NOT CHANGE**: Test file imports and references

#### Reason
These changes would break:
- Go module resolution
- `go get` commands
- Upstream PR submissions
- Cherry-pick merges
- Import compatibility

### Documentation Updates

#### Repository Description
```markdown
# Atheon-Enhanced

**Enhanced testing fork of [HoraDomu/Atheon](https://github.com/HoraDomu/Atheon)**

This is an experimental feature-rich build for testing pattern matching limits.
**NOT a competing project** - all improvements intended for upstream submission.

## Repository Location
- **Technical Path**: `github.com/aliasfoxkde/Atheon` (unchanged for compatibility)
- **Display Name**: Atheon-Enhanced
- **Upstream**: [HoraDomu/Atheon](https://github.com/HoraDomu/Atheon)

## Compatibility
- ✅ Go module: `github.com/aliasfoxkde/Atheon` (maintained for upstream compatibility)
- ✅ PR submissions: Full compatibility with HoraDomu/Atheon
- ✅ Import paths: All imports work unchanged
- ✅ CI/CD: Works with both repositories
```

#### Installation Instructions
```bash
# Clone the repository
git clone https://github.com/aliasfoxkde/Atheon.git
cd Atheon

# Build and install
go build -o atheon .
./atheon --version
```

### Upstream Contribution Workflow

#### Submitting PRs to Upstream
```bash
# Start from upstream stable branch
git fetch upstream stable/clean
git checkout stable/clean

# Create feature branch
git checkout -b feat/my-feature

# Make changes
# ... implement feature ...

# Add upstream remote (if not already added)
git remote add upstream https://github.com/HoraDomu/Atheon.git

# Push to your fork
git push origin feat/my-feature

# Create PR to HoraDomu/Atheon
# Use GitHub web interface or gh cli:
gh pr create --repo HoraDomu/Atheon --base stable/clean
```

### Verification Checklist

#### Pre-Rename Verification
- [ ] All tests pass with current setup
- [ ] Go imports work correctly
- [ ] CI/CD workflows run successfully
- [ ] Can submit PRs to upstream

#### Post-Rename Verification
- [ ] All tests still pass
- [ ] Go imports still work
- [ ] CI/CD workflows run successfully
- [ ] Can still submit PRs to upstream
- [ ] Documentation updates are clear
- [ ] Repository description is accurate

### Files to Update

#### Display Name Changes
1. `README.md` - Title and description
2. `docs/README.md` - Documentation hub
3. `INSTALL.md` - Installation guide
4. `docs/contributing.md` - Contribution guidelines
5. `docs/FAQ.md` - Frequently asked questions

#### Technical References (Keep Paths, Update Descriptions)
1. GitHub Actions workflows (keep paths, update comments)
2. Test files (keep imports, update comments)
3. Go code (keep imports, update comments)

#### No Changes Required
1. `go.mod` - Keep module path unchanged
2. Import statements - Keep all imports unchanged
3. Git remotes - Keep remote URLs unchanged
4. GitHub Actions repository references - Keep technical paths

### Implementation Order

#### Step 1: Documentation Updates (First - Safe)
```bash
# Update README and documentation
# These changes are safe and don't affect functionality
```

#### Step 2: Code Comment Updates (Safe)
```bash
# Update comments in test files and source code
# Keep all imports and technical references unchanged
```

#### Step 3: GitHub Repository Rename (Manual)
```bash
# Use GitHub web interface to rename repository
# Settings → General → Repository Name → Atheon-Enhanced
```

#### Step 4: Verification (Critical)
```bash
# Run full test suite
# The -p 1 flag is MANDATORY: see core init() state note in docs/development/SETUP.md.
go test ./... -p 1

# Check imports work
go mod tidy

# Verify CI/CD
# Push a test commit and check workflows
```

### Rollback Plan

If any issues arise:
1. Revert documentation changes
2. Restore original repository name on GitHub
3. Verify all functionality restored
4. Investigate and fix issues

### Success Criteria

✅ **Repository renamed to "Atheon-Enhanced"**
✅ **All documentation updated correctly**
✅ **Go module compatibility maintained**
✅ **Upstream PR submission works**
✅ **All tests pass**
✅ **CI/CD workflows functional**
✅ **No breaking changes to imports or functionality**

---

**Status**: Ready for implementation
**Estimated Time**: 1-2 hours
**Risk Level**: Low (documentation-only changes)
**Rollback**: Easy (documentation reverts)