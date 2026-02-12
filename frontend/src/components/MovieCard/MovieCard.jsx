import { useState } from 'react'
import { formatTimeToPacific } from '../../utils/dateTime'
import './MovieCard.css'

function MovieCard({ movie, showtimes, theaters }) {
    const [isExpanded, setIsExpanded] = useState(false)
    const movieShowtimes = showtimes.filter((st) => st.tmdb_id === movie.tmdb_id)

    // Get unique theater names for this movie
    const movieTheaters = [...new Set(movieShowtimes.map(st => st.theater_id))]
        .map(theaterId => theaters.find(t => t.id === theaterId))
        .filter(Boolean)

    // Get the first showtime link if available
    const ticketLink = movieShowtimes.find(st => st.link)?.link

    // Debug logging
    if (movieShowtimes.length > 0) {
        console.log(`Movie: ${movie.title}, First showtime:`, movieShowtimes[0], 'ticketLink:', ticketLink)
    }

    // Truncate overview if it's longer than 120 characters
    const shouldTruncate = movie.overview && movie.overview.length > 120
    const displayOverview = isExpanded || !shouldTruncate
        ? movie.overview
        : movie.overview?.substring(0, 120) + '...'

    return (
        <div className="movie-card">
            <div className="movie-poster">
                {movie.poster_path ? (
                    <img src={movie.poster_path} alt={movie.title} />
                ) : (
                    <div className="no-poster">No Poster</div>
                )}
            </div>
            <div className="movie-details">
                <h3 className="movie-title">{movie.title}</h3>
                <div className="movie-meta">
                    <span className="rating">{movie.rating}</span>
                    <span className="runtime">{movie.runtime} min</span>
                    <span className="tmdb-rating">⭐ {movie.tmdb_rating}</span>
                </div>
                <p className="movie-overview">
                    {displayOverview}
                    {shouldTruncate && (
                        <button
                            className="more-button"
                            onClick={() => setIsExpanded(!isExpanded)}
                        >
                            {isExpanded ? 'Less' : 'More'}
                        </button>
                    )}
                </p>
                <div className="showtimes-preview">
                    {movieShowtimes.slice(0, 3).map((st, idx) => (
                        <span key={idx} className="showtime-tag">
                            {formatTimeToPacific(st.time)}
                        </span>
                    ))}
                    {movieShowtimes.length > 3 && <span className="more">+{movieShowtimes.length - 3} more</span>}
                </div>
                <div className="theater-labels">
                    {movieTheaters.map((theater, idx) => (
                        <span key={idx} className="theater-label">
                            {theater.name}
                        </span>
                    ))}
                </div>
                {ticketLink && (
                    <div className="ticket-link-container">
                        <a href={ticketLink} target="_blank" rel="noopener noreferrer" className="ticket-link">
                            Tickets
                        </a>
                    </div>
                )}
            </div>
        </div>
    )
}

export default MovieCard
