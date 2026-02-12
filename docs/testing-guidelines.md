# Testing Guidelines

This document outlines test-driven development (TDD) practices and testing best practices for the Theater Showtimes project.

## Table of Contents

- [General TDD Principles](#general-tdd-principles)
- [Go Testing Best Practices](#go-testing-best-practices)
- [React Testing Best Practices](#react-testing-best-practices)
- [Integration Testing](#integration-testing)
- [Testing Tools](#testing-tools)

---

## General TDD Principles

### The TDD Cycle (Red-Green-Refactor)

1. **Red**: Write a failing test that defines the desired functionality
2. **Green**: Write the minimum code necessary to make the test pass
3. **Refactor**: Clean up the code while keeping tests green

### Core TDD Best Practices

- **Write tests first**: Always write the test before implementing the feature
- **One test at a time**: Focus on one failing test before moving to the next
- **Small steps**: Make incremental changes and run tests frequently
- **Test behavior, not implementation**: Focus on what the code does, not how it does it
- **Keep tests simple**: Tests should be easier to understand than the code they test
- **Fast feedback**: Tests should run quickly to encourage frequent execution
- **Independent tests**: Each test should be able to run in isolation
- **Descriptive test names**: Test names should clearly describe what they're testing
- **Arrange-Act-Assert (AAA)**: Structure tests with setup, execution, and verification phases

### Benefits of TDD

- Catch bugs early in development
- Improve code design and architecture
- Provide living documentation
- Enable safe refactoring
- Increase confidence in code changes
- Reduce debugging time
- Force you to think about requirements before coding

---

## Go Testing Best Practices

### Test File Organization

```go
// File structure: place test files alongside source files
// example.go
// example_test.go

package mypackage

import "testing"
```

### Naming Conventions

```go
// Test function names: Test + FunctionName + Scenario
func TestCalculateTotal_WithValidItems_ReturnsCorrectSum(t *testing.T) {
    // test implementation
}

// Benchmark function names: Benchmark + FunctionName
func BenchmarkCalculateTotal(b *testing.B) {
    // benchmark implementation
}

// Example function names: Example + FunctionName
func ExampleCalculateTotal() {
    // example with output
    // Output: expected output
}
```

### Table-Driven Tests

Use table-driven tests to test multiple scenarios efficiently:

```go
func TestCalculateDiscount(t *testing.T) {
    tests := []struct {
        name           string
        price          float64
        discountPercent int
        expected       float64
        expectError    bool
    }{
        {
            name:           "10% discount on $100",
            price:          100.0,
            discountPercent: 10,
            expected:       90.0,
            expectError:    false,
        },
        {
            name:           "invalid negative discount",
            price:          100.0,
            discountPercent: -10,
            expected:       0,
            expectError:    true,
        },
        {
            name:           "zero discount",
            price:          50.0,
            discountPercent: 0,
            expected:       50.0,
            expectError:    false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := CalculateDiscount(tt.price, tt.discountPercent)
            
            if tt.expectError {
                if err == nil {
                    t.Errorf("expected error but got none")
                }
                return
            }
            
            if err != nil {
                t.Errorf("unexpected error: %v", err)
            }
            
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### Test Helpers and Fixtures

```go
// Helper functions for common test setup
func setupTestDB(t *testing.T) *sql.DB {
    t.Helper() // Marks this as a helper function
    
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("failed to open test database: %v", err)
    }
    
    t.Cleanup(func() {
        db.Close()
    })
    
    return db
}

// Using the helper
func TestUserRepository_Create(t *testing.T) {
    db := setupTestDB(t)
    repo := NewUserRepository(db)
    
    // Test implementation
}
```

### Mocking and Interfaces

```go
// Define interfaces for dependencies
type EmailSender interface {
    Send(to, subject, body string) error
}

// Mock implementation for testing
type MockEmailSender struct {
    SendCalled bool
    SendError  error
}

func (m *MockEmailSender) Send(to, subject, body string) error {
    m.SendCalled = true
    return m.SendError
}

// Test using the mock
func TestNotifyUser_SendsEmail(t *testing.T) {
    mockSender := &MockEmailSender{}
    service := NewNotificationService(mockSender)
    
    err := service.NotifyUser("user@example.com", "Welcome")
    
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    
    if !mockSender.SendCalled {
        t.Error("expected Send to be called")
    }
}
```

### Testing Error Handling

```go
func TestValidateInput_EmptyString_ReturnsError(t *testing.T) {
    _, err := ValidateInput("")
    
    if err == nil {
        t.Fatal("expected error for empty input, got nil")
    }
    
    // Check error message if needed
    expectedMsg := "input cannot be empty"
    if err.Error() != expectedMsg {
        t.Errorf("got error %q, want %q", err.Error(), expectedMsg)
    }
}

// Testing error wrapping
func TestProcessData_DatabaseError_WrapsError(t *testing.T) {
    result, err := ProcessData(invalidData)
    
    if err == nil {
        t.Fatal("expected error, got nil")
    }
    
    // Check if error is wrapped correctly
    if !errors.Is(err, ErrDatabaseFailure) {
        t.Errorf("expected ErrDatabaseFailure in error chain")
    }
}
```

### Test Coverage

```bash
# Run tests with coverage
go test -cover ./...

# Generate detailed coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Set minimum coverage threshold in CI
go test -cover ./... | grep "coverage:" | awk '{if ($3 < 80) exit 1}'
```

### Subtests and Parallel Tests

```go
func TestUserService(t *testing.T) {
    t.Run("Create", func(t *testing.T) {
        t.Parallel() // Run this subtest in parallel
        
        // Test user creation
    })
    
    t.Run("Update", func(t *testing.T) {
        t.Parallel()
        
        // Test user update
    })
    
    t.Run("Delete", func(t *testing.T) {
        t.Parallel()
        
        // Test user deletion
    })
}
```

### Benchmarking

```go
func BenchmarkCalculateTotal(b *testing.B) {
    items := []int{1, 2, 3, 4, 5}
    
    b.ResetTimer() // Reset timer after setup
    
    for i := 0; i < b.N; i++ {
        calculateTotal(items)
    }
}

// Run benchmarks
// go test -bench=. -benchmem
```

### Go Testing Best Practices Summary

- ✅ Use table-driven tests for multiple scenarios
- ✅ Test exported functions and methods
- ✅ Use `t.Helper()` for test helper functions
- ✅ Use `t.Cleanup()` for resource cleanup
- ✅ Use subtests with `t.Run()` for better organization
- ✅ Test error cases as thoroughly as happy paths
- ✅ Use interfaces to enable mocking
- ✅ Avoid testing implementation details
- ✅ Keep tests fast and independent
- ✅ Use meaningful test names that describe the scenario
- ✅ Aim for high test coverage but focus on critical paths
- ✅ Use `t.Parallel()` for tests that can run concurrently
- ❌ Don't use external dependencies in unit tests (use mocks)
- ❌ Don't test private functions directly
- ❌ Don't write tests that depend on execution order

---

## React Testing Best Practices

### Testing Philosophy

**Test user behavior, not implementation details**

- Test what users see and do
- Avoid testing component state or internal methods
- Focus on rendered output and user interactions

### React Testing Library Principles

```jsx
// ✅ Good: Testing user behavior
import { render, screen, fireEvent } from '@testing-library/react'
import userEvent from '@testing-library/user-event'

test('user can submit login form', async () => {
  render(<LoginForm />)
  
  const user = userEvent.setup()
  await user.type(screen.getByLabelText(/username/i), 'john')
  await user.type(screen.getByLabelText(/password/i), 'secret')
  await user.click(screen.getByRole('button', { name: /log in/i }))
  
  expect(screen.getByText(/welcome/i)).toBeInTheDocument()
})

// ❌ Bad: Testing implementation details
test('form updates state on input change', () => {
  const wrapper = shallow(<LoginForm />)
  wrapper.find('input').simulate('change', { target: { value: 'john' } })
  expect(wrapper.state('username')).toBe('john') // Don't do this
})
```

### Component Testing Patterns

#### Testing Rendering

```jsx
import { render, screen } from '@testing-library/react'
import MovieCard from './MovieCard'

describe('MovieCard', () => {
  test('renders movie title and rating', () => {
    const movie = {
      title: 'The Matrix',
      rating: 'R',
      runtime: 136,
    }
    
    render(<MovieCard movie={movie} />)
    
    expect(screen.getByText('The Matrix')).toBeInTheDocument()
    expect(screen.getByText('R')).toBeInTheDocument()
    expect(screen.getByText('136 min')).toBeInTheDocument()
  })
  
  test('renders poster image when provided', () => {
    const movie = {
      title: 'The Matrix',
      posterPath: '/poster.jpg',
    }
    
    render(<MovieCard movie={movie} />)
    
    const image = screen.getByRole('img', { name: /the matrix/i })
    expect(image).toHaveAttribute('src', '/poster.jpg')
  })
})
```

#### Testing User Interactions

```jsx
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import FilterBar from './FilterBar'

describe('FilterBar', () => {
  test('calls onTheaterChange when theater is selected', async () => {
    const mockOnChange = jest.fn()
    const theaters = [
      { id: '1', name: 'Cinema 1' },
      { id: '2', name: 'Cinema 2' },
    ]
    
    render(
      <FilterBar
        theaters={theaters}
        selectedTheater={null}
        onTheaterChange={mockOnChange}
      />
    )
    
    const user = userEvent.setup()
    const select = screen.getByLabelText(/theater/i)
    
    await user.selectOptions(select, '2')
    
    expect(mockOnChange).toHaveBeenCalledWith('2')
    expect(mockOnChange).toHaveBeenCalledTimes(1)
  })
})
```

#### Testing Async Behavior

```jsx
import { render, screen, waitFor } from '@testing-library/react'
import Home from './Home'
import * as api from '../services/api'

jest.mock('../services/api')

describe('Home', () => {
  test('displays movies after loading', async () => {
    const mockMovies = [
      { tmdb_id: 1, title: 'Movie 1' },
      { tmdb_id: 2, title: 'Movie 2' },
    ]
    
    api.getMovies.mockResolvedValue(mockMovies)
    api.getTheaters.mockResolvedValue([])
    api.getShowtimes.mockResolvedValue([])
    
    render(<Home />)
    
    // Should show loading state initially
    expect(screen.getByText(/loading/i)).toBeInTheDocument()
    
    // Wait for movies to load
    await waitFor(() => {
      expect(screen.getByText('Movie 1')).toBeInTheDocument()
      expect(screen.getByText('Movie 2')).toBeInTheDocument()
    })
    
    // Loading state should be gone
    expect(screen.queryByText(/loading/i)).not.toBeInTheDocument()
  })
  
  test('displays error message when API fails', async () => {
    api.getMovies.mockRejectedValue(new Error('API Error'))
    api.getTheaters.mockResolvedValue([])
    api.getShowtimes.mockResolvedValue([])
    
    render(<Home />)
    
    await waitFor(() => {
      expect(screen.getByText(/failed to load/i)).toBeInTheDocument()
    })
  })
})
```

#### Testing Custom Hooks

```jsx
import { renderHook, waitFor } from '@testing-library/react'
import { useMovies } from './useMovies'
import * as api from '../services/api'

jest.mock('../services/api')

describe('useMovies', () => {
  test('fetches and returns movies', async () => {
    const mockMovies = [{ id: 1, title: 'Test Movie' }]
    api.getMovies.mockResolvedValue(mockMovies)
    
    const { result } = renderHook(() => useMovies())
    
    expect(result.current.loading).toBe(true)
    
    await waitFor(() => {
      expect(result.current.loading).toBe(false)
    })
    
    expect(result.current.movies).toEqual(mockMovies)
    expect(result.current.error).toBe(null)
  })
})
```

#### Testing Forms

```jsx
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import SearchForm from './SearchForm'

describe('SearchForm', () => {
  test('submits form with entered search term', async () => {
    const mockOnSubmit = jest.fn()
    
    render(<SearchForm onSubmit={mockOnSubmit} />)
    
    const user = userEvent.setup()
    const input = screen.getByPlaceholderText(/search movies/i)
    const button = screen.getByRole('button', { name: /search/i })
    
    await user.type(input, 'Matrix')
    await user.click(button)
    
    expect(mockOnSubmit).toHaveBeenCalledWith('Matrix')
  })
  
  test('does not submit empty form', async () => {
    const mockOnSubmit = jest.fn()
    
    render(<SearchForm onSubmit={mockOnSubmit} />)
    
    const user = userEvent.setup()
    const button = screen.getByRole('button', { name: /search/i })
    
    await user.click(button)
    
    expect(mockOnSubmit).not.toHaveBeenCalled()
    expect(screen.getByText(/please enter a search term/i)).toBeInTheDocument()
  })
})
```

### Mocking API Calls

```jsx
// Option 1: Mock the entire module
import * as api from '../services/api'

jest.mock('../services/api')

test('component uses API', async () => {
  api.getMovies.mockResolvedValue([{ id: 1, title: 'Test' }])
  // test implementation
})

// Option 2: Use MSW (Mock Service Worker) for more realistic mocking
import { rest } from 'msw'
import { setupServer } from 'msw/node'

const server = setupServer(
  rest.get('/api/movies', (req, res, ctx) => {
    return res(ctx.json([{ id: 1, title: 'Test' }]))
  })
)

beforeAll(() => server.listen())
afterEach(() => server.resetHandlers())
afterAll(() => server.close())
```

### Testing Router Components

```jsx
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import Theater from './Theater'

test('displays theater showtimes', () => {
  render(
    <MemoryRouter initialEntries={['/theater/cinema-1']}>
      <Theater />
    </MemoryRouter>
  )
  
  // test implementation
})
```

### Snapshot Testing (Use Sparingly)

```jsx
import { render } from '@testing-library/react'
import MovieCard from './MovieCard'

test('matches snapshot', () => {
  const movie = {
    title: 'The Matrix',
    rating: 'R',
    runtime: 136,
  }
  
  const { container } = render(<MovieCard movie={movie} />)
  expect(container.firstChild).toMatchSnapshot()
})

// Only use snapshots for:
// - Testing component structure that rarely changes
// - Detecting unintended changes
// Don't use snapshots as a substitute for meaningful assertions
```

### React Testing Best Practices Summary

- ✅ Use `@testing-library/react` over Enzyme
- ✅ Query by accessible roles and labels (getByRole, getByLabelText)
- ✅ Use `userEvent` instead of `fireEvent` for realistic interactions
- ✅ Test loading and error states
- ✅ Mock API calls and external dependencies
- ✅ Use `waitFor` for async operations
- ✅ Test accessibility (use semantic HTML and ARIA roles)
- ✅ Write tests that resemble how users interact with your app
- ✅ Keep tests maintainable and readable
- ✅ Test critical user paths thoroughly
- ❌ Don't test implementation details (state, methods)
- ❌ Don't overuse snapshot tests
- ❌ Don't query by class names or test IDs unless necessary
- ❌ Don't test third-party libraries

---

## Integration Testing

### Go Integration Tests

```go
// Use build tags to separate unit and integration tests
// integration_test.go
//go:build integration
// +build integration

package mypackage_test

import (
    "testing"
    "database/sql"
)

func TestDatabaseIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    db := setupRealDatabase(t)
    defer db.Close()
    
    // Integration test implementation
}

// Run integration tests:
// go test -tags=integration ./...
// Skip integration tests:
// go test -short ./...
```

### React Integration Tests

```jsx
// Test entire user flows
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import App from './App'

describe('Movie browsing flow', () => {
  test('user can filter movies by theater', async () => {
    render(<App />)
    
    const user = userEvent.setup()
    
    // Wait for movies to load
    await screen.findByText('Movie 1')
    
    // Select a theater
    const theaterSelect = screen.getByLabelText(/theater/i)
    await user.selectOptions(theaterSelect, 'cinema-1')
    
    // Verify filtered results
    expect(screen.getByText('Movie 1')).toBeInTheDocument()
    expect(screen.queryByText('Movie 2')).not.toBeInTheDocument()
  })
})
```

### End-to-End Testing

Consider tools like:
- **Playwright** or **Cypress** for full E2E testing
- Test complete user journeys across frontend and backend
- Run against staging environment before production

---

## Testing Tools

### Go Testing Tools

- **Built-in `testing` package**: Standard Go testing
- **testify**: Assertion and mocking library
  ```go
  import "github.com/stretchr/testify/assert"
  assert.Equal(t, expected, actual)
  ```
- **gomock**: Mock generation tool
- **httptest**: Testing HTTP handlers
- **sqlmock**: Mocking database interactions

### React Testing Tools

- **Jest**: JavaScript testing framework (comes with Create React App)
- **React Testing Library**: DOM testing utilities
- **@testing-library/user-event**: User interaction simulation
- **MSW (Mock Service Worker)**: API mocking
- **React Hooks Testing Library**: Testing custom hooks

### Test Configuration

#### Jest Configuration (React)

```javascript
// package.json
{
  "jest": {
    "testEnvironment": "jsdom",
    "setupFilesAfterEnv": ["<rootDir>/src/setupTests.js"],
    "collectCoverageFrom": [
      "src/**/*.{js,jsx}",
      "!src/index.js",
      "!src/reportWebVitals.js"
    ],
    "coverageThreshold": {
      "global": {
        "branches": 70,
        "functions": 70,
        "lines": 70,
        "statements": 70
      }
    }
  }
}
```

#### Setup File (React)

```javascript
// src/setupTests.js
import '@testing-library/jest-dom'

// Mock window.matchMedia
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: jest.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: jest.fn(),
    removeListener: jest.fn(),
    addEventListener: jest.fn(),
    removeEventListener: jest.fn(),
    dispatchEvent: jest.fn(),
  })),
})
```

---

## TDD Workflow Examples

### Go TDD Example

```go
// 1. RED: Write failing test
func TestCalculateTotal_EmptySlice_ReturnsZero(t *testing.T) {
    result := calculateTotal([]int{})
    if result != 0 {
        t.Errorf("got %d, want 0", result)
    }
}

