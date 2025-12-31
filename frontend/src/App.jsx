import './App.css'
import React, { useState, useEffect, useRef } from 'react'
import ChatWindow from './components/ChatWindow'
import SessionSetup from './components/SessionSetup'
import { initializeIOSFixes } from './utils/iosHelper'

function App() {
  const [sessionId, setSessionId] = useState(null)
  const [userId, setUserId] = useState(null)
  const [specialty, setSpecialty] = useState(null)
  const [connected, setConnected] = useState(false)

  // 初始化iOS Safari修复
  useEffect(() => {
    initializeIOSFixes()
  }, [])

  const handleSessionCreated = (newSessionId, newUserId, newSpecialty) => {
    setSessionId(newSessionId)
    setUserId(newUserId)
    setSpecialty(newSpecialty)
    setConnected(true)
  }

  const handleSessionClosed = () => {
    setSessionId(null)
    setUserId(null)
    setSpecialty(null)
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
          specialty={specialty}
          onSessionClosed={handleSessionClosed}
        />
      )}
    </div>
  )
}

export default App
