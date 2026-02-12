# Frontend Unit Tests

## Overview

Comprehensive unit tests for the Portland Movie Theater Showtimes frontend application.

## Test Suite

### Total Coverage
- **43 tests** across 3 test files
- **100% passing** ✅

### Test Files

#### 1. `src/utils/dateTime.test.js` (21 tests)
Tests for date/time utility functions with Pacific timezone support:

- **formatTimeToPacific**: Converts 24-hour time to 12-hour format with "PT" suffix
  - ✓ Formats time to 12-hour format with PT
  - ✓ Handles noon and midnight correctly
  - ✓ Handles invalid input gracefully
  - ✓ Preserves minutes

- **formatDateToPacific**: Formats dates in readable format
  - ✓ Formats date correctly
  - ✓ Handles invalid input gracefully

- **formatDateToYYYYMMDD**: Converts Date objects to YYYY-MM-DD format in PT
  - ✓ Formats date to YYYY-MM-DD
  - ✓ Handles invalid input gracefully
  - ✓ Pads single digit months and days

- **filterShowtimesByDate**: Filters showtimes by selected date
  - ✓ Filters showtimes by date
  - ✓ Returns empty array for date with no showtimes
  - ✓ Handles invalid input gracefully
  - ✓ Does not modify original array

- **filterShowtimesByTheater**: Filters showtimes by theater ID
  - ✓ Filters showtimes by theater
  - ✓ Returns empty array for unknown theater
  - ✓ Handles invalid input gracefully

- **getMoviesWithShowtimes**: Gets movies that have showtimes on a specific date
  - ✓ Returns movies with showtimes on selected date
  - ✓ Returns all movies when no date selected
  - ✓ Returns empty array when no movies match
  - ✓ Handles invalid input gracefully
  - ✓ Returns unique movies only

#### 2. `src/components/MovieCard/MovieCard.test.jsx` (13 tests)
Tests for the MovieCard component:

- ✓ Renders movie title
- ✓ Renders movie metadata (rating, runtime, TMDB rating)
- ✓ Renders movie overview (truncated)
- ✓ Renders poster image when poster_path exists
- ✓ Renders "No Poster" placeholder when poster_path is missing
- ✓ Renders showtimes in Pacific Time format
- ✓ Shows only first 3 showtimes
- ✓ Shows "+N more" indicator when there are more than 3 showtimes
- ✓ Does not show "+N more" when there are 3 or fewer showtimes
- ✓ Renders "View Details" link with correct href
- ✓ Filters showtimes to only show those matching the movie
- ✓ Handles empty showtimes array
- ✓ Handles missing overview gracefully

#### 3. `src/pages/Home.test.jsx` (9 tests)
Tests for the Home page component with date/theater filtering:

- ✓ Shows loading state initially
- ✓ Loads and displays data from API
- ✓ Shows error message when API fails
- ✓ Filters movies by selected date
- ✓ Filters movies by theater
- ✓ Combines date and theater filters
- ✓ Switches to theater list view
- ✓ Passes filtered showtimes to theater list view
- ✓ Renders DateSelector and FilterBar

## Bugs Fixed

### 1. Date Filtering Not Working ✅
**Problem**: The date selector was displayed but didn't actually filter movies by date.

**Solution**: 
- Added `useMemo` hooks in [Home.jsx](../src/pages/Home.jsx) to filter showtimes by selected date and theater
- Created utility functions in [dateTime.js](../src/utils/dateTime.js) for date filtering logic
- Movies now show only those with showtimes on the selected date

### 2. Times Not in Pacific Timezone ✅
**Problem**: Showtimes were displayed in 24-hour format without timezone indication.

**Solution**:
- Created `formatTimeToPacific()` function that converts 24-hour time to 12-hour AM/PM format
- Added "PT" suffix to all showtime displays
- Example: `19:00` → `7:00 PM PT`

## Running Tests

### Run all tests once
```bash
npm test -- --run
```

### Run tests in watch mode (re-runs on file changes)
```bash
npm test
```

### Run tests with coverage
```bash
npm test:coverage
```

### Run tests with UI
```bash
npm test:ui
```

## Test Configuration

- **Framework**: Vitest 4.0
- **Testing Library**: @testing-library/react 16.1
- **Environment**: jsdom (simulates browser)
- **Config**: [vitest.config.js](../vitest.config.js)
- **Setup**: [src/test/setup.js](../src/test/setup.js)

## Dependencies

```json
{
  "vitest": "^4.0.18",
  "@testing-library/react": "^16.1.0",
  "@testing-library/jest-dom": "^6.6.3",
  "@testing-library/user-event": "^14.5.2",
  "jsdom": "^24.0.0"
}
```

## Date/Time Utilities

All date/time functions use **Pacific Time (America/Los_Angeles)** as the timezone:

- Showtimes are displayed in PT
- Date filtering works in PT
- Today's date defaults to current time in PT

## Test Patterns

### Component Testing
- Uses React Testing Library for component rendering
- Tests user interactions with `userEvent`
- Waits for async updates with `waitFor`
- Mocks child components to isolate tests

### Utility Testing
- Pure function tests with various inputs
- Edge case handling (null, undefined, invalid formats)
- Immutability checks (doesn't modify original arrays)

### API Mocking
- Mocks API calls with Vitest's `vi.fn()`
- Tests both success and error states
- Verifies API is called with correct parameters

## Future Improvements

- Add integration tests for full user flows
- Add visual regression tests for UI components
- Add E2E tests with Playwright or Cypress
- Increase code coverage to >90%
- Add performance tests for filtering large datasets
