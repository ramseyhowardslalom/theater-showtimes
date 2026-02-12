import { format, parseISO } from 'date-fns'

/**
 * Formats a time string (HH:mm) to Pacific Time (PT) with AM/PM
 * Assumes the input time is in PT already from the scraper
 * @param {string} timeString - Time in HH:mm format (e.g., "19:00")
 * @returns {string} Formatted time with AM/PM (e.g., "7:00 PM PT")
 */
export function formatTimeToPacific(timeString) {
    if (!timeString || typeof timeString !== 'string') {
        return ''
    }

    // Parse the time string (HH:mm)
    const [hours, minutes] = timeString.split(':').map(Number)

    if (isNaN(hours) || isNaN(minutes)) {
        return timeString
    }

    // Create a date object for today with the given time
    const date = new Date()
    date.setHours(hours, minutes, 0, 0)

    // Format to 12-hour time with AM/PM
    const formattedTime = date.toLocaleTimeString('en-US', {
        hour: 'numeric',
        minute: '2-digit',
        hour12: true,
        timeZone: 'America/Los_Angeles'
    })

    return `${formattedTime} PT`
}

/**
 * Formats a date string to Pacific Time
 * @param {string} dateString - Date in YYYY-MM-DD format
 * @returns {string} Formatted date
 */
export function formatDateToPacific(dateString) {
    if (!dateString) return ''

    try {
        const date = parseISO(dateString)
        return format(date, 'MMM d, yyyy')
    } catch (error) {
        return dateString
    }
}

/**
 * Converts a date to YYYY-MM-DD format in Pacific timezone
 * @param {Date} date - JavaScript Date object
 * @returns {string} Date string in YYYY-MM-DD format
 */
export function formatDateToYYYYMMDD(date) {
    if (!(date instanceof Date) || isNaN(date.getTime())) {
        return ''
    }

    // Format in Pacific timezone
    const year = date.toLocaleString('en-US', { year: 'numeric', timeZone: 'America/Los_Angeles' })
    const month = date.toLocaleString('en-US', { month: '2-digit', timeZone: 'America/Los_Angeles' })
    const day = date.toLocaleString('en-US', { day: '2-digit', timeZone: 'America/Los_Angeles' })

    return `${year}-${month}-${day}`
}

/**
 * Filters showtimes by date
 * @param {Array} showtimes - Array of showtime objects
 * @param {Date} selectedDate - Selected date to filter by
 * @returns {Array} Filtered showtimes
 */
export function filterShowtimesByDate(showtimes, selectedDate) {
    if (!Array.isArray(showtimes) || !selectedDate) {
        return showtimes || []
    }

    const targetDate = formatDateToYYYYMMDD(selectedDate)

    return showtimes.filter((showtime) => {
        return showtime.date === targetDate
    })
}

/**
 * Filters showtimes by theater
 * @param {Array} showtimes - Array of showtime objects
 * @param {string} theaterId - Theater ID to filter by
 * @returns {Array} Filtered showtimes
 */
export function filterShowtimesByTheater(showtimes, theaterId) {
    if (!Array.isArray(showtimes) || !theaterId) {
        return showtimes || []
    }

    return showtimes.filter((showtime) => {
        return showtime.theater_id === theaterId
    })
}

/**
 * Gets unique movies that have showtimes on a given date
 * @param {Array} movies - Array of movie objects
 * @param {Array} showtimes - Array of showtime objects
 * @param {Date} selectedDate - Selected date
 * @returns {Array} Movies with showtimes on the selected date
 */
export function getMoviesWithShowtimes(movies, showtimes, selectedDate) {
    if (!Array.isArray(movies) || !Array.isArray(showtimes)) {
        return movies || []
    }

    // If no date selected, return all movies
    if (!selectedDate) {
        return movies
    }

    const filteredShowtimes = filterShowtimesByDate(showtimes, selectedDate)
    const movieIds = new Set(filteredShowtimes.map(st => st.tmdb_id))

    return movies.filter(movie => movieIds.has(movie.tmdb_id))
}
