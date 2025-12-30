# MedSeek - Online Doctor Consultation Platform

A web-based chat application that provides online doctor consultation services powered by DeepSeek AI.

## Features

- 💬 **Real-time Chat Interface** - WebSocket-based instant messaging
- 🏥 **AI Doctor Consultation** - Powered by DeepSeek medical model
- 🔒 **Secure Sessions** - Individual chat sessions for each consultation
- 📱 **Responsive Design** - Works on desktop, tablet, and mobile devices
- 🎨 **Modern UI** - Beautiful gradient design with smooth animations
- 📊 **Chat History** - Persistent message history for each session

## Tech Stack

### Backend
- **Go 1.22.2** - Server runtime
- **Gorilla WebSocket** - Real-time WebSocket communication
- **DeepSeek API** - AI-powered medical consultation

### Frontend
- **React 18** - UI framework
- **Vite** - Build tool and dev server
- **Axios** - HTTP client
- **CSS3** - Modern styling with animations

## Project Structure

```
medseek/
├── cmd/
│   └── server/
│       └── main.go              # Server entry point
├── internal/
│   ├── models/
│   │   └── models.go            # Data models
│   ├── service/
│   │   └── chat_service.go      # Chat business logic
│   ├── handlers/
│   │   └── handlers.go          # HTTP/WebSocket handlers
│   ├── deepseek/
│   │   └── client.go            # DeepSeek API client
│   ├── websocket/
│   │   └── hub.go               # WebSocket hub
│   └── db/
│       └── (database layer - placeholder)
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── ChatWindow.jsx   # Main chat component
│   │   │   └── SessionSetup.jsx # Session initialization
│   │   ├── utils/
│   │   │   └── api.js           # API client
│   │   ├── styles/
│   │   │   └── (CSS files)
│   │   ├── App.jsx              # Root component
│   │   └── main.jsx             # Entry point
│   ├── index.html               # HTML template
│   ├── vite.config.js           # Vite configuration
│   └── package.json             # Frontend dependencies
├── go.mod                       # Go module
├── .env.example                 # Environment template
└── README.md                    # This file
```

## Prerequisites

- Go 1.22.2 or higher
- Node.js 16+ and npm
- DeepSeek API key (get one at https://www.deepseek.com)

## Installation

### 1. Clone and Setup

```bash
cd /home/oliver/projects/medseek
```

### 2. Backend Setup

```bash
# Download Go dependencies
go mod download

# Create .env file with your DeepSeek API key
cp .env.example .env
# Edit .env and add your actual DeepSeek API key
nano .env
```

### 3. Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Build frontend
npm run build

cd ..
```

## Running the Application

### Development Mode

**Terminal 1 - Frontend Dev Server:**
```bash
cd frontend
npm run dev
```
Frontend will be available at `http://localhost:3000`

**Terminal 2 - Backend Server:**
```bash
go run cmd/server/main.go
```
Backend will be running at `http://localhost:8080`

### Production Mode

```bash
# Build frontend (if not already built)
cd frontend
npm run build
cd ..

# Run server (frontend will be served from /frontend/dist)
go run cmd/server/main.go
```

Visit `http://localhost:8080` in your browser.

## API Endpoints

### REST API

- `POST /api/session/create` - Create a new chat session
  - Request: `{ "user_id": "user@example.com" }`
  - Response: `{ "session_id": "xxx", "status": "active" }`

- `GET /api/session/messages` - Get session messages
  - Query: `?session_id=xxx`
  - Response: Array of messages

- `POST /api/session/close` - Close a session
  - Query: `?session_id=xxx`
  - Response: `{ "status": "closed" }`

- `GET /health` - Health check
  - Response: `{ "status": "ok" }`

### WebSocket

- `WS /ws?session_id=xxx&user_id=yyy` - Real-time chat connection
  - Message format: `{ "type": "message", "content": "..." }`

## Environment Variables

Create a `.env` file in the root directory:

```
DEEPSEEK_API_KEY=your_actual_api_key_here
PORT=8080
```

## How to Use

1. **Start a Session**
   - Enter your name and email
   - Click "Start Consultation"

2. **Chat with Doctor**
   - Describe your symptoms or health concerns
   - Get instant AI-powered medical guidance
   - Ask follow-up questions

3. **End Session**
   - Click "End Session" button
   - Session will be closed and archived

## System Prompt

The AI doctor operates with the following guidelines:
- Provides general medical information and guidance
- Asks clarifying questions about symptoms
- Recommends professional medical care when needed
- Maintains patient privacy and confidentiality
- Provides evidence-based medical information

## Security Considerations

- ⚠️ This is an AI assistant, not a substitute for professional medical advice
- Emergency symptoms should trigger recommendations for immediate professional care
- All data should be encrypted in production
- Implement proper authentication for production
- Store conversations securely in a database

## Future Enhancements

- [ ] Database integration (PostgreSQL)
- [ ] User authentication and authorization
- [ ] Doctor profiles and ratings
- [ ] Prescription management
- [ ] Medical report generation
- [ ] Video consultation feature
- [ ] Payment integration
- [ ] Multi-language support
- [ ] Mobile app (React Native)
- [ ] Admin dashboard

## Troubleshooting

### WebSocket Connection Issues
- Ensure backend server is running on port 8080
- Check browser console for connection errors
- Verify CORS settings if using different domains

### DeepSeek API Errors
- Verify your API key is correct in `.env`
- Check API rate limits
- Ensure API key has necessary permissions

### Frontend Build Issues
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
npm run build
```

## Contributing

Feel free to fork and submit pull requests!

## License

MIT License - feel free to use this project for your own purposes.

## Support

For issues and questions:
1. Check the troubleshooting section
2. Review DeepSeek API documentation
3. Open an issue on GitHub

---

**Built with ❤️ for healthcare accessibility**
