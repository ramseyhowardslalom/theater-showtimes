import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import MovieCard from './MovieCard'

describe('MovieCard', () => {
    const mockMovie = {
        tmdb_id: 123,
        title: 'Test Movie',
        rating: 'PG-13',
        runtime: 120,
        tmdb_rating: 7.5,
        overview: 'This is a test movie overview that should be truncated if it is too long to fit in the card display area. This needs to be longer than 120 characters.',
        poster_path: 'https://image.tmdb.org/t/p/w500/test.jpg',
    }

    const mockShowtimes = [
        { id: '1', tmdb_id: 123, time: '14:00', date: '2026-02-11', theater_id: 'theater-1' },
        { id: '2', tmdb_id: 123, time: '19:00', date: '2026-02-11', theater_id: 'theater-1' },
        { id: '3', tmdb_id: 123, time: '21:30', date: '2026-02-11', theater_id: 'theater-2' },
        { id: '4', tmdb_id: 123, time: '23:00', date: '2026-02-11', theater_id: 'theater-2' },
    ]

    const mockTheaters = [
        { id: 'theater-1', name: 'Clinton Street Theater' },
        { id: 'theater-2', name: 'Hollywood Theatre' },
    ]

    it('should render movie title', () => {
        render(<MovieCard movie={mockMovie} showtimes={[]} theaters={[]} />)
        expect(screen.getByText('Test Movie')).toBeInTheDocument()
    })

    it('should render movie metadata (rating, runtime, tmdb_rating)', () => {
        render(<MovieCard movie={mockMovie} showtimes={[]} theaters={[]} />)
        expect(screen.getByText('PG-13')).toBeInTheDocument()
        expect(screen.getByText('120 min')).toBeInTheDocument()
        expect(screen.getByText(/⭐ 7.5/)).toBeInTheDocument()
    })

    it('should render truncated movie overview with More button', () => {
        render(<MovieCard movie={mockMovie} showtimes={[]} theaters={[]} />)
        const overview = screen.getByText(/This is a test movie overview/)
        expect(overview).toBeInTheDocument()
        expect(overview.textContent).toContain('...')
        expect(screen.getByRole('button', { name: /More/i })).toBeInTheDocument()
    })

    it('should expand overview when More button is clicked', async () => {
        const user = userEvent.setup()
        render(<MovieCard movie={mockMovie} showtimes={[]} theaters={[]} />)

        const moreButton = screen.getByRole('button', { name: /More/i })
        await user.click(moreButton)

        expect(screen.getByRole('button', { name: /Less/i })).toBeInTheDocument()
        const overview = screen.getByText(/This is a test movie overview that should be truncated if it is too long to fit in the card display area./)
        expect(overview).toBeInTheDocument()
    })

    it('should render poster image when poster_path exists', () => {
        render(<MovieCard movie={mockMovie} showtimes={[]} theaters={[]} />)
        const img = screen.getByAltText('Test Movie')
        expect(img).toBeInTheDocument()
        expect(img).toHaveAttribute('src', mockMovie.poster_path)
    })

    it('should render "No Poster" placeholder when poster_path is missing', () => {
        const movieWithoutPoster = { ...mockMovie, poster_path: null }
        render(<MovieCard movie={movieWithoutPoster} showtimes={[]} theaters={[]} />)
        expect(screen.getByText('No Poster')).toBeInTheDocument()
    })

    it('should render showtimes in Pacific Time format', () => {
        render(<MovieCard movie={mockMovie} showtimes={mockShowtimes} theaters={mockTheaters} />)
        expect(screen.getByText('2:00 PM PT')).toBeInTheDocument()
        expect(screen.getByText('7:00 PM PT')).toBeInTheDocument()
        expect(screen.getByText('9:30 PM PT')).toBeInTheDocument()
    })

    it('should show only first 3 showtimes', () => {
        render(<MovieCard movie={mockMovie} showtimes={mockShowtimes} theaters={mockTheaters} />)
        expect(screen.getByText('2:00 PM PT')).toBeInTheDocument()
        expect(screen.getByText('7:00 PM PT')).toBeInTheDocument()
        expect(screen.getByText('9:30 PM PT')).toBeInTheDocument()
        expect(screen.queryByText('11:00 PM PT')).not.toBeInTheDocument()
    })

    it('should show "+N more" indicator when there are more than 3 showtimes', () => {
        render(<MovieCard movie={mockMovie} showtimes={mockShowtimes} theaters={mockTheaters} />)
        expect(screen.getByText('+1 more')).toBeInTheDocument()
    })

    it('should not show "+N more" when there are 3 or fewer showtimes', () => {
        const fewerShowtimes = mockShowtimes.slice(0, 2)
        render(<MovieCard movie={mockMovie} showtimes={fewerShowtimes} theaters={mockTheaters} />)
        expect(screen.queryByText(/\+\d+ more/)).not.toBeInTheDocument()
    })

    it('should render theater labels', () => {
        render(<MovieCard movie={mockMovie} showtimes={mockShowtimes} theaters={mockTheaters} />)
        expect(screen.getByText('Clinton Street Theater')).toBeInTheDocument()
        expect(screen.getByText('Hollywood Theatre')).toBeInTheDocument()
    })

    it('should render unique theater labels only', () => {
        const singleTheaterShowtimes = mockShowtimes.filter(st => st.theater_id === 'theater-1')
        render(<MovieCard movie={mockMovie} showtimes={singleTheaterShowtimes} theaters={mockTheaters} />)
        expect(screen.getByText('Clinton Street Theater')).toBeInTheDocument()
        expect(screen.queryByText('Hollywood Theatre')).not.toBeInTheDocument()
    })

    it('should filter showtimes to only show those matching the movie', () => {
        const mixedShowtimes = [
            ...mockShowtimes,
            { id: '5', tmdb_id: 456, time: '10:00', date: '2026-02-11', theater_id: 'theater-1' }, // Different movie
        ]
        render(<MovieCard movie={mockMovie} showtimes={mixedShowtimes} theaters={mockTheaters} />)

        // Should still show showtimes for movie 123, not 456
        expect(screen.getByText('2:00 PM PT')).toBeInTheDocument()
        expect(screen.queryByText('10:00 AM PT')).not.toBeInTheDocument()
    })

    it('should handle empty showtimes array', () => {
        render(<MovieCard movie={mockMovie} showtimes={[]} theaters={[]} />)
        expect(screen.getByText('Test Movie')).toBeInTheDocument()
        expect(screen.queryByText(/PT/)).not.toBeInTheDocument()
    })

    it('should handle missing overview gracefully', () => {
        const movieWithoutOverview = { ...mockMovie, overview: null }
        render(<MovieCard movie={movieWithoutOverview} showtimes={[]} theaters={[]} />)
        expect(screen.getByText('Test Movie')).toBeInTheDocument()
    })

    it('should not show More button for short overview', () => {
        const movieWithShortOverview = { ...mockMovie, overview: 'Short overview' }
        render(<MovieCard movie={movieWithShortOverview} showtimes={[]} theaters={[]} />)
        expect(screen.queryByRole('button', { name: /More/i })).not.toBeInTheDocument()
    })
})
