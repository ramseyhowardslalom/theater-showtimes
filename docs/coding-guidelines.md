# Coding Guidelines

This document outlines the coding standards and best practices for the Theater Showtimes project.

## General Principles

- Write clean, readable, and maintainable code
- Follow the principle of least surprise
- Keep functions and components small and focused
- Favor composition over inheritance
- Write self-documenting code with clear naming
- Add comments only when necessary to explain "why", not "what"

## Code Quality Tools

### ESLint
- Use ESLint to enforce code quality and style rules
- Fix all linter errors and warnings before committing code
- Configure ESLint rules in `.eslintrc` to match project standards
- Run `npm run lint` or `yarn lint` before creating pull requests

### Prettier
- Use Prettier for automatic code formatting where possible
- Configure Prettier in `.prettierrc` for consistent formatting
- Enable "format on save" in your IDE
- Run Prettier on staged files using pre-commit hooks

### Pre-commit Checks
- Ensure all linter errors are resolved
- Verify code formatting is consistent
- Run tests before committing

## TypeScript

### Why TypeScript

This project uses **TypeScript** for both frontend and backend development to provide:
- Type safety and early error detection
- Better IDE support with autocomplete and IntelliSense
- Self-documenting code through type definitions
- Easier refactoring and maintenance
- Reduced runtime errors

### TypeScript Best Practices

- **Use strict mode**: Enable `"strict": true` in tsconfig.json
- **Explicit types**: Define interfaces and types for all data structures
- **Avoid `any`**: Use `unknown` or proper types instead of `any`
- **Type inference**: Let TypeScript infer simple types, be explicit for complex ones
- **Utility types**: Use built-in utility types (`Partial`, `Pick`, `Omit`, `Record`, etc.)
- **Enums vs Union Types**: Prefer union types over enums for most cases
- **Non-null assertion (`!`)**: Use sparingly, only when you're certain a value exists

### TypeScript File Extensions

- **React components**: `.tsx` (TypeScript + JSX)
- **TypeScript files**: `.ts`
- **Type definitions**: `.d.ts`
- **JavaScript files**: Avoid `.js` in new code, migrate to `.ts`/`.tsx`

### Type Definition Examples

```typescript
// Interface for component props
interface MovieCardProps {
  movie: Movie;
  showtimes: Showtime[];
  onMovieClick?: (movieId: number) => void;
}

// Type for API response
type ApiResponse<T> = {
  data: T;
  error: string | null;
  loading: boolean;
}

// Union type for states
type ViewMode = 'movies' | 'theaters';

// Type guards
function isMovie(item: Movie | Theater): item is Movie {
  return 'tmdb_id' in item;
}
```

## React Best Practices

### Component Structure

- **Use Functional Components**: Prefer functional components with hooks over class components
- **Component Organization**: One component per file, matching the filename
- **File Naming**: Use PascalCase for component files (e.g., `UserProfile.tsx`)
- **Always Use TypeScript**: Write all React components in TypeScript (.tsx)

### Hooks

- **Follow Rules of Hooks**: 
  - Only call hooks at the top level
  - Only call hooks from React functions
- **Custom Hooks**: Extract reusable logic into custom hooks
- **Dependency Arrays**: Always specify complete dependency arrays for `useEffect`, `useCallback`, and `useMemo`
- **Avoid Unnecessary Re-renders**: Use `React.memo`, `useMemo`, and `useCallback` judiciously

### State Management

- **Local State First**: Start with local component state
- **Lift State Up**: Share state by lifting it to the closest common ancestor
- **Context for Global State**: Use React Context for truly global state
- **Consider State Libraries**: For complex state, evaluate Redux, Zustand, or Jotai

### Component Best Practices

```tsx
// Good: Clear props interface, destructured props, early returns
interface UserCardProps {
  userId: string;
  userName: string;
  onUserClick?: (userId: string) => void;
}

export const UserCard: React.FC<UserCardProps> = ({ 
  userId, 
  userName, 
  onUserClick 
}) => {
  if (!userId) return null;
  
  const handleClick = () => {
    onUserClick?.(userId);
  };
  
  return (
    <div onClick={handleClick}>
      <h3>{userName}</h3>
    </div>
  );
};
```

### Performance

- **Lazy Loading**: Use `React.lazy()` for code splitting
- **Avoid Index as Key**: Use stable, unique identifiers for list keys
- **Optimize Images**: Use appropriate formats and sizes
- **Virtualize Long Lists**: Use libraries like `react-window` for large lists

### Testing

- **Write Tests**: Aim for meaningful test coverage
- **Test User Behavior**: Focus on testing what users see and do
- **Use Testing Library**: Prefer `@testing-library/react` over Enzyme
- **Mock External Dependencies**: Mock API calls and external services

### Code Organization

```
src/
├── components/          # Reusable UI components
│   ├── common/         # Shared components (Button, Input, etc.)
│   └── features/       # Feature-specific components
├── hooks/              # Custom React hooks
├── pages/              # Page components
├── services/           # API and external service integrations
├── utils/              # Utility functions
├── types/              # TypeScript type definitions
└── constants/          # Application constants
```

## Go Best Practices

### Code Structure

