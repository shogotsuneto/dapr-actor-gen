# Changelog

All notable changes to this project will be documented in this file.

## [v0.0.4] - 2025-08-19

### 💥 BREAKING CHANGES

- **Optional object references are now generated as pointers**: Optional fields that use `$ref` to reference other schemas are now generated as pointer types instead of value types to enable correct `omitempty` JSON behavior

---

*Changes are documented here starting from v0.0.4. For earlier versions, please refer to the release notes.*