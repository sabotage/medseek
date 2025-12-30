import React, { useState } from 'react'
import { createSession } from '../utils/api'
import './SessionSetup.css'

export default function SessionSetup({ onSessionCreated }) {
  const [userEmail, setUserEmail] = useState('')
  const [userName, setUserName] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      if (!userEmail || !userName) {
        throw new Error('Please fill in all fields')
      }

      const { session_id } = await createSession(userEmail)
      onSessionCreated(session_id, userEmail)
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="session-setup">
      <div className="setup-container">
        <div className="setup-card">
          <h1>MedSeek</h1>
          <p className="subtitle">Online Doctor Consultation</p>

          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label htmlFor="name">Your Name</label>
              <input
                type="text"
                id="name"
                value={userName}
                onChange={(e) => setUserName(e.target.value)}
                placeholder="Enter your full name"
                disabled={loading}
              />
            </div>

            <div className="form-group">
              <label htmlFor="email">Email Address</label>
              <input
                type="email"
                id="email"
                value={userEmail}
                onChange={(e) => setUserEmail(e.target.value)}
                placeholder="Enter your email"
                disabled={loading}
              />
            </div>

            {error && <div className="error-message">{error}</div>}

            <button
              type="submit"
              disabled={loading}
              className="start-button"
            >
              {loading ? 'Creating Session...' : 'Start Consultation'}
            </button>
          </form>

          <div className="info-box">
            <h3>How It Works</h3>
            <ul>
              <li>Chat with our AI doctor for initial consultation</li>
              <li>Get medical guidance based on your symptoms</li>
              <li>Receive recommendations for professional care if needed</li>
              <li>All conversations are private and secure</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  )
}
