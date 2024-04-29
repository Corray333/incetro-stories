import axios from 'axios'

function getCookie(name) {
    const value = `; ${document.cookie}`
    const parts = value.split(`; ${name}=`)
    if (parts.length === 2) return parts.pop().split(';').shift()
}

const refreshTokens = async () => {
    try {
        console.log(localStorage.getItem('Refresh'))
        let { data } = await axios.get('http://localhost:3001/api/users/refresh', {
            headers: {
                'Refresh': localStorage.getItem('Refresh'),
            }
        })

        localStorage.setItem('Authorization', data.authorization)
        localStorage.setItem('Refresh', data.refresh)
    } catch (error) {
        alert('Error refreshing tokens')
    }
}

export { getCookie, refreshTokens };