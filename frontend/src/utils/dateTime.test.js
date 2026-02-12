import { describe, it, expect } from 'vitest'
import {
    formatTimeToPacific,
    formatDateToPacific,
    formatDateToYYYYMMDD,
    filterShowtimesByDate,
    filterShowtimesByTheater,
    getMoviesWithShowtimes,
} from './dateTime'

describe('formatTimeToPacific', () => {
    it('should format time to 12-hour format with PT', () => {
        expect(formatTimeToPacific('14:00')).toBe('2:00 PM PT')
        expect(formatTimeToPacific('09:30')).toBe('9:30 AM PT')
        expect(formatTimeToPacific('00:00')).toBe('12:00 AM PT')
        expect(formatTimeToPacific('23:45')).toBe('11:45 PM PT')
    })

    it('should handle noon and midnight correctly', () => {
        expect(formatTimeToPacific('12:00')).toBe('12:00 PM PT')
        expect(formatTimeToPacific('00:00')).toBe('12:00 AM PT')
    })

    it('should handle invalid input gracefully', () => {
        expect(formatTimeToPacific('')).toBe('')
        expect(formatTimeToPacific(null)).toBe('')
        expect(formatTimeToPacific(undefined)).toBe('')
        expect(formatTimeToPacific('invalid')).toBe('invalid')
        expect(formatTimeToPacific('25:00')).toContain('1:00 AM PT') // Wraps around
    })

    it('should preserve minutes', () => {
        expect(formatTimeToPacific('15:05')).toBe('3:05 PM PT')
        expect(formatTimeToPacific('08:15')).toBe('8:15 AM PT')
    })
})

describe('formatDateToPacific', () => {
    it('should format date correctly', () => {
        expect(formatDateToPacific('2026-02-11')).toBe('Feb 11, 2026')
        expect(formatDateToPacific('2026-12-31')).toBe('Dec 31, 2026')
        expect(formatDateToPacific('2026-01-01')).toBe('Jan 1, 2026')
    })

    it('should handle invalid input gracefully', () => {
        expect(formatDateToPacific('')).toBe('')
        expect(formatDateToPacific(null)).toBe('')
        expect(formatDateToPacific('invalid')).toBe('invalid')
    })
})

describe('formatDateToYYYYMMDD', () => {
    it('should format date to YYYY-MM-DD', () => {
        const date = new Date('2026-02-11T12:00:00')
        const formatted = formatDateToYYYYMMDD(date)
        expect(formatted).toMatch(/^\d{4}-\d{2}-\d{2}$/)
        expect(formatted).toContain('2026')
    })

    it('should handle invalid input gracefully', () => {
        expect(formatDateToYYYYMMDD(null)).toBe('')
        expect(formatDateToYYYYMMDD(undefined)).toBe('')
        expect(formatDateToYYYYMMDD('invalid')).toBe('')
        expect(formatDateToYYYYMMDD(new Date('invalid'))).toBe('')
    })

    it('should pad single digit months and days', () => {
        const date = new Date('2026-01-05T12:00:00')
        const formatted = formatDateToYYYYMMDD(date)
        expect(formatted).toMatch(/^\d{4}-\d{2}-\d{2}$/)
        expect(formatted).toContain('2026')
        expect(formatted).toContain('01')
        expect(formatted).toContain('05')
    })
})