- **Package Naming**: Use short, lowercase, single-word package names
- **File Organization**: Group related functionality in the same package
- **Exported vs Unexported**: 
  - Capitalize exported identifiers
  - Keep unexported what doesn't need to be public

### Error Handling

```go
// Good: Explicit error handling
func getUserByID(id string) (*User, error) {
    user, err := db.FindUser(id)
    if err != nil {
        return nil, fmt.Errorf("failed to find user %s: %w", id, err)
    }
    return user, nil
}

// Avoid: Ignoring errors
func getUserByID(id string) *User {
    user, _ := db.FindUser(id)  // BAD: Error ignored
    return user
}
```

### Best Practices

- **Error Wrapping**: Use `fmt.Errorf` with `%w` to wrap errors with context
- **Defer for Cleanup**: Use `defer` to ensure cleanup happens
- **Avoid Panic**: Reserve `panic` for truly unrecoverable situations
- **Use Interfaces**: Define interfaces where you use them, not where you implement them
- **Keep It Simple**: Prefer clarity over cleverness

### Concurrency

```go
// Good: Proper error handling with goroutines
func processItems(items []Item) error {
    errChan := make(chan error, len(items))
    
    for _, item := range items {
        item := item // Capture loop variable
        go func() {
            errChan <- processItem(item)
        }()
    }
    
    for i := 0; i < len(items); i++ {
        if err := <-errChan; err != nil {
            return err
        }
    }
    return nil
}
```

- **Use Channels**: Communicate by sharing memory, don't share memory to communicate
- **Context for Cancellation**: Pass `context.Context` for cancellation and timeouts
- **Avoid Goroutine Leaks**: Ensure all goroutines can exit
- **Sync Primitives**: Use `sync.WaitGroup`, `sync.Mutex` appropriately

### Naming Conventions

- **Variables**: Use camelCase, be descriptive (`userCount`, not `uc`)
- **Constants**: Use camelCase or MixedCaps, not SCREAMING_SNAKE_CASE
- **Interfaces**: Single-method interfaces end in "-er" (`Reader`, `Writer`)
- **Getters**: Don't use "Get" prefix (`user.Name()`, not `user.GetName()`)

### Project Structure

```
project/
├── cmd/                # Application entry points
│   └── server/
│       └── main.go
├── internal/           # Private application code
│   ├── handlers/       # HTTP handlers
│   ├── models/         # Data models
│   ├── repository/     # Data access layer
│   └── service/        # Business logic
├── pkg/                # Public libraries
└── go.mod
```

### Testing

```go
// Good: Table-driven tests
func TestCalculateTotal(t *testing.T) {
    tests := []struct {
        name     string
        input    []int
        expected int
    }{
        {"empty slice", []int{}, 0},
        {"single item", []int{5}, 5},
        {"multiple items", []int{1, 2, 3}, 6},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := calculateTotal(tt.input)
            if result != tt.expected {
                t.Errorf("got %d, want %d", result, tt.expected)
            }
        })
    }
}
```

- **Table-Driven Tests**: Use for testing multiple scenarios
- **Test Coverage**: Aim for high coverage, but focus on critical paths
- **Benchmarks**: Use `testing.B` for performance testing
- **Test Names**: Use descriptive names that explain what's being tested

## Version Control

### Commits

- Write clear, descriptive commit messages
- Use present tense ("Add feature" not "Added feature")
- Reference issue numbers when applicable
- Keep commits atomic and focused

### Pull Requests

- Create small, focused PRs
- Write descriptive PR descriptions
- Request reviews from appropriate team members
- Address all review comments before merging
- Ensure CI/CD checks pass

### .gitignore Best Practices

**Never commit the following:**

- **Secrets and credentials**: API keys, passwords, tokens, certificates
- **Environment files**: `.env`, `.env.local`, `.env.production`
- **Dependencies**: `node_modules/`, `vendor/`
- **Build outputs**: `dist/`, `build/`, compiled binaries
- **IDE files**: `.vscode/`, `.idea/`, `*.swp`
- **OS files**: `.DS_Store`, `Thumbs.db`
- **Logs**: `*.log`, debug logs
- **Test coverage**: `coverage/`, `*.out`
- **Temporary files**: `tmp/`, `temp/`, `*.tmp`
- **Data files**: Database files, JSON storage (development data)

**Go Projects (.gitignore)**
```
# Binaries and build output
*.exe
*.dll
*.so
*.dylib
*.test
*.out
/api
/scraper

# Dependencies
vendor/

# Data and logs
/data/
*.log

# Coverage
coverage.out
```

**React Projects (.gitignore)**
```
# Dependencies
node_modules/

# Build output
/build
/dist

# Environment files
.env.local
.env.*.local

# Testing
/coverage

# Logs
npm-debug.log*
yarn-error.log*
```

**Important:**
- Keep `.gitignore` files up to date as project evolves
- Use `.gitignore` templates for your language/framework
- Add `.env.example` with dummy values to show required variables
- Commit lock files (`package-lock.json`, `go.sum`) for reproducible builds

## Documentation

- Document public APIs and exported functions
- Keep README.md up to date
- Add inline comments for complex logic
- Update documentation when changing behavior

## Security

- Never commit secrets, API keys, or passwords
- Use environment variables for configuration
- Validate and sanitize user input
- Keep dependencies up to date
- Follow OWASP best practices for web applications
