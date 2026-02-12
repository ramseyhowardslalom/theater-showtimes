import { useParams } from 'react-router-dom'
import { useState, useEffect } from 'react'
import { getTheaterShowtimes } from '../services/api'
import './Theater.css'

function Theater() {
    const { id } = useParams()
    const [showtimes, setShowtimes] = useState([])
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        loadTheaterShowtimes()
    }, [id])

    const loadTheaterShowtimes = async () => {
        try {
            setLoading(true)
            const data = await getTheaterShowtimes(id)
            setShowtimes(data)
        } catch (err) {
            console.error(err)
        } finally {
            setLoading(false)
        }
    }

    if (loading) return <div>Loading...</div>

    return (
        <div className="theater-page">
            <h1>Theater Showtimes</h1>
            <div className="showtimes-list">
                {showtimes.map((showtime) => (
                    <div key={showtime.id} className="showtime-item">
                        <h3>{showtime.movie_title}</h3>
                        <p>{showtime.date} - {showtime.time}</p>
                        <p>{showtime.format}</p>
                    </div>
                ))}
            </div>
        </div>
    )
}

export default Theater