describe('filterShowtimesByDate', () => {
    const mockShowtimes = [
        { id: '1', date: '2026-02-11', time: '14:00', tmdb_id: 1 },
        { id: '2', date: '2026-02-11', time: '19:00', tmdb_id: 2 },
        { id: '3', date: '2026-02-12', time: '20:00', tmdb_id: 1 },
        { id: '4', date: '2026-02-13', time: '15:00', tmdb_id: 3 },
    ]

    it('should filter showtimes by date', () => {
        const date = new Date('2026-02-11T12:00:00')
        const filtered = filterShowtimesByDate(mockShowtimes, date)
        expect(filtered).toHaveLength(2)
        expect(filtered[0].id).toBe('1')
        expect(filtered[1].id).toBe('2')
    })

    it('should return empty array for date with no showtimes', () => {
        const date = new Date('2026-02-20T12:00:00')
        const filtered = filterShowtimesByDate(mockShowtimes, date)
        expect(filtered).toHaveLength(0)
    })

    it('should handle invalid input gracefully', () => {
        expect(filterShowtimesByDate(null, new Date())).toEqual([])
        expect(filterShowtimesByDate(undefined, new Date())).toEqual([])
        expect(filterShowtimesByDate([], new Date())).toEqual([])
        expect(filterShowtimesByDate(mockShowtimes, null)).toEqual(mockShowtimes)
    })

    it('should not modify original array', () => {
        const date = new Date('2026-02-11T12:00:00')
        const originalLength = mockShowtimes.length
        filterShowtimesByDate(mockShowtimes, date)
        expect(mockShowtimes).toHaveLength(originalLength)
    })
})

describe('filterShowtimesByTheater', () => {
    const mockShowtimes = [
        { id: '1', theater_id: 'theater-a', tmdb_id: 1 },
        { id: '2', theater_id: 'theater-b', tmdb_id: 2 },
        { id: '3', theater_id: 'theater-a', tmdb_id: 3 },
        { id: '4', theater_id: 'theater-c', tmdb_id: 1 },
    ]

    it('should filter showtimes by theater', () => {
        const filtered = filterShowtimesByTheater(mockShowtimes, 'theater-a')
        expect(filtered).toHaveLength(2)
        expect(filtered[0].id).toBe('1')
        expect(filtered[1].id).toBe('3')
    })

    it('should return empty array for unknown theater', () => {
        const filtered = filterShowtimesByTheater(mockShowtimes, 'theater-z')
        expect(filtered).toHaveLength(0)
    })

    it('should handle invalid input gracefully', () => {
        expect(filterShowtimesByTheater(null, 'theater-a')).toEqual([])
        expect(filterShowtimesByTheater(undefined, 'theater-a')).toEqual([])
        expect(filterShowtimesByTheater(mockShowtimes, null)).toEqual(mockShowtimes)
        expect(filterShowtimesByTheater(mockShowtimes, '')).toEqual(mockShowtimes)
    })
})

describe('getMoviesWithShowtimes', () => {
    const mockMovies = [
        { tmdb_id: 1, title: 'Movie 1' },
        { tmdb_id: 2, title: 'Movie 2' },
        { tmdb_id: 3, title: 'Movie 3' },
    ]

    const mockShowtimes = [
        { id: '1', date: '2026-02-11', tmdb_id: 1 },
        { id: '2', date: '2026-02-11', tmdb_id: 2 },
        { id: '3', date: '2026-02-12', tmdb_id: 1 },
    ]

    it('should return movies with showtimes on selected date', () => {
        const date = new Date('2026-02-11T12:00:00')
        const result = getMoviesWithShowtimes(mockMovies, mockShowtimes, date)
        expect(result).toHaveLength(2)
        expect(result.map(m => m.tmdb_id)).toEqual([1, 2])
    })

    it('should return all movies when no date selected', () => {
        const result = getMoviesWithShowtimes(mockMovies, mockShowtimes, null)
        expect(result).toHaveLength(3)
    })

    it('should return empty array when no movies match', () => {
        const date = new Date('2026-02-20T12:00:00')
        const result = getMoviesWithShowtimes(mockMovies, mockShowtimes, date)
        expect(result).toHaveLength(0)
    })

    it('should handle invalid input gracefully', () => {
        expect(getMoviesWithShowtimes(null, mockShowtimes, new Date())).toEqual([])
        expect(getMoviesWithShowtimes(mockMovies, null, new Date())).toEqual(mockMovies)
        expect(getMoviesWithShowtimes([], mockShowtimes, new Date())).toEqual([])
    })

    it('should return unique movies only', () => {
        const date = new Date('2026-02-12T12:00:00')
        const result = getMoviesWithShowtimes(mockMovies, mockShowtimes, date)
        expect(result).toHaveLength(1)
        expect(result[0].tmdb_id).toBe(1)
    })
})
