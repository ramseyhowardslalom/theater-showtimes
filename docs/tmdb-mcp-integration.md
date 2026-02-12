# TMDB MCP Server Integration

This document explains how the TMDB MCP (Model Context Protocol) server is integrated into the Portland Movie Theater Showtimes project.

## Overview

The TMDB MCP server provides access to The Movie Database (TMDB) API through the Model Context Protocol, enabling AI assistants (like GitHub Copilot) to search for movies, get recommendations, and retrieve detailed movie information.

## Setup

The TMDB MCP server is located in `mcp-servers/mcp-server-tmdb/` and is configured in `.vscode/mcp.json`.

### Installation

The server has already been set up with the following steps:

1. **Cloned the repository**:
   ```bash
   mkdir -p mcp-servers
   cd mcp-servers
   git clone https://github.com/Laksh-star/mcp-server-tmdb.git
   ```

2. **Installed dependencies and built**:
   ```bash
   cd mcp-server-tmdb
   npm install
   # Build automatically runs via prepare script
   ```

3. **Configured in VS Code**:
   The server is registered in `.vscode/mcp.json` with the TMDB API key.

### Configuration

The MCP server is configured in `.vscode/mcp.json`:

```json
{
  "servers": {
    "tmdb": {
      "command": "node",
      "args": ["/private/tmp/theater-showtimes/mcp-servers/mcp-server-tmdb/dist/index.js"],
      "env": {
        "TMDB_API_KEY": "your_api_key_here"
      }
    }
  }
}
```

**Note**: The TMDB API key is stored in the configuration. Keep this file secure and do not commit it to version control.

## Available Tools

The TMDB MCP server provides three tools:

### 1. search_movies

Search for movies by title or keywords.

**Input**:
- `query` (string): Search query for movie titles

**Returns**: List of movies with:
- Title and release year
- Movie ID
- Rating (out of 10)
- Overview/description

**Example usage**:
```
"Search for movies about space exploration"
"Find movies with 'Batman' in the title"
```

### 2. get_recommendations

Get movie recommendations based on a specific movie.

**Input**:
- `movieId` (string): TMDB movie ID

**Returns**: Top 5 recommended movies with details

**Example usage**:
```
"Get movie recommendations based on movie ID 550" (Fight Club)
"What movies are similar to movie 680?" (Pulp Fiction)
```

### 3. get_trending

Get trending movies for a specified time window.

**Input**:
- `timeWindow` (string): Either "day" or "week"

**Returns**: Top 10 trending movies with details

**Example usage**:
```
"What are the trending movies today?"
"Show me this week's trending movies"
```

## Available Resources

The server also provides access to detailed movie information via resources:

**Resource URI Format**: `tmdb:///movie/<movie_id>`

**Returns**: Comprehensive movie details including:
- Title and release date
- Rating and overview
- Genres
- Poster URL
- Cast information (top 5 actors)
- Director
- Selected reviews (top 3)

## Using the MCP Server

### In GitHub Copilot Chat

When using GitHub Copilot in VS Code, you can interact with the TMDB server through natural language:

1. **Search for movies**:
   - "Search for movies about artificial intelligence"
   - "Find movies released in 2023"

2. **Get trending movies**:
   - "What are today's trending movies?"
   - "Show me this week's popular movies"

3. **Get recommendations**:
   - "Recommend movies similar to The Matrix"
   - First, search for the movie to get its ID, then use that ID for recommendations

4. **Get movie details**:
   - "Tell me about movie ID 550"
   - "Get details for The Shawshank Redemption"

### Integration with Backend

The backend Go code in `backend/internal/tmdb/client.go` can be enhanced to use the MCP server for movie data enrichment. Currently, it has placeholder HTTP client code, but it can be updated to leverage the MCP server when running in a development environment with GitHub Copilot.

For production use, the backend should use direct TMDB API calls via HTTP.

## API Key Management

### Getting a TMDB API Key

1. Create a free account at [TMDB](https://www.themoviedb.org/)
2. Go to your account settings
3. Navigate to the API section
4. Request an API key for developer use
5. Wait for approval (usually quick for developer accounts)

### Security

- **Never commit API keys** to version control
- The `.vscode/` directory should be added to `.gitignore` if it contains sensitive credentials
- Consider using environment variables for production deployments
- For team collaboration, document in the README that developers need to obtain their own TMDB API key

## Troubleshooting

### Server not responding

1. Check that Node.js is installed: `node --version`
2. Verify the server was built: Check for `dist/index.js` in `mcp-servers/mcp-server-tmdb/`
3. Rebuild if needed: `cd mcp-servers/mcp-server-tmdb && npm run build`

### Invalid API key errors

1. Verify your TMDB API key is correct in `.vscode/mcp.json`
2. Check that your TMDB account has API access enabled
3. Ensure the API key is properly approved by TMDB

### Path issues

The configuration uses an absolute path. If you move the project, update the path in `.vscode/mcp.json`:

```json
"args": ["/full/path/to/theater-showtimes/mcp-servers/mcp-server-tmdb/dist/index.js"]
```

## Further Reading

- [TMDB API Documentation](https://developers.themoviedb.org/3)
- [Model Context Protocol Specification](https://modelcontextprotocol.io/)
- [Original mcp-server-tmdb Repository](https://github.com/Laksh-star/mcp-server-tmdb)