// 2. GREEN: Write minimal code to pass
func calculateTotal(items []int) int {
    return 0 // Simplest implementation
}

// 3. Add more tests
func TestCalculateTotal_SingleItem_ReturnsItem(t *testing.T) {
    result := calculateTotal([]int{5})
    if result != 5 {
        t.Errorf("got %d, want 5", result)
    }
}

// 4. Implement full solution
func calculateTotal(items []int) int {
    total := 0
    for _, item := range items {
        total += item
    }
    return total
}

// 5. REFACTOR: Clean up if needed
```

### React TDD Example

```jsx
// 1. RED: Write failing test
test('Counter displays initial count', () => {
  render(<Counter initialCount={5} />)
  expect(screen.getByText('Count: 5')).toBeInTheDocument()
})

// 2. GREEN: Write minimal component
function Counter({ initialCount }) {
  return <div>Count: {initialCount}</div>
}

// 3. Add interaction test
test('Counter increments when button clicked', async () => {
  render(<Counter initialCount={0} />)
  const user = userEvent.setup()
  
  await user.click(screen.getByRole('button', { name: /increment/i }))
  
  expect(screen.getByText('Count: 1')).toBeInTheDocument()
})

// 4. Implement state
function Counter({ initialCount }) {
  const [count, setCount] = useState(initialCount)
  
  return (
    <div>
      <div>Count: {count}</div>
      <button onClick={() => setCount(count + 1)}>Increment</button>
    </div>
  )
}

