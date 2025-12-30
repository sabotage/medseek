import axios from 'axios'

const API_BASE_URL = '/api'

export const createSession = async (userId) => {
  try {
    const response = await axios.post(`${API_BASE_URL}/session/create`, {
      user_id: userId,
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
  return new WebSocket(wsUrl)
}
