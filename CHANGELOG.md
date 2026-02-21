# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.0.3] - 2025-02-21

### Added
- MIT LICENSE file for pkg.go.dev compatibility
- Codecov integration for coverage reporting

## [0.0.2] - 2025-02-21

### Added
- GoDoc documentation for all exported functions
- Framework examples for Gin, Echo, Chi, and Fiber in README
- CONTRIBUTING.md with contribution guidelines
- Issue templates (bug report, feature request)
- Pull request template
- CHANGELOG.md following Keep a Changelog format
- Example functions for pkg.go.dev documentation
- Go Report Card badge in README

## [0.0.1] - 2025-02-21

### Added
- Core filter engine with type-safe generics
- Comparison operators: `Eq`, `Ne`, `Gt`, `Lt`, `Gte`, `Lte`
- String operations: `Contains`, `StringMatch`, `RegexMatch`
- Collection filters: `In`, `Between`, `ArrayContains`, `ArrayContainsAny`, `ArrayContainsAll`
- Logical composition: `And`, `Or`, `Not`
- Nil/Zero checks: `IsNil`, `IsNotNil`, `IsZero`, `IsNotZero`
- Date filters: `DateBefore`, `DateAfter`, `DateBetween`
- Geospatial filters: `WithinRadius`, `OutsideRadius`, `WithinBoundingBox`, `SortByDistance`
- Map field filters: `HasKey`, `HasValue`, `KeyValueEquals`, `MapContainsAll`, `MapContainsAny`, `MapSizeEquals`, `MapSizeGreaterThan`, `MapSizeLessThan`
- Query package for HTTP query parameter parsing
- Django-style query syntax (`field_gt`, `field_contains`, `field_in`, etc.)
- Automatic type coercion from query strings
- Pagination support with `ApplyPaginated`
- Struct tag system for secure field exposure (`gofilter:"filterable,sortable"`)
- Custom filter function support
- Nested field access with dot notation
- Sorting functionality with `Sort` and `SortByDistance`

[Unreleased]: https://github.com/sidneip/gofilter/compare/v0.0.3...HEAD
[0.0.3]: https://github.com/sidneip/gofilter/compare/v0.0.2...v0.0.3
[0.0.2]: https://github.com/sidneip/gofilter/compare/v0.0.1...v0.0.2
[0.0.1]: https://github.com/sidneip/gofilter/releases/tag/v0.0.1
