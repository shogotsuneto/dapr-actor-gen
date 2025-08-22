# Changelog

This changelog documents all significant changes that users should be aware of when upgrading, including breaking changes, new features, bug fixes, improvements, and deprecations.

## [v0.0.4] - 2025-08-19

### 💥 BREAKING CHANGES

- **Optional object references are now generated as pointers**: Optional fields that use `$ref` to reference other schemas are now generated as pointer types (e.g., `*CounterStateData`, `*UserId`, `*OperationStatus`) instead of value types. This enables correct `omitempty` JSON behavior where `nil` pointers are excluded from serialization.

  **Migration required**: 
  - Regenerate all actors using the updated tool
  - Update existing actor implementations to handle pointer types
  - When setting optional fields, use pointer assignment (e.g., `&someValue` or appropriate helper functions)
  
  **Before (v0.0.3)**:
  ```go
  type CounterState struct {
      Data  CounterStateData `json:"data,omitempty"`    // Empty struct {} included in JSON
      UserId UserId          `json:"userId,omitempty"`   // Empty string "" included in JSON
  }
  ```
  
  **After (v0.0.4)**:
  ```go
  type CounterState struct {
      Data  *CounterStateData `json:"data,omitempty"`   // nil pointer omitted from JSON
      UserId *UserId          `json:"userId,omitempty"`  // nil pointer omitted from JSON
  }
  ```

---

*Breaking changes are documented here starting from v0.0.4. For earlier versions, please refer to the release notes.*