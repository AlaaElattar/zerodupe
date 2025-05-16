# Server Components

The server components follow a clean architecture with clear separation of concerns:

## Architecture

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│             │     │             │     │             │
│     API     │────▶│   Handler   │────▶│   Storage   │
│             │     │             │     │             │
└─────────────┘     └─────────────┘     └─────────────┘
```

### Components

1. **API**: HTTP routes and server setup
2. **Handler**: Contains logic for handling requests
3. **Storage**: Interface for data persistence
4. **Model**: Data structures shared across components
5. **Config**: Server configuration