export function authHeader() {
  const user = JSON.parse(localStorage.getItem('user'))
  if (user && user.token) {
    return { Authorization: 'Bearer ' + user.token }
  }
  return {}
}

export function isAdmin() {
  const user = JSON.parse(localStorage.getItem('user'))
  return user && user.role >= 10
}

export function isRoot() {
  const user = JSON.parse(localStorage.getItem('user'))
  return user && user.role >= 100
}

export function getUser() {
  return JSON.parse(localStorage.getItem('user'))
}
