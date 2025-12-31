import React, { useState, useEffect, useRef } from 'react'
import { connectWebSocket, closeSession } from '../utils/api'
import { getSpecialtyInfo } from '../utils/specialties'
import { scrollToBottom, onIOSKeyboardToggle, isIOSSafari } from '../utils/iosHelper'
import './ChatWindow.css'

export default function ChatWindow({ sessionId, userId, specialty, onSessionClosed }) {
  const [messages, setMessages] = useState([])
  const [inputValue, setInputValue] = useState('')
  const [loading, setLoading] = useState(false)
  const [connected, setConnected] = useState(false)
  const ws = useRef(null)
  const messagesEndRef = useRef(null)
  const inputRef = useRef(null)
  const formRef = useRef(null)
  
  const info = getSpecialtyInfo(specialty)

  useEffect(() => {
    // 处理iOS键盘弹出
    if (isIOSSafari()) {
      const unsubscribe = onIOSKeyboardToggle((state) => {
        setTimeout(() => {
          const container = document.querySelector('.messages-container')
          if (container) {
            scrollToBottom(container, false)
          }
        }, 100)
      })
      return () => unsubscribe()
    }
  }, [])

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
    const container = document.querySelector('.messages-container')
    scrollToBottom(container, true)
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

    // Send to WebSocket (server will echo it back)
    ws.current.send(JSON.stringify(message))

    setInputValue('')
    setLoading(true)

    // iOS修复：保持焦点在输入框上
    if (inputRef.current && isIOSSafari()) {
      setTimeout(() => {
        inputRef.current?.focus()
        inputRef.current?.scrollIntoView({ behavior: 'smooth', block: 'center' })
      }, 100)
    }
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
          <h1>{info.emoji} {info.title}</h1>
          <p className="user-info">当前咨询：{userId}</p>
        </div>
        <div className="header-status">
          <span className={`status-indicator ${connected ? 'connected' : 'disconnected'}`}>
            {connected ? '● 在线' : '● 已断开'}
          </span>
          <button onClick={handleCloseSession} className="close-button">
            结束咨询
          </button>
        </div>
      </div>

      <div className="messages-container">
        {messages.length === 0 && (
          <div className="welcome-message">
            <h2>{info.emoji} {info.welcome}</h2>
            <p>{info.description}</p>
            <div className="guidelines">
              <h3>咨询须知：</h3>
              <ul>
                {info.guidelines.map((guideline, idx) => (
                  <li key={idx}>{guideline}</li>
                ))}
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
              {msg.user_id === 'assistant' && <span className="message-role">{info.emoji} {info.doctorRole}</span>}
              {msg.user_id !== 'assistant' && <span className="message-role">👤 患者</span>}
              <p>{msg.content}</p>
            </div>
          </div>
        ))}

        {loading && (
          <div className="message assistant-message">
            <div className="message-content">
              <span className="message-role">{info.emoji} 医生正在分析...</span>
              <div className="thinking-indicator">
                <div className="pulse-dot"></div>
                <div className="pulse-dot"></div>
                <div className="pulse-dot"></div>
              </div>
              <p className="thinking-text">医生正在根据您的症状进行分析，请稍候...</p>
            </div>
          </div>
        )}
      </div>

      <form onSubmit={handleSendMessage} className="message-form" ref={formRef}>
        <input
          ref={inputRef}
          type="text"
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          onFocus={() => {
            // iOS修复：处理焦点时的滚动
            if (isIOSSafari()) {
              setTimeout(() => {
                const container = document.querySelector('.messages-container')
                if (container) {
                  scrollToBottom(container, false)
                }
              }, 300)
            }
          }}
          placeholder="请描述您的症状或咨询的问题..."
          disabled={!connected || loading}
          className="message-input"
        />
        <button type="submit" disabled={!connected || loading} className="send-button">
          {loading ? '发送中...' : '发送'}
        </button>
      </form>
    </div>
  )
}