// 5. REFACTOR: Extract logic, improve structure
```

---

## Continuous Integration

### Run Tests in CI/CD

```yaml
# Example GitHub Actions workflow
name: Tests

on: [push, pull_request]

jobs:
  backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - run: go test -v -race -coverprofile=coverage.out ./...
      - run: go tool cover -func=coverage.out

  frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
      - run: npm ci
      - run: npm run lint
      - run: npm test -- --coverage
```

---

## Summary: Test-Driven Development Checklist

### Before Writing Code
- [ ] Write a test that describes the desired behavior
- [ ] Run the test and verify it fails (RED)

### While Writing Code
- [ ] Write minimal code to make the test pass (GREEN)
- [ ] Run all tests to ensure nothing broke
- [ ] Refactor code while keeping tests green

### After Writing Code
- [ ] Ensure all tests pass
- [ ] Check test coverage
- [ ] Review test quality and readability
- [ ] Commit code with tests

### General Guidelines
- [ ] Tests are fast and can run frequently
- [ ] Tests are independent and can run in any order
- [ ] Tests have clear, descriptive names
- [ ] Tests follow AAA pattern (Arrange, Act, Assert)
- [ ] Edge cases and error conditions are tested
- [ ] Tests focus on behavior, not implementation
- [ ] Mocks are used appropriately for external dependencies
- [ ] Test coverage is meaningful (not just high percentage)

---

## Resources

### Go Testing
- [Go Testing Package Documentation](https://pkg.go.dev/testing)
- [Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests/)
- [Testify Documentation](https://github.com/stretchr/testify)

### React Testing
- [React Testing Library Documentation](https://testing-library.com/react)
- [Jest Documentation](https://jestjs.io/)
- [Testing Library Best Practices](https://kentcdodds.com/blog/common-mistakes-with-react-testing-library)

### TDD Philosophy
- *Test Driven Development: By Example* by Kent Beck
- *Growing Object-Oriented Software, Guided by Tests* by Steve Freeman
