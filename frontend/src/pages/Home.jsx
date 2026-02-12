import { useState, useEffect, useMemo } from 'react'
import { getShowtimes, getMovies, getTheaters } from '../services/api'
import MovieCard from '../components/MovieCard/MovieCard'
import FilterBar from '../components/FilterBar/FilterBar'
import DateSelector from '../components/DateSelector/DateSelector'
import TheaterList from '../components/TheaterList/TheaterList'
import LoadingState from '../components/LoadingState/LoadingState'
import { filterShowtimesByDate, filterShowtimesByTheater, getMoviesWithShowtimes } from '../utils/dateTime'
import './Home.css'

function Home() {
    const [showtimes, setShowtimes] = useState([])
    const [movies, setMovies] = useState([])
    const [theaters, setTheaters] = useState([])
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState(null)
    const [selectedDate, setSelectedDate] = useState(new Date())
    const [selectedTheater, setSelectedTheater] = useState(null)
    const [viewMode, setViewMode] = useState('movies') // 'movies' or 'theaters'

    useEffect(() => {
        loadData()
    }, [])

    const loadData = async () => {
        try {
            setLoading(true)
            const [showtimesData, moviesData, theatersData] = await Promise.all([
                getShowtimes(),
                getMovies(),
                getTheaters(),
            ])
            setShowtimes(showtimesData)
            setMovies(moviesData)
            setTheaters(theatersData)
            setError(null)
        } catch (err) {
            setError('Failed to load data. Please try again.')
            console.error(err)
        } finally {
            setLoading(false)
        }
    }

    // Filter showtimes by selected date and theater
    const filteredShowtimes = useMemo(() => {
        let filtered = showtimes

        // Filter by date
        if (selectedDate) {
            filtered = filterShowtimesByDate(filtered, selectedDate)
        }

        // Filter by theater
        if (selectedTheater) {
            filtered = filterShowtimesByTheater(filtered, selectedTheater)
        }

        return filtered
    }, [showtimes, selectedDate, selectedTheater])

    // Get movies that have showtimes for the filtered criteria
    const filteredMovies = useMemo(() => {
        if (!selectedDate && !selectedTheater) {
            return movies
        }

        // Get unique movie IDs from filtered showtimes
        const movieIds = new Set(filteredShowtimes.map(st => st.tmdb_id))
        return movies.filter(movie => movieIds.has(movie.tmdb_id))
    }, [movies, filteredShowtimes, selectedDate, selectedTheater])

    if (loading) return <LoadingState />
    if (error) return <div className="error-state">{error}</div>

    return (
        <div className="home">
            <div className="controls">
                <DateSelector selectedDate={selectedDate} onChange={setSelectedDate} />
                <FilterBar
                    theaters={theaters}
                    selectedTheater={selectedTheater}
                    onTheaterChange={setSelectedTheater}
                    viewMode={viewMode}
                    onViewModeChange={setViewMode}
                />
            </div>

            {viewMode === 'movies' ? (
                <div className="movies-grid">
                    {filteredMovies.map((movie) => (
                        <MovieCard
                            key={movie.tmdb_id}
                            movie={movie}
                            showtimes={filteredShowtimes} theaters={theaters} selectedDate={selectedDate}
                        />
                    ))}
                </div>
            ) : (
                <TheaterList theaters={theaters} showtimes={filteredShowtimes} />
            )}
        </div>
    )
}

export default Home
