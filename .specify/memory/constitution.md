<!--
SYNC IMPACT REPORT
==================
Version Change: Template → 1.0.0
Constitution Type: Initial Ratification
Principles Added:
  - I. Test-Driven Development (NON-NEGOTIABLE)
  - II. TypeScript-First Development
  - III. Modular Architecture
  - IV. Code Quality & Tooling
  - V. Performance & Caching
  - VI. Graceful Degradation
  - VII. Security & Secrets Management

Templates Reviewed:
  ✅ .specify/templates/plan-template.md - Constitution Check section aligns
  ✅ .specify/templates/spec-template.md - Requirements align with principles
  ✅ .specify/templates/tasks-template.md - Test-first approach matches TDD principle
  ✅ .specify/templates/checklist-template.md - Quality gates reference constitution

Documentation Alignment:
  ✅ docs/testing-guidelines.md - Source for TDD principle
  ✅ docs/coding-guidelines.md - Source for TypeScript, tooling, security principles
  ✅ docs/functional-requirements.md - Source for architecture, performance, degradation principles

Follow-up Actions: None required
-->

# Portland Theater Showtimes Constitution

## Core Principles

### I. Test-Driven Development (NON-NEGOTIABLE)

**TDD is mandatory for all feature development.**

- Tests MUST be written before implementation code
- Follow Red-Green-Refactor cycle strictly:
  - **Red**: Write a failing test that defines desired functionality
  - **Green**: Write minimum code necessary to make the test pass
  - **Refactor**: Clean up code while keeping tests green
- Test behavior, not implementation details
- Tests MUST run independently and in any order
- Tests MUST be fast (unit tests < 100ms, integration tests < 2s)
- All tests MUST pass before code can be committed

**Rationale**: TDD catches bugs early, improves design, provides living documentation, and enables safe refactoring. This is non-negotiable because it fundamentally shapes code quality and maintainability.

### II. TypeScript-First Development

**All frontend and shared code MUST use TypeScript with strict mode enabled.**

- Enable `"strict": true` in tsconfig.json
- Define explicit interfaces and types for all data structures
- Avoid `any` type; use `unknown` or proper types instead
- Use type inference for simple cases, be explicit for complex ones
- React components MUST use `.tsx` extension
- Prefer union types over enums for most cases
- Never use non-null assertion (`!`) except when provably safe

**Rationale**: Type safety eliminates entire classes of runtime errors, improves IDE support, makes refactoring safer, and provides self-documenting code.

### III. Modular Architecture

**Code MUST be organized in self-contained, independently testable modules.**

- Each theater scraper MUST be its own module
- Scrapers MUST not depend on each other
- Clear separation between: scrapers, API handlers, storage, TMDB client, data enrichment
- Go packages MUST have single, clear responsibilities
- React components MUST be organized by feature, not by type
- Shared utilities MUST be in dedicated packages/modules
- Configuration MUST be externalized and environment-based

**Rationale**: Modular architecture enables independent development, easier testing, simpler maintenance, and supports scaling to 20+ theaters.

### IV. Code Quality & Tooling

**Code quality tools MUST be configured and enforced at all stages.**

- ESLint MUST be configured for TypeScript/React with project standards
- Prettier MUST be used for automatic code formatting
- Pre-commit hooks MUST run linters and formatters
- All linter errors and warnings MUST be resolved before committing
- Go code MUST pass `go vet` and `golint` checks
- Table-driven tests MUST be used for Go testing multiple scenarios
- Test coverage MUST be tracked; critical paths require high coverage (>80%)

**Rationale**: Automated tooling ensures consistency, catches errors early, reduces review friction, and maintains codebase health.

### V. Performance & Caching

**System MUST meet performance targets through strategic caching.**

- API response time MUST be < 500ms with cached TMDB data
- Frontend initial load MUST be < 3s
- TMDB data MUST be cached for 7 days minimum
- Cache hit ratio MUST exceed 90%
- Scraping MUST respect rate limits (1 request per 2-3 seconds per domain)
- Concurrent scraping MUST be used where appropriate
- Images MUST use lazy loading and responsive sizing
- TMDB API calls MUST be batched and queued to avoid rate limits

**Rationale**: Performance directly impacts user experience; caching reduces external API costs and ensures reliability.

### VI. Graceful Degradation

**System MUST remain functional when external dependencies fail.**

- If TMDB match fails, display movie with title and placeholder poster
- If TMDB server unavailable, use cached data
- If scraper fails, display last successful data with timestamp
- All TMDB failures MUST be logged for manual review
- Error states MUST be user-friendly and actionable
- Loading states MUST be clear and informative
- Fallback placeholder images MUST be provided for missing posters

**Rationale**: External services (TMDB, theater websites) are unreliable; the app must function despite failures.

