// Common types for the Theater Showtimes application

export interface Theater {
    id: string
    name: string
    address: string
    city: string
    zip: string
    website: string
    phone?: string
}

export interface Movie {
    tmdb_id: number
    title: string
    original_title: string
    overview: string
    runtime: number
    rating: string
    genres: string[]
    release_date: string
    poster_path: string
    backdrop_path: string
    tmdb_rating: number
    vote_count: number
    popularity: number
    cast: string[]
    director: string
}

export interface Showtime {
    id: string
    theater_id: string
    movie_title: string
    tmdb_id: number
    date: string
    time: string
    format: string
    price?: number
    booking_url?: string
    screen?: string
}

export interface ScrapeMetadata {
    last_updated: string
    theater_id: string
    status: 'success' | 'error'
    error_message?: string
    movies_scraped: number
    showtimes_scraped: number
}

export type ViewMode = 'movies' | 'theaters'

export interface ApiError {
    message: string
    code?: string
}

export interface ApiResponse<T> {
    data?: T
    error?: ApiError
}
