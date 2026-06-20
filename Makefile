# Atheon Makefile
#
# Single source of truth for build / test / audit / lint targets.
# Used locally and from CI workflows (see .github/workflows/).
#
# Phony targets are declared per-target so `make` doesn't get
# confused when a file in the repo shadows a target name.

GO               ?= go
GOFLAGS          ?=
COVER_PROFILE    ?= .tmp-cov/cover.out
COVER_THRESHOLD  ?= 54
PKGS             ?= ./...

.PHONY: help build test race coverage vet staticcheck audit \
        audit-dead-code audit-nolint audit-fixmes audit-sentinels \
        audit-snapshot fmt clean

# ---- meta ----------------------------------------------------------------

help: ## list available targets
	@awk 'BEGIN {FS = ":.*##"; printf "Targets:\n"} \
	/^[a-zA-Z_-]+:.*##/ {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}' \
	$(MAKEFILE_LIST)

# ---- build / test --------------------------------------------------------

build: ## compile every package
	$(GO) build $(GOFLAGS) $(PKGS)

test: ## run unit tests (no race)
	$(GO) test $(GOFLAGS) $(PKGS)

race: ## run unit tests with -race
	$(GO) test $(GOFLAGS) -race $(PKGS)

coverage: ## run tests with coverage, print summary
	@mkdir -p $(dir $(COVER_PROFILE))
	$(GO) test $(GOFLAGS) -race -coverprofile=$(COVER_PROFILE) -timeout 15m $(PKGS)
	@$(GO) tool cover -func=$(COVER_PROFILE) | grep -E '^total:'
	@COV=$$($(GO) tool cover -func=$(COVER_PROFILE) | awk '/^total:/ {print $$3}' | sed 's/%//'); \
	if [ $${COV%.*} -lt $(COVER_THRESHOLD) ]; then \
		echo "coverage $$COV% < threshold $(COVER_THRESHOLD)%"; exit 1; \
	fi

fmt: ## run gofmt + goimports (if available)
	@gofmt -l .
	@if command -v goimports >/dev/null 2>&1; then goimports -l .; fi

clean: ## remove generated artefacts
	rm -rf .tmp-cov/ coverage.out

# ---- individual quality gates -------------------------------------------

vet: ## go vet static analysis
	$(GO) vet $(PKGS)

staticcheck: ## staticcheck if installed (warn and skip when missing)
	@if command -v staticcheck >/dev/null 2>&1; then \
		staticcheck $(PKGS); \
	else \
		echo "staticcheck not installed; install with: go install honnef.co/go/tools/cmd/staticcheck@latest"; \
	fi

# ---- audit gates ---------------------------------------------------------
#
# Each audit gate is independently runnable. The umbrella `audit`
# target runs them all in sequence and reports a roll-up.

audit: vet staticcheck audit-dead-code audit-nolint audit-fixmes audit-sentinels ## run every audit gate

audit-dead-code: ## find unexported helpers with no production callers
	@echo "--- audit: dead code (unexported helpers with no callers)"
	@bash scripts/audit-dead-code.sh all

audit-nolint: ## enumerate every //nolint annotation
	@echo "--- audit: //nolint annotations"
	@HITS=$$(grep -rn '//nolint' --include='*.go' . 2>/dev/null || true); \
	if [ -z "$$HITS" ]; then echo "  none"; else echo "$$HITS"; fi

audit-fixmes: ## enumerate TODO/FIXME/XXX markers
	@echo "--- audit: TODO/FIXME/XXX"
	@HITS=$$(grep -rnE '// *(TODO|FIXME|XXX)' --include='*.go' . 2>/dev/null || true); \
	if [ -z "$$HITS" ]; then echo "  none"; else echo "$$HITS"; fi

audit-sentinels: ## verify every exported sentinel error has callers
	@echo "--- audit: exported sentinel errors must have callers"
	@FAILED=0; \
	for err in $$(grep -hE '^var Err[A-Z][A-Za-z0-9_]* *=' core/*.go 2>/dev/null \
		| sed -E 's/^var ([A-Za-z0-9_]+).*/\1/' | sort -u); do \
		USES=$$(grep -rlE "\\b$$err\\b" --include='*.go' . 2>/dev/null \
			| grep -v 'errors.go' | wc -l); \
		if [ "$$USES" -lt 1 ]; then \
			echo "  $$err: no callers"; FAILED=1; \
		fi; \
	done; \
	if [ $$FAILED -eq 0 ]; then echo "  OK"; fi

# ---- snapshots ----------------------------------------------------------

audit-snapshot: ## regenerate docs/audits/DEAD_CODE_AUDIT.md
	@echo "snapshot: regenerate docs/audits/DEAD_CODE_AUDIT.md (manual edit recommended)"
	@echo "see scripts/audit-snapshot.sh for the helper"