### VII. Security & Secrets Management

**Secrets MUST never be committed to version control.**

- API keys, passwords, tokens MUST use environment variables
- `.env` files MUST be in `.gitignore`
- `.env.example` MUST be provided with dummy values
- TMDB API key MUST be stored securely (environment variables, secure config)
- User input MUST be validated and sanitized
- Dependencies MUST be kept up to date (regular audits)
- CORS MUST be properly configured for production
- HTTPS MUST be used in production deployment

**Rationale**: Security breaches can expose user data and API costs; prevention is critical and non-negotiable.

## Technical Standards

### Language & Framework Requirements

**Backend (Go)**:
- Go 1.21 or later
- Colly v2 for web scraping
- Standard library HTTP server or Gin framework
- JSON encoding for data persistence

**Frontend (React)**:
- React 18+
- TypeScript 5.0+
- React Router for navigation
- Vitest for testing
- Axios or Fetch API for HTTP requests

**TMDB Integration**:
- MCP Server (`mcp-server-tmdb`) via Model Context Protocol
- Node.js runtime for MCP server
- HTTP client for MCP communication

### Project Structure

**Standard directory layout MUST be maintained**:

```
theater-showtimes/
├── backend/
│   ├── cmd/          # Executables (api, scraper)
│   ├── internal/     # Private application code
│   ├── configs/      # Configuration files
│   └── data/         # JSON storage
├── frontend/
│   ├── src/
│   │   ├── components/  # Reusable UI components
│   │   ├── pages/       # Page components
│   │   ├── services/    # API services
│   │   └── utils/       # Utility functions
│   └── tests/
├── mcp-servers/     # TMDB MCP server
└── docs/            # Project documentation
```

### Testing Framework Requirements

**Go Testing**:
- Built-in `testing` package
- Table-driven tests for multiple scenarios
- `t.Helper()` for test helpers
- `t.Cleanup()` for resource cleanup
- Subtests with `t.Run()` for organization

**React Testing**:
- Vitest as test runner
- React Testing Library for component testing
- `@testing-library/user-event` for interactions
- Mock Service Worker (MSW) for API mocking

## Quality Standards

### Test Coverage Requirements

- Critical paths MUST have >80% test coverage
- All API endpoints MUST have contract tests
- All scrapers MUST have integration tests
- All React components MUST have behavioral tests
- Edge cases and error conditions MUST be tested

### Code Review Requirements

- All code MUST be reviewed before merging
- All tests MUST pass in CI/CD before merge
- Linter errors MUST be resolved
- Performance impact MUST be considered
- Security implications MUST be reviewed

### Documentation Requirements

- Public APIs MUST be documented
- Complex logic MUST have explanatory comments
- README.md MUST be kept up to date
- Breaking changes MUST be documented
- Setup instructions MUST be clear and complete

## Development Workflow

### Commit Standards

- Use clear, descriptive commit messages
- Use present tense ("Add feature" not "Added feature")
- Reference issue numbers when applicable
- Keep commits atomic and focused
- Commit messages MUST explain "why" for non-obvious changes

### Pull Request Process

- Create small, focused PRs
- Write descriptive PR descriptions
- Address all review comments before merging
- Ensure CI/CD checks pass
- Follow PR template if available

### Version Control Discipline

**Never commit**:
- Secrets and credentials
- Environment files (`.env`, `.env.local`)
- Dependencies (`node_modules/`, `vendor/`)
- Build outputs (`dist/`, `build/`, binaries)
- IDE-specific files (`.vscode/`, `.idea/`)
- OS files (`.DS_Store`, `Thumbs.db`)
- Log files (`*.log`)
- Test coverage outputs
- Data files (JSON storage with development data)

## Governance

### Amendment Process

1. Proposed amendments MUST be documented with rationale
2. Version MUST be incremented according to semantic versioning:
   - **MAJOR**: Backward incompatible governance/principle removals or redefinitions
   - **MINOR**: New principle/section added or materially expanded guidance
   - **PATCH**: Clarifications, wording, typo fixes, non-semantic refinements
3. All affected templates and documentation MUST be updated
4. Sync Impact Report MUST be generated and prepended to this file

### Compliance Verification

- All specifications MUST reference applicable principles
- All implementation plans MUST include Constitution Check section
- All pull requests SHOULD verify compliance with relevant principles
- Violations MUST be justified in Complexity Tracking section of plan

### Living Document

- This constitution is the single source of truth for project governance
- Guidelines in `docs/` supplement but do not override this constitution
- When conflicts arise, this constitution takes precedence
- Use `.specify/templates/` for structured workflows and commands

**Version**: 1.0.0 | **Ratified**: 2026-02-11 | **Last Amended**: 2026-02-11
