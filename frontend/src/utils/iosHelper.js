/**
 * iOS Safari 兼容性辅助函数
 * 处理iOS特定的浏览器问题
 */

/**
 * 初始化iOS Safari修复
 * - 防止缩放问题
 * - 处理键盘弹出
 * - 处理视口高度
 */
export function initializeIOSFixes() {
  // 检测是否为iOS
  const isIOS = /iPad|iPhone|iPod/.test(navigator.userAgent)
  
  if (!isIOS) return

  // 修复1: 处理键盘弹出时的视口变化
  let lastWindowHeight = window.innerHeight
  
  window.addEventListener('resize', () => {
    const currentHeight = window.innerHeight
    
    // 检测键盘是否弹出（高度变化超过50px）
    if (Math.abs(currentHeight - lastWindowHeight) > 50) {
      // 键盘弹出或关闭时，延迟滚动到底部
      setTimeout(() => {
        const messagesContainer = document.querySelector('.messages-container')
        if (messagesContainer) {
          messagesContainer.scrollTop = messagesContainer.scrollHeight
        }
      }, 100)
    }
    
    lastWindowHeight = currentHeight
  })

  // 修复2: 处理输入框焦点时的问题
  const inputs = document.querySelectorAll('input, textarea')
  inputs.forEach(input => {
    input.addEventListener('focus', () => {
      // 延迟滚动到输入框
      setTimeout(() => {
        input.scrollIntoView({ behavior: 'smooth', block: 'center' })
      }, 300)
    })
  })

  // 修复3: 防止默认的iOS缩放行为
  document.addEventListener('touchmove', (e) => {
    // 仅在非input/textarea元素上防止默认行为
    if (!isInputElement(e.target)) {
      // 允许滚动，但防止缩放
    }
  }, false)

  // 修复4: 处理iOS Safari的100vh问题
  const setViewportHeight = () => {
    const vh = window.innerHeight * 0.01
    document.documentElement.style.setProperty('--vh', `${vh}px`)
  }
  
  setViewportHeight()
  window.addEventListener('resize', setViewportHeight)
  window.addEventListener('orientationchange', setViewportHeight)

  // 修复5: 防止iOS Safari的表单自动缩放
  document.addEventListener('touchstart', function(e) {
    if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA') {
      e.target.style.fontSize = '16px'
    }
  }, false)
}

/**
 * 检查元素是否为输入元素
 */
function isInputElement(element) {
  const tagName = element.tagName.toLowerCase()
  return tagName === 'input' || tagName === 'textarea' || tagName === 'select'
}

/**
 * 强制滚动到底部（用于消息容器）
 */
export function scrollToBottom(element, smooth = true) {
  if (!element) return
  
  setTimeout(() => {
    element.scrollTop = element.scrollHeight
    
    // iOS Safari的额外处理
    const isIOS = /iPad|iPhone|iPod/.test(navigator.userAgent)
    if (isIOS && smooth) {
      // iOS在某些情况下需要额外的时间来完成滚动
      setTimeout(() => {
        element.scrollTop = element.scrollHeight
      }, 100)
    }
  }, 0)
}

/**
 * 处理iOS键盘的显示/隐藏
 */
export function onIOSKeyboardToggle(callback) {
  const isIOS = /iPad|iPhone|iPod/.test(navigator.userAgent)
  if (!isIOS) return
  
  let lastHeight = window.innerHeight
  
  const handleResize = () => {
    const currentHeight = window.innerHeight
    const keyboardHeight = lastHeight - currentHeight
    
    if (keyboardHeight > 50) {
      callback('show', keyboardHeight)
    } else if (keyboardHeight < -50) {
      callback('hide', 0)
    }
    
    lastHeight = currentHeight
  }
  
  window.addEventListener('resize', handleResize)
  
  return () => {
    window.removeEventListener('resize', handleResize)
  }
}

/**
 * 获取安全区域内距 (用于有刘海的设备)
 */
export function getSafeAreaInsets() {
  const style = getComputedStyle(document.documentElement)
  
  return {
    top: parseFloat(style.getPropertyValue('env(safe-area-inset-top)')) || 0,
    right: parseFloat(style.getPropertyValue('env(safe-area-inset-right)')) || 0,
    bottom: parseFloat(style.getPropertyValue('env(safe-area-inset-bottom)')) || 0,
    left: parseFloat(style.getPropertyValue('env(safe-area-inset-left)')) || 0,
  }
}

/**
 * 检查是否为iOS设备
 */
export function isIOSDevice() {
  return /iPad|iPhone|iPod/.test(navigator.userAgent)
}

/**
 * 检查是否为iOS Safari浏览器
 */
export function isIOSSafari() {
  const userAgent = navigator.userAgent
  const isIOS = /iPad|iPhone|iPod/.test(userAgent)
  const isNotChrome = !/CriOS/.test(userAgent)
  const isNotFirefox = !/FxiOS/.test(userAgent)
  
  return isIOS && isNotChrome && isNotFirefox
}
