import React, { useState } from 'react'
import { createSession } from '../utils/api'
import './SessionSetup.css'

export default function SessionSetup({ onSessionCreated }) {
  const [userEmail, setUserEmail] = useState('')
  const [userName, setUserName] = useState('')
  const [specialty, setSpecialty] = useState('obstetrics')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      if (!userEmail || !userName) {
        throw new Error('请填写所有必填项')
      }

      const { session_id } = await createSession(userEmail, specialty)
      onSessionCreated(session_id, userEmail, specialty)
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
          <h1>信臣健康互联网医院</h1>
          <p className="subtitle">在线医生咨询服务</p>

          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label htmlFor="name">您的姓名</label>
              <input
                type="text"
                id="name"
                value={userName}
                onChange={(e) => setUserName(e.target.value)}
                placeholder="请输入您的姓名"
                disabled={loading}
              />
            </div>

            <div className="form-group">
              <label htmlFor="email">手机号/邮箱</label>
              <input
                type="email"
                id="email"
                value={userEmail}
                onChange={(e) => setUserEmail(e.target.value)}
                placeholder="请输入手机号或邮箱"
                disabled={loading}
              />
            </div>

            <div className="form-group">
              <label htmlFor="specialty">选择医生科室</label>
              <select
                id="specialty"
                value={specialty}
                onChange={(e) => setSpecialty(e.target.value)}
                disabled={loading}
                className="specialty-select"
              >
                <option value="obstetrics">👩‍⚕️ 妇产科 - 妇科、产科、产后恢复</option>
                <option value="pediatrics">👶 儿科 - 儿童疾病、生长发育</option>
                <option value="internal_medicine">🫀 内科 - 内脏器官、代谢、感染</option>
                <option value="dermatology">🩹 皮肤科 - 皮肤病、痤疮、湿疹</option>
                <option value="ent">👂 耳鼻喉科 - 鼻炎、喉咙痛、耳痛</option>
                <option value="cardiology">🫀 心脑血管科 - 胸痛、心悸、头晕、脑卒中</option>
                <option value="respiratory">💨 呼吸科 - 咳嗽、哮喘、呼吸困难</option>
              </select>
            </div>

            {error && <div className="error-message">{error}</div>}

            <button
              type="submit"
              disabled={loading}
              className="start-button"
            >
              {loading ? '正在创建会话...' : '开始咨询'}
            </button>
          </form>

          <div className="info-box">
            <h3>服务说明</h3>
            <ul>
              <li>提供妇产科、儿科、内科、皮肤科等多科室专业咨询服务</li>
              <li>各科室值班医生在线实时解答患者问题</li>
              <li>常见疾病初步诊断和个性化治疗建议</li>
              <li>严重症状立即建议到医院进一步检查和治疗</li>
              <li>所有患者对话内容严格保密，保护隐私</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  )
}
