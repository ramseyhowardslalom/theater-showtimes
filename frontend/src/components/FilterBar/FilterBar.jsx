import './FilterBar.css'

function FilterBar({ theaters, selectedTheater, onTheaterChange, viewMode, onViewModeChange }) {
    return (
        <div className="filter-bar">
            <div className="view-mode-toggle">
                <button
                    className={viewMode === 'movies' ? 'active' : ''}
                    onClick={() => onViewModeChange('movies')}
                >
                    By Movie
                </button>
                <button
                    className={viewMode === 'theaters' ? 'active' : ''}
                    onClick={() => onViewModeChange('theaters')}
                >
                    By Theater
                </button>
            </div>

            <div className="theater-filter">
                <label htmlFor="theater-select">Theater:</label>
                <select
                    id="theater-select"
                    value={selectedTheater || ''}
                    onChange={(e) => onTheaterChange(e.target.value || null)}
                >
                    <option value="">All Theaters</option>
                    {theaters.map((theater) => (
                        <option key={theater.id} value={theater.id}>
                            {theater.name}
                        </option>
                    ))}
                </select>
            </div>
        </div>
    )
}

export default FilterBar
