# Dapr Actor Code Generator

Dapr Actor Code Generator is a Go CLI tool that generates Go interfaces, types, and factory functions from OpenAPI 3.0 specifications for Dapr actors.

Always reference these instructions first and fallback to search or bash commands when you encounter unexpected information that does not match the info here, or when you need to explore unfamiliar parts of the codebase.

## Working Effectively

- **Bootstrap and build the repository:**
  - `go mod download` -- downloads dependencies, takes <1 second (cached). Set timeout to 30+ seconds.
  - `make build` -- builds the binary, takes <1 second (cached) to 16 seconds (fresh). Set timeout to 45+ seconds. NEVER CANCEL.
  - Binary is created at `./bin/dapr-actor-gen`

- **Run tests:**
  - `make test` -- runs all tests, takes <1 second (cached) to 3 seconds (fresh). Set timeout to 15+ seconds. NEVER CANCEL.
  - All tests pass consistently

- **Code formatting and linting:**
  - `go fmt ./...` -- formats all Go code, takes <1 second
  - `go vet ./...` -- lints code, takes <1 second
  - No golangci-lint or other advanced linting tools are configured

- **Clean build artifacts:**
  - `make clean` -- removes bin/, dist/, generated/, and test-output/ directories

## Core Functionality Testing

- **Generate code from OpenAPI specs:**
  - Basic generation: `./bin/dapr-actor-gen examples/multi-actors/openapi.yaml ./output` -- generates interfaces only
  - With implementation stubs: `./bin/dapr-actor-gen --generate-impl examples/multi-actors/openapi.yaml ./output` -- REQUIRED for compilable code
  - With example app: `./bin/dapr-actor-gen --generate-example examples/multi-actors/openapi.yaml ./output`
  - Full generation: `./bin/dapr-actor-gen --generate-impl --generate-example examples/multi-actors/openapi.yaml ./output`
  - Code generation is very fast (<1 second)

- **Validate generated code:**
  - Generated interfaces alone do NOT compile - you must use `--generate-impl` or provide your own implementation
  - For compilable code: `cd ./output && go mod init test-module && go mod tidy` -- downloads Dapr dependencies, takes ~2 seconds. Set timeout to 30+ seconds.
  - `cd ./output && go build ./...` -- compiles generated code, takes ~1 second after tidy. Set timeout to 15+ seconds.

## Testing Guidelines

- **Unit and integration tests:** Use OpenAPI specs from `test/integration/testdata/` directory, NOT from `examples/`
- **Manual testing and demos:** Use OpenAPI specs from `examples/` directory
- **The `examples/` directory is for:**
  - Manual testing and validation of the CLI tool
  - Demonstration purposes and README examples
  - Generating sample output for documentation
- **The `testdata/` directory is for:**
  - Automated unit and integration tests
  - CI/CD validation
  - Regression testing

## Validation Scenarios

- **ALWAYS run these validation steps after making changes:**
  - Build the tool: `make build`
  - Run all tests: `make test`
  - Generate sample code with implementation: `./bin/dapr-actor-gen --generate-impl examples/multi-actors/openapi.yaml ./test-output`
  - Verify generated code compiles: `cd ./test-output && go mod init test-module && go mod tidy && go build ./...`
  - Format code: `go fmt ./...`
  - Lint code: `go vet ./...`

- **When generator logic or example OpenAPI specs change:**
  - **Before regenerating:** Remove any existing impl.go files: `find examples/multi-actors/generated -name "impl.go" -delete`
  - Regenerate example code: `./bin/dapr-actor-gen --generate-impl --generate-example examples/multi-actors/openapi.yaml examples/multi-actors/generated`
  - **Compare with reference implementations:** For each actor type, compare exported functions between `actor.go` (reference) and `impl.go` (newly generated):
    - `go doc -u ./examples/multi-actors/generated/[actortype]/` to see all exported functions
    - Ensure `actor.go` implements all methods that `impl.go` expects from the API interface
    - Update `actor.go` if the API interface has changed (new methods, different signatures)
  - Verify the regenerated example compiles: `cd examples/multi-actors/generated && go mod tidy && go build ./...`
  - Update any documentation that references the generated code structure

- **Manual testing scenarios:**
  - Test CLI help: `./bin/dapr-actor-gen` (should show usage and exit with code 1)
  - Test basic generation with provided example OpenAPI spec (generates interfaces only - will not compile)
  - Test generation with `--generate-impl` flag (generates compilable code)
  - Test generation with `--generate-example` flag (generates complete example application)
  - Verify generated code structure matches expected format
  - ALWAYS test compilation of generated code with implementation stubs

