import { useParams } from 'react-router-dom'
import { useState, useEffect } from 'react'
import { getMovieDetails } from '../services/api'
import './Movie.css'

function Movie() {
    const { id } = useParams()
    const [movie, setMovie] = useState(null)
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        loadMovieDetails()
    }, [id])

    const loadMovieDetails = async () => {
        try {
            setLoading(true)
            const data = await getMovieDetails(id)
            setMovie(data)
        } catch (err) {
            console.error(err)
        } finally {
            setLoading(false)
        }
    }

    if (loading) return <div>Loading...</div>
    if (!movie) return <div>Movie not found</div>

    return (
        <div className="movie-page">
            <div className="movie-header">
                <img src={movie.poster_path} alt={movie.title} />
                <div className="movie-info">
                    <h1>{movie.title}</h1>
                    <p className="overview">{movie.overview}</p>
                    <div className="movie-meta">
                        <span>Rating: {movie.rating}</span>
                        <span>Runtime: {movie.runtime} min</span>
                        <span>TMDB: {movie.tmdb_rating}/10</span>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default Movie
