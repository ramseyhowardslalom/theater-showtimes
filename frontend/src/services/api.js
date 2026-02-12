import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api'

const api = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
})

// Theaters
export const getTheaters = async () => {
    const response = await api.get('/theaters')
    return response.data
}

// Showtimes
export const getShowtimes = async (params = {}) => {
    const response = await api.get('/showtimes', { params })
    return response.data
}

export const getTheaterShowtimes = async (theaterId) => {
    const response = await api.get(`/showtimes/${theaterId}`)
    return response.data
}

// Movies
export const getMovies = async () => {
    const response = await api.get('/movies')
    return response.data
}

export const getMovieDetails = async (movieId) => {
    const response = await api.get(`/movies/${movieId}`)
    return response.data
}

// Scraper
export const triggerScrape = async (theaterIds = []) => {
    const response = await api.post('/scrape', { theater_ids: theaterIds })
    return response.data
}

// Health & Meta
export const getHealth = async () => {
    const response = await api.get('/health')
    return response.data
}

export const getLastUpdated = async () => {
    const response = await api.get('/last-updated')
    return response.data
}

export default api