## Documentation Maintenance

- **Before finishing any development session:**
  - Review and update `copilot-instructions.md` if any workflows or processes have changed
  - Keep instructions concise while maintaining essential information
  - Update README.md if CLI usage, features, or examples have changed
  - Verify that any referenced file paths, commands, or outputs are still accurate
  - Remove outdated information and add new essential knowledge gained during the session

- **When making changes that affect examples or generated code:**
  - Regenerate example code in `examples/multi-actors/generated/` if the generator logic has changed
  - Update any documentation that shows expected file structure or generated code samples
  - Verify that CLI help text matches documented usage

## Common Tasks

The following are outputs from frequently run commands. Reference them instead of viewing, searching, or running bash commands to save time.

### Repository structure
```
.
├── .github/
│   ├── ISSUE_TEMPLATE/
│   └── workflows/           # CI/CD pipelines
├── cmd/
│   └── main.go             # CLI entry point
├── examples/
│   ├── multi-actors/       # Example OpenAPI specifications
│   └── generated-complete/ # Example generated code
├── pkg/
│   ├── generator/          # Code generation logic
│   └── parser/             # OpenAPI parsing logic
├── test/
│   └── integration/        # Integration tests
├── Dockerfile              # Multi-stage Docker build
├── Makefile               # Build automation
├── go.mod                 # Go module definition
└── README.md              # Complete documentation
```

### Available Make targets
```
make help          # Show available targets
make build         # Build the binary (<1s cached, up to 16s fresh)
make test          # Run tests (<1s cached, up to 3s fresh)
make clean         # Clean artifacts
make tidy          # Tidy go modules
make build-all     # Cross-compile for all platforms
make dev-build     # Development build (tidy + build)
make dev-test      # Development test (tidy + test)
```

### CLI usage
```
dapr-actor-gen [flags] <openapi-file> <output-directory>

Flags:
  -generate-impl    Generate partial implementation stubs with not-implemented errors
  -generate-example Generate example main.go, go.mod and other files for a complete app
```

### Expected generated structure
```
output/
├── {actortype}/
│   ├── api.go          # Actor interface with Dapr integration
│   ├── types.go        # Type definitions from OpenAPI schemas
│   ├── factory.go      # Factory functions for registration
│   ├── impl.go         # Implementation stubs (if --generate-impl) - gitignored
│   └── actor.go        # Reference implementation (manually maintained)
├── main.go             # Example application (if --generate-example)
└── go.mod              # Go module for example (if --generate-example)
```

**Note**: `impl.go` files are generated stubs and gitignored. Use `actor.go` for reference implementations that should be maintained manually and not overwritten by code generation.

## Key Go Dependencies

- `github.com/getkin/kin-openapi` - OpenAPI 3.0 parsing
- `github.com/dapr/go-sdk` - Used in generated code (not in generator itself)
- No other external dependencies

## CI/CD Information

- Uses GitHub Actions with Go 1.23.6
- Runs on `ubuntu-latest`
- Pipeline steps: checkout → setup-go → download → build → test → clean
- Triggered on pushes to `main`/`develop` branches and PRs to `develop`

## Docker Usage

- Multi-stage build using golang:1.23.6-alpine and alpine:latest
- Binary available at `/root/dapr-actor-gen` in container
- Example: `docker run --rm -v $(pwd)/examples:/examples -v $(pwd)/output:/output image-name /examples/multi-actors/openapi.yaml /output`
- Docker builds may fail in restricted environments due to certificate issues

## Project Context

- **Purpose**: Schema-first development for Dapr actors
- **Input**: OpenAPI 3.0 specifications with actor types parsed from path patterns
- **Output**: Ready-to-implement Go actor packages
- **Design**: Generates actor-specific packages with embedded types (type duplication by design)
- **Integration**: Works with Dapr Go SDK for actor registration and execution

## Troubleshooting

- If build fails: Run `go mod tidy` first
- If generated code doesn't compile: Ensure you used `--generate-impl` flag or provide your own implementation structs
- If factory.go shows undefined struct errors: You need implementation structs - use `--generate-impl` or create your own
- If generated code can't find modules: Run `go mod init <module-name> && go mod tidy` in output directory
- If tests fail: Check that OpenAPI files exist in `test/integration/testdata/` directory (NOT examples/)
- If Docker build fails: Expected in restricted environments, try local build instead
- Always clean with `make clean` if encountering persistent build issues
- **When in doubt:** Explore the codebase and run commands to understand current state rather than relying solely on these instructions