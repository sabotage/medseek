import './App.css'
import React, { useState, useEffect, useRef } from 'react'
import ChatWindow from './components/ChatWindow'
import SessionSetup from './components/SessionSetup'

function App() {
  const [sessionId, setSessionId] = useState(null)
  const [userId, setUserId] = useState(null)
  const [connected, setConnected] = useState(false)

  const handleSessionCreated = (newSessionId, newUserId) => {
    setSessionId(newSessionId)
    setUserId(newUserId)
    setConnected(true)
  }

  const handleSessionClosed = () => {
    setSessionId(null)
    setUserId(null)
    setConnected(false)
  }

  return (
    <div className="app">
      {!connected ? (
        <SessionSetup onSessionCreated={handleSessionCreated} />
      ) : (
        <ChatWindow
          sessionId={sessionId}
          userId={userId}
          onSessionClosed={handleSessionClosed}
        />
      )}
    </div>
  )
}

export default App
