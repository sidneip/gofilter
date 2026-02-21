---
name: Feature Request
about: Suggest an idea for gofilter
title: '[FEATURE] '
labels: enhancement
assignees: ''
---

## Problem Statement

A clear description of the problem you're trying to solve.

Example: "I want to filter users by email domain, but there's no suffix match operator..."

## Proposed Solution

Describe how you'd like this to work.

```go
// Example API usage
result := filter.Apply(users, filter.EndsWith[User]("Email", "@gmail.com"))
```

## Alternatives Considered

Any alternative solutions or workarounds you've considered.

## Additional Context

- Use cases where this would be helpful
- Links to similar features in other libraries
- Any other context
