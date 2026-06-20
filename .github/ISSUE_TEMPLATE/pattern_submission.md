---
name: Pattern submission
about: Propose a new detection pattern for the community bundle
title: "pattern: "
labels: ["pattern", "needs-triage"]
assignees: []
---

## Pattern name

<!-- Lowercase, hyphenated, machine-friendly. e.g. `gitlab-deploy-token` -->

## Category

<!-- Pick one: secrets, pii, finance, healthcare, code-quality,
accessibility, networking, cloud, devops. Add a justification if
none of the existing ones fit. -->

## What does it detect?

<!-- One paragraph describing what the pattern catches and why
that's a real problem worth detecting. -->

## Sample positive matches

<!-- Real-world-looking strings that should match. Redact or alter
the actual secret portion if you copied from a real leak — the
pattern engine will still match the prefix. -->

```text
<!-- paste here -->
```

## Sample negatives

<!-- Strings that look similar but should NOT match. Negative
samples are critical for false-positive review. -->

```text
<!-- paste here -->
```

## Reference

<!-- Link to documentation, RFC, or vendor docs describing the
format. -->

## Implementation

<!-- Sketch of the YAML, especially the `match:` regex. If you
have it locally, paste the full file. -->

```yaml
name: <pattern-name>
category: <category>
enabled: true
patterns:
  - pattern: '<regex>'
```

## Notes

<!-- Anything else — performance concerns, known false positives,
related patterns in the bundle. -->
