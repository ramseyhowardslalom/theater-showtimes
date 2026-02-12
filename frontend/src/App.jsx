import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import Home from './pages/Home'
import Theater from './pages/Theater'
import Movie from './pages/Movie'
import Header from './components/Header/Header'
import './styles/colors.css'
import './styles/neon.css'

function App() {
    return (
        <Router>
            <div className="app">
                <Header />
                <main className="main-content">
                    <Routes>
                        <Route path="/" element={<Home />} />
                        <Route path="/theater/:id" element={<Theater />} />
                        <Route path="/movie/:id" element={<Movie />} />
                    </Routes>
                </main>
            </div>
        </Router>
    )
}

export default App
