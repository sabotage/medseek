import axios from 'axios'

const API_BASE_URL = '/api'

export const createSession = async (userId, specialty = 'obstetrics') => {
  try {
    const response = await axios.post(`${API_BASE_URL}/session/create`, {
      user_id: userId,
      specialty: specialty,
    })
    return response.data
  } catch (error) {
    throw new Error(`Failed to create session: ${error.message}`)
  }
}

export const getSessionMessages = async (sessionId) => {
  try {
    const response = await axios.get(`${API_BASE_URL}/session/messages`, {
      params: { session_id: sessionId },
    })
    return response.data || []
  } catch (error) {
    throw new Error(`Failed to fetch messages: ${error.message}`)
  }
}

export const closeSession = async (sessionId) => {
  try {
    const response = await axios.post(`${API_BASE_URL}/session/close`, null, {
      params: { session_id: sessionId },
    })
    return response.data
  } catch (error) {
    throw new Error(`Failed to close session: ${error.message}`)
  }
}

export const connectWebSocket = (sessionId, userId) => {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const wsUrl = `${protocol}//${window.location.host}/ws?session_id=${sessionId}&user_id=${userId}`
  
  const ws = new WebSocket(wsUrl)
  
  // iOS Safari specific handling
  if (/iPad|iPhone|iPod/.test(navigator.userAgent)) {
    // Set a longer timeout for iOS Safari
    ws._iosReconnectAttempts = 0
    ws._iosMaxReconnectAttempts = 5
    
    ws.addEventListener('close', function handleIOSClose(event) {
      // Auto-reconnect for iOS Safari (except on normal close)
      if (!event.wasClean && this._iosReconnectAttempts < this._iosMaxReconnectAttempts) {
        const delay = Math.min(1000 * Math.pow(2, this._iosReconnectAttempts), 10000)
        console.log(`iOS WebSocket disconnected, attempting reconnect in ${delay}ms...`)
        
        setTimeout(() => {
          // Note: This is a simple retry, full reconnection logic should be in the component
          this._iosReconnectAttempts++
        }, delay)
      }
    })
  }
  
  return ws
}
