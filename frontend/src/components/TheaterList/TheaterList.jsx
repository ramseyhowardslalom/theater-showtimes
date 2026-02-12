import { Link } from 'react-router-dom'
import './TheaterList.css'

function TheaterList({ theaters, showtimes }) {
    return (
        <div className="theater-list">
            {theaters.map((theater) => {
                const theaterShowtimes = showtimes.filter((st) => st.theater_id === theater.id)

                return (
                    <div key={theater.id} className="theater-item">
                        <div className="theater-header">
                            <h2>{theater.name}</h2>
                            <p className="theater-address">
                                {theater.address}, {theater.city}
                            </p>
                        </div>
                        <div className="theater-showtimes">
                            {theaterShowtimes.length > 0 ? (
                                <p>{theaterShowtimes.length} showtimes available</p>
                            ) : (
                                <p className="no-showtimes">No showtimes available</p>
                            )}
                            <Link to={`/theater/${theater.id}`} className="view-link">
                                View All Showtimes →
                            </Link>
                        </div>
                    </div>
                )
            })}
        </div>
    )
}

export default TheaterList
