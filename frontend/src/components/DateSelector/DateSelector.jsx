import { useState } from 'react'
import { format, addDays, startOfToday } from 'date-fns'
import './DateSelector.css'

function DateSelector({ selectedDate, onChange }) {
    const [weekOffset, setWeekOffset] = useState(0)
    const today = startOfToday()

    // Calculate the start of the current week view (today + weekOffset * 7 days)
    const weekStart = addDays(today, weekOffset * 7)
    const dates = Array.from({ length: 7 }, (_, i) => addDays(weekStart, i))

    const handlePrevious = () => {
        if (weekOffset > 0) {
            setWeekOffset(weekOffset - 1)
        }
    }

    const handleNext = () => {
        setWeekOffset(weekOffset + 1)
    }

    return (
        <div className="date-selector">
            <button
                className="nav-button"
                onClick={handlePrevious}
                disabled={weekOffset === 0}
            >
                ‹
            </button>
            {dates.map((date) => (
                <button
                    key={date.toISOString()}
                    className={`date-button ${format(selectedDate, 'yyyy-MM-dd') === format(date, 'yyyy-MM-dd') ? 'active' : ''}`}
                    onClick={() => onChange(date)}
                >
                    <span className="day">{format(date, 'EEE')}</span>
                    <span className="date">{format(date, 'MMM d')}</span>
                </button>
            ))}
            <button
                className="nav-button"
                onClick={handleNext}
            >
                ›
            </button>
        </div>
    )
}

export default DateSelector
