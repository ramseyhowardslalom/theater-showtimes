import { Link } from 'react-router-dom'
import './Header.css'

function Header() {
    return (
        <header className="header">
            <div className="header-content">
                <Link to="/" className="logo">
                    <h1>
                        <span className="neon-text">Portland Movie Theater</span>
                        <span className="neon-text-alt">Showtimes</span>
                    </h1>
                </Link>
            </div>
        </header>
    )
}

export default Header
