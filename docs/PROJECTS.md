# Projects Using Atheon

A non-exhaustive list of organisations and projects that
use Atheon in production. Listed in alphabetical order.
Inclusion here is a public statement of use; it is **not**
an endorsement. If you would like to be added or removed,
open a PR against this file.

> **Note:** most users run Atheon in a CI pipeline, an
> editor plugin, or a pre-commit hook, so usage is rarely
> visible at the project level. The list below mixes
> self-reported use with public CI configurations we have
> seen.

## Open source projects

| Project | Use case | How they use it |
|---|---|---|
| _(none reported yet)_ | — | — |

We will list the first three open source projects that
volunteer to be here. If your project uses Atheon, please
file a PR adding a row.

## Companies and teams

| Organisation | Use case | Contact |
|---|---|---|
| _(none reported yet)_ | — | — |

The BDFL runs Atheon in a private capacity; that does not
count. If you work somewhere that uses it, please file a
PR.

## Plugins and integrations

Projects that wrap or extend Atheon:

| Project | Description |
|---|---|
| _(none reported yet)_ | — |

## Editor and IDE support

Atheon is editor-agnostic — it is a CLI and a Go library.
The following editors have first-class support, either via
an official plugin or via the Language Server Protocol
bridge:

- **VS Code** — community extension (see Discussions).
- **Neovim** — `null-ls` integration; example config in
  `docs/guides/neovim.md`.
- **JetBrains** — use the bundled `File Watcher` template
  with the `atheon scan --stdin` command.
- **Emacs** — `flycheck` integration; example in
  `docs/guides/emacs.md`.

## CI providers

Atheon has been tested in:

- GitHub Actions
- GitLab CI
- CircleCI
- Buildkite
- Jenkins (declarative pipelines)

Example CI configurations live in `docs/guides/ci/`.

## How to be added

1. Open a PR against `docs/PROJECTS.md`.
2. Add a row to the table.
3. Include a link to a public artefact that shows the
   usage (a CI config in the repo, a blog post, a talk
   recording, etc.) — this is what makes the listing
   honest.
4. The PR will be reviewed by a maintainer; we may ask
   for clarification but we do not gatekeep on company
   size or project popularity.

If you would rather not be public, that is fine. Atheon
is used by individuals and small teams who have no
interest in being listed; we respect that.

## Removal

If a project is listed here and the maintainers no longer
want to be, open a PR (or an issue if you do not have
push access) and the row will be removed. The git history
will retain the prior version.
