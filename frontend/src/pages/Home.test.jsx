import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import userEvent from '@testing-library/user-event'
import Home from './Home'
import * as api from '../services/api'

// Mock the API module
vi.mock('../services/api', () => ({
    getShowtimes: vi.fn(),
    getMovies: vi.fn(),
    getTheaters: vi.fn(),
}))

// Mock child components to simplify testing
vi.mock('../components/MovieCard/MovieCard', () => ({
    default: ({ movie }) => <div data-testid={`movie-${movie.tmdb_id}`}>{movie.title}</div>,
}))

vi.mock('../components/FilterBar/FilterBar', () => ({
    default: ({ onTheaterChange, onViewModeChange }) => (
        <div data-testid="filter-bar">
            <button onClick={() => onTheaterChange('theater-1')}>Filter Theater</button>
            <button onClick={() => onViewModeChange('theaters')}>View Theaters</button>
        </div>
    ),
}))

vi.mock('../components/DateSelector/DateSelector', () => ({
    default: ({ onChange }) => (
        <div data-testid="date-selector">
            <button onClick={() => onChange(new Date('2026-02-11T12:00:00'))}>Feb 11</button>
            <button onClick={() => onChange(new Date('2026-02-12T12:00:00'))}>Feb 12</button>
        </div>
    ),
}))

vi.mock('../components/TheaterList/TheaterList', () => ({
    default: ({ showtimes }) => (
        <div data-testid="theater-list">
            Showing {showtimes.length} showtimes
        </div>
    ),
}))

vi.mock('../components/LoadingState/LoadingState', () => ({
    default: () => <div data-testid="loading">Loading...</div>,
}))

const renderWithRouter = (component) => {
    return render(<BrowserRouter>{component}</BrowserRouter>)
}

