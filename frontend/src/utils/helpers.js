import axios from 'axios'

function getCookie(name) {
    const value = `; ${document.cookie}`
    const parts = value.split(`; ${name}=`)
    if (parts.length === 2) return parts.pop().split(';').shift()
}

const refreshTokens = async () => {
    try {
        let { data } = await axios.get('http://localhost:3001/api/users/refresh', {
            headers: {
                'Refresh': getCookie('Refresh'),
            }
        })

        document.cookie = `Authorization=${data.authorization};`
        document.cookie = `Refresh=${data.refresh};`
    } catch (error) {
        alert('Error refreshing tokens')
    }
}

export { getCookie, refreshTokens };