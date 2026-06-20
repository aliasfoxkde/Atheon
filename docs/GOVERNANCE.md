# Governance

How decisions get made, who is responsible for what, and how the
project runs day-to-day. Pair this with
[`MAINTAINERS.md`](MAINTAINERS.md), which lists the people in
each role.

## Principles

1. **Open by default.** Discussions, issues, and design docs
   happen on GitHub. Private channels are reserved for security
   disclosures and personnel matters.
2. **Consensus first.** Maintainers prefer to reach agreement
   in a public thread. If consensus fails, the BDFL falls back
   to the decision process below.
3. **Reversible until shipped.** Once a release is tagged, the
   changes in it are permanent unless a follow-up release
   reverts them. Until then, any merged PR can be reverted.
4. **Smallest possible change.** When solving a problem, prefer
   the patch that touches the fewest files, the smallest API,
   and the most explicit configuration. Optimisation and
   reorganisation come in follow-up PRs.

## Roles

| Role | Responsibility |
|---|---|
| **User** | Files issues, opens PRs, runs the scanner. |
| **Contributor** | Has a merged PR. May review future PRs. |
| **Pattern author** | Maintains one or more YAML patterns in `community/`. |
| **Triager** | Triages new issues and PRs. Does not have merge rights. |
| **Maintainer** | Has merge rights. Owns a subsystem. |
| **Benevolent Dictator For Life (BDFL)** | Final tie-breaker on cross-cutting decisions. Currently `@aliasfoxkde`. |

A person may hold more than one role.

## Decision process

Most decisions are made in the PR that implements them. The
PR's review thread is the discussion; merging the PR is the
decision.

For decisions that do not fit in a single PR (e.g. adopting a
new dependency, deprecating a public API, changing the release
cadence), the process is:

1. **Discussion.** Open a GitHub Discussion under the
   *Ideas* category. Mark it as *roadmap* if it concerns the
   long-term plan.
2. **Lazy consensus.** After 7 days, if no maintainer objects
   in the thread, the proposal is accepted.
3. **Objections.** If a maintainer objects, the BDFL makes the
   final call after the discussion has run for at least 14
   days.
4. **Emergency.** Security fixes and broken-CI fixes may be
   merged without waiting for the timer. The PR description
   must call out why the timer was skipped.

## Subsystems and owners

Each subsystem has a maintainer who is the primary reviewer
for changes in it. The CODEOWNERS file at
`.github/CODEOWNERS` is the source of truth; the table below
is human-readable.

| Subsystem | Path | Owner |
|---|---|---|
| Core scanner | `core/` | `@aliasfoxkde` |
| CLI | `cmd/atheon/` | `@aliasfoxkde` |
| MCP server | `cmd/mcp/` | `@aliasfoxkde` |
| Pattern bundle generator | `bundler/` | `@aliasfoxkde` |
| Community patterns | `community/` | `@aliasfoxkde` |
| CI/CD | `.github/workflows/`, `.githooks/` | `@aliasfoxkde` |
| Documentation | `docs/`, `*.md` | `@aliasfoxkde` |
| Releases | goreleaser config | `@aliasfoxkde` |

When adding new subsystems, update both CODEOWNERS and this
table in the same PR.

## Becoming a maintainer

Contributors who have shipped several high-quality PRs across
multiple subsystems may be invited to become maintainers. The
process is:

1. A current maintainer nominates the contributor in a
   private maintainer channel.
2. Two other maintainers second the nomination.
3. The BDFL confirms.
4. The new maintainer is added to CODEOWNERS, granted merge
   rights, and announced in the Discussions thread.

There is no fixed timeline. People are invited when their
contribution pattern shows they will say "yes" to the role.

## Stepping down

Maintainers may step down at any time by:

1. Opening a PR that removes their name from CODEOWNERS and
   `MAINTAINERS.md`.
2. Sending a brief note in the maintainer channel so the
   remaining maintainers can rebalance ownership.

There is no expectation to give advance warning. Maintainership
is a volunteer role, not an obligation.

## Conflict of interest

Maintainers who have a financial interest in a PR (employer,
contractor, customer of a feature) must disclose it in the PR
description and recuse themselves from the review if a
non-conflicted maintainer is available.

## Removing a maintainer

The removal process mirrors the addition process:

1. A maintainer raises the concern in the maintainer channel,
   with specific evidence.
2. The discussion runs for at least 14 days.
3. Two-thirds of the other maintainers agree that removal is
   appropriate.
4. The BDFL confirms.

The removed maintainer is informed privately before the
public PR that updates CODEOWNERS and MAINTAINERS.md.

## Commercial support

Atheon itself is open source and free. Commercial support,
custom integrations, and SLA-backed hosting are available
through the maintainers; reach out via the channels listed at
the bottom of [`MAINTAINERS.md`](MAINTAINERS.md).