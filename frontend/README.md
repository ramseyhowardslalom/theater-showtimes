# Portland Movie Theater Showtimes - Frontend

React frontend for the Portland Movie Theater Showtimes application with retro neon aesthetics.

## Features

- Browse movies and showtimes from Portland theaters
- Filter by theater, date, and genre
- View movie details with TMDB data
- Retro theater-themed UI with neon colors
- Responsive design for mobile and desktop

## Tech Stack

- React 18 with TypeScript
- React Router
- Axios for API calls
- Vite for build tooling
- Framer Motion for animations
- date-fns for date handling

## Setup

1. Install dependencies:
```bash
npm install
```

2. Configure environment variables:
```bash
cp .env.example .env
```

Edit `.env` to point to your backend API.

3. Start development server:
```bash
npm run dev
```

The app will be available at `http://localhost:3000`

## Development

### Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production (includes TypeScript compilation)
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint
- `npm run format` - Format code with Prettier
- `npm run type-check` - Check TypeScript types without emitting files
- `npm run format` - Format code with Prettier

### Project Structure

```
src/
├── components/          # Reusable UI components
│   ├── Header/
│   ├── MovieCard/
│   ├── TheaterList/
│   ├── FilterBar/
│   ├── DateSelector/
│   └── LoadingState/
├── pages/              # Page components
│   ├── Home.tsx
│   ├── Theater.tsx
│   └── Movie.tsx
├── services/           # API services
│   └── api.ts
├── styles/             # Global styles
│   ├── colors.css      # Color scheme
│   ├── neon.css        # Neon effects
│   └── global.css
├── types/              # TypeScript type definitions
│   └── index.ts
├── vite-env.d.ts       # Vite environment types
└── App.tsx

```

## Design System

### Colors
- Primary: Neon Pink (#FF1493), Neon Blue (#00D9FF)
- Accents: Neon Orange (#FF6B35), Neon Yellow (#FFD700)
- Backgrounds: Dark Navy (#0A0E27), Deep Purple (#16213E)

### Typography
- Headers: Bebas Neue
- Special text: Righteous
- Body: System fonts

## Code Quality

- TypeScript for type safety
- ESLint configured for React and TypeScript best practices
- Prettier for consistent formatting
- Pre-commit hooks recommended
- Strict mode enabled in tsconfig.json

## Building for Production

```bash
npm run build
```

The built files will be in the `dist/` directory.