describe('Home', () => {
    const mockMovies = [
        { tmdb_id: 1, title: 'Movie 1' },
        { tmdb_id: 2, title: 'Movie 2' },
        { tmdb_id: 3, title: 'Movie 3' },
    ]

    const mockShowtimes = [
        { id: '1', tmdb_id: 1, date: '2026-02-11', time: '14:00', theater_id: 'theater-1' },
        { id: '2', tmdb_id: 2, date: '2026-02-11', time: '19:00', theater_id: 'theater-1' },
        { id: '3', tmdb_id: 1, date: '2026-02-12', time: '20:00', theater_id: 'theater-2' },
        { id: '4', tmdb_id: 3, date: '2026-02-12', time: '15:00', theater_id: 'theater-1' },
    ]

    const mockTheaters = [
        { id: 'theater-1', name: 'Theater 1' },
        { id: 'theater-2', name: 'Theater 2' },
    ]

    beforeEach(() => {
        vi.clearAllMocks()
        api.getShowtimes.mockResolvedValue(mockShowtimes)
        api.getMovies.mockResolvedValue(mockMovies)
        api.getTheaters.mockResolvedValue(mockTheaters)
    })

    it('should show loading state initially', () => {
        renderWithRouter(<Home />)
        expect(screen.getByTestId('loading')).toBeInTheDocument()
    })

    it('should load and display data from API', async () => {
        renderWithRouter(<Home />)

        await waitFor(() => {
            expect(api.getShowtimes).toHaveBeenCalled()
            expect(api.getMovies).toHaveBeenCalled()
            expect(api.getTheaters).toHaveBeenCalled()
        })

        // By default, should show movies with showtimes for today (Feb 11, 2026)
        // Today's date is Feb 11 based on context, so only Movie 1 and Movie 2 should show
        await waitFor(() => {
            expect(screen.getByText('Movie 1')).toBeInTheDocument()
            expect(screen.getByText('Movie 2')).toBeInTheDocument()
            // Movie 3 only has showtimes on Feb 12
            expect(screen.queryByText('Movie 3')).not.toBeInTheDocument()
        })
    })

    it('should show error message when API fails', async () => {
        api.getShowtimes.mockRejectedValue(new Error('Network error'))

        renderWithRouter(<Home />)

        await waitFor(() => {
            expect(screen.getByText('Failed to load data. Please try again.')).toBeInTheDocument()
        })
    })

    it('should filter movies by selected date', async () => {
        renderWithRouter(<Home />)
        const user = userEvent.setup()

        // Wait for initial load - defaults to today (Feb 11), showing Movie 1 and Movie 2
        await waitFor(() => {
            expect(screen.getByText('Movie 1')).toBeInTheDocument()
            expect(screen.getByText('Movie 2')).toBeInTheDocument()
        })

        // Click Feb 12 date button to change filter
        const feb12Button = screen.getByText('Feb 12')
        await user.click(feb12Button)

        // Should show only movies with showtimes on Feb 12 (Movie 1 and Movie 3)
        await waitFor(() => {
            expect(screen.getByTestId('movie-1')).toBeInTheDocument()
            expect(screen.getByTestId('movie-3')).toBeInTheDocument()
            // Movie 2 only has showtimes on Feb 11, so it shouldn't be displayed
            expect(screen.queryByTestId('movie-2')).not.toBeInTheDocument()
        })
    })

    it('should filter movies by theater', async () => {
        renderWithRouter(<Home />)
        const user = userEvent.setup()

        // Wait for initial load - should show Feb 11 movies by default
        await waitFor(() => {
            expect(screen.getByText('Movie 1')).toBeInTheDocument()
        })

        // Click filter theater button to filter by theater-1
        const filterButton = screen.getByText('Filter Theater')
        await user.click(filterButton)

        // Should show only movies in theater-1 ON FEB 11 (current date)
        // Feb 11, theater-1: Movie 1 and Movie 2
        await waitFor(() => {
            expect(screen.getByTestId('movie-1')).toBeInTheDocument()
            expect(screen.getByTestId('movie-2')).toBeInTheDocument()
            // Movie 3's Feb 11 showtime would not exist, it only has Feb 12
            expect(screen.queryByTestId('movie-3')).not.toBeInTheDocument()
        })
    })

    it('should combine date and theater filters', async () => {
        renderWithRouter(<Home />)
        const user = userEvent.setup()

        // Wait for initial load
        await waitFor(() => {
            expect(screen.getByText('Movie 1')).toBeInTheDocument()
        })

        // Select Feb 12
        const feb12Button = screen.getByText('Feb 12')
        await user.click(feb12Button)

        // Filter by theater-1
        const filterButton = screen.getByText('Filter Theater')
        await user.click(filterButton)

        // Should show only movies in theater-1 on Feb 12
        // Feb 12, theater-1: Only Movie 3 (id 4)
        // Movie 1 on Feb 12 is in theater-2, not theater-1
        await waitFor(() => {
            expect(screen.getByTestId('movie-3')).toBeInTheDocument()
            expect(screen.queryByTestId('movie-1')).not.toBeInTheDocument()
            expect(screen.queryByTestId('movie-2')).not.toBeInTheDocument()
        })
    })

    it('should switch to theater list view', async () => {
        renderWithRouter(<Home />)
        const user = userEvent.setup()

        // Wait for initial load
        await waitFor(() => {
            expect(screen.getByText('Movie 1')).toBeInTheDocument()
        })

        // Switch to theaters view
        const viewTheatersButton = screen.getByText('View Theaters')
        await user.click(viewTheatersButton)

        await waitFor(() => {
            expect(screen.getByTestId('theater-list')).toBeInTheDocument()
            // Movie cards should not be visible
            expect(screen.queryByTestId('movie-1')).not.toBeInTheDocument()
        })
    })

    it('should pass filtered showtimes to theater list view', async () => {
        renderWithRouter(<Home />)
        const user = userEvent.setup()

        // Wait for initial load - defaults to today (Feb 11)
        await waitFor(() => {
            expect(screen.getByText('Movie 1')).toBeInTheDocument()
        })

        // Theater list in Feb 11 view should show 2 showtimes (both in theater-1)
        // Switch to theaters view
        const viewTheatersButton = screen.getByText('View Theaters')
        await user.click(viewTheatersButton)

        // Theater list should show only Feb 11 showtimes (2 showtimes)
        await waitFor(() => {
            expect(screen.getByText('Showing 2 showtimes')).toBeInTheDocument()
        })
    })

    it('should render DateSelector and FilterBar', async () => {
        renderWithRouter(<Home />)

        await waitFor(() => {
            expect(screen.getByTestId('date-selector')).toBeInTheDocument()
            expect(screen.getByTestId('filter-bar')).toBeInTheDocument()
        })
    })
})
