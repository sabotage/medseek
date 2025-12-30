# MedSeek Development Setup Guide

## Quick Start

### Prerequisites
- Go 1.22.2+
- Node.js 16+
- DeepSeek API key

### Setup Steps

1. **Environment Configuration**
   ```bash
   cp .env.example .env
   # Edit .env with your DeepSeek API key
   ```

2. **Backend Development**
   ```bash
   go mod download
   go run cmd/server/main.go
   ```

3. **Frontend Development**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

## Project Architecture

### Backend (Go)
- Models: Data structures for users, messages, sessions
- Service: Chat business logic and message processing
- Handlers: HTTP and WebSocket request handlers
- DeepSeek Client: API integration for AI responses
- WebSocket Hub: Real-time message broadcasting

### Frontend (React)
- SessionSetup: User authentication and session creation
- ChatWindow: Main chat interface
- API Utils: Backend communication

## Environment Variables

```env
DEEPSEEK_API_KEY=your_api_key_here
PORT=8080
```

## API Flow

1. User creates session → POST /api/session/create
2. User connects WebSocket → WS /ws?session_id=xxx&user_id=yyy
3. User sends message → WebSocket message
4. Message processed by chat service
5. Response from DeepSeek API
6. Response sent back via WebSocket
7. User closes session → POST /api/session/close

## WebSocket Message Format

### User Message
```json
{
  "type": "message",
  "content": "I have a headache"
}
```

### Assistant Response
```json
{
  "type": "message",
  "content": "I understand you have a headache. Can you tell me more about the severity and duration?",
  "user_id": "assistant"
}
```

## Testing

### Test Backend
```bash
# Health check
curl http://localhost:8080/health

# Create session
curl -X POST http://localhost:8080/api/session/create \
  -H "Content-Type: application/json" \
  -d '{"user_id":"test@example.com"}'

# Get messages
curl "http://localhost:8080/api/session/messages?session_id=xxx"
```

### Test Frontend
```bash
cd frontend
npm run dev
# Open http://localhost:3000
```

## Build for Production

### Frontend
```bash
cd frontend
npm run build
# Builds to frontend/dist
```

### Backend
```bash
# Go binaries are platform-specific
# Linux
GOOS=linux GOARCH=amd64 go build -o medseek cmd/server/main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o medseek cmd/server/main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o medseek.exe cmd/server/main.go
```

## Key Features Implemented

✅ Real-time WebSocket communication
✅ DeepSeek API integration
✅ Medical consultation prompting
✅ Session management
✅ Beautiful responsive UI
✅ Message history
✅ Error handling

## Next Steps

1. Implement PostgreSQL database
2. Add user authentication
3. Add doctor profiles
4. Implement encryption
5. Add mobile support
6. Deploy to production
