import React, { useState, useEffect, useRef } from 'react'
import { connectWebSocket, closeSession } from '../utils/api'
import './ChatWindow.css'

export default function ChatWindow({ sessionId, userId, onSessionClosed }) {
  const [messages, setMessages] = useState([])
  const [inputValue, setInputValue] = useState('')
  const [loading, setLoading] = useState(false)
  const [connected, setConnected] = useState(false)
  const ws = useRef(null)
  const messagesEndRef = useRef(null)

  useEffect(() => {
    // Connect to WebSocket
    ws.current = connectWebSocket(sessionId, userId)

    ws.current.onopen = () => {
      setConnected(true)
      console.log('WebSocket connected')
    }

    ws.current.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data)
        setMessages((prev) => [...prev, message])
        setLoading(false)
      } catch (err) {
        console.error('Failed to parse message:', err)
      }
    }

    ws.current.onerror = (error) => {
      console.error('WebSocket error:', error)
      setConnected(false)
    }

    ws.current.onclose = () => {
      setConnected(false)
      console.log('WebSocket disconnected')
    }

    return () => {
      if (ws.current && ws.current.readyState === WebSocket.OPEN) {
        ws.current.close()
      }
    }
  }, [sessionId, userId])

  // Auto-scroll to bottom
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  const handleSendMessage = (e) => {
    e.preventDefault()

    if (!inputValue.trim() || !connected) {
      return
    }

    const message = {
      type: 'message',
      content: inputValue,
    }

    // Add user message to display
    setMessages((prev) => [
      ...prev,
      {
        type: 'message',
        content: inputValue,
        user_id: userId,
      },
    ])

    // Send to WebSocket
    ws.current.send(JSON.stringify(message))

    setInputValue('')
    setLoading(true)
  }

  const handleCloseSession = async () => {
    try {
      await closeSession(sessionId)
      if (ws.current) {
        ws.current.close()
      }
      onSessionClosed()
    } catch (err) {
      console.error('Failed to close session:', err)
    }
  }

  return (
    <div className="chat-window">
      <div className="chat-header">
        <div className="header-content">
          <h1>MedSeek Doctor Consultation</h1>
          <p className="user-info">Session with {userId}</p>
        </div>
        <div className="header-status">
          <span className={`status-indicator ${connected ? 'connected' : 'disconnected'}`}>
            {connected ? '● Connected' : '● Disconnected'}
          </span>
          <button onClick={handleCloseSession} className="close-button">
            End Session
          </button>
        </div>
      </div>

      <div className="messages-container">
        {messages.length === 0 && (
          <div className="welcome-message">
            <h2>Welcome to MedSeek</h2>
            <p>Start your consultation by describing your health concerns below.</p>
            <div className="guidelines">
              <h3>Guidelines:</h3>
              <ul>
                <li>Describe your symptoms in detail</li>
                <li>Mention any relevant medical history</li>
                <li>Be honest about your health concerns</li>
                <li>Remember: This is an AI assistant, not a substitute for professional medical advice</li>
              </ul>
            </div>
          </div>
        )}

        {messages.map((msg, idx) => (
          <div
            key={idx}
            className={`message ${msg.user_id === 'assistant' ? 'assistant-message' : 'user-message'}`}
          >
            <div className="message-content">
              {msg.user_id === 'assistant' && <span className="message-role">🏥 Doctor</span>}
              {msg.user_id !== 'assistant' && <span className="message-role">👤 You</span>}
              <p>{msg.content}</p>
            </div>
          </div>
        ))}

        {loading && (
          <div className="message assistant-message">
            <div className="message-content">
              <span className="message-role">🏥 Doctor</span>
              <p className="typing-indicator">
                <span></span>
                <span></span>
                <span></span>
              </p>
            </div>
          </div>
        )}

        <div ref={messagesEndRef} />
      </div>

      <form onSubmit={handleSendMessage} className="message-form">
        <input
          type="text"
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          placeholder="Describe your symptoms or ask a medical question..."
          disabled={!connected || loading}
          className="message-input"
        />
        <button type="submit" disabled={!connected || loading} className="send-button">
          {loading ? 'Sending...' : 'Send'}
        </button>
      </form>
    </div>
  )
}
