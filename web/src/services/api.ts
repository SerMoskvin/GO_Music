const API_BASE = 'http://localhost:8080'

interface ApiResponse {
  status: string
  data?: any
  error?: string
}

export const api = {
  async request(endpoint: string, options: RequestInit = {}): Promise<any> {
    const token = localStorage.getItem('token')
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    }
    
    // Добавляем Authorization если есть токен
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }

    // Мерджим заголовки
    const finalHeaders = {
      ...headers,
      ...options.headers as Record<string, string>,
    }

    const response = await fetch(`${API_BASE}${endpoint}`, {
      ...options,
      headers: finalHeaders,
    })

    const result: ApiResponse = await response.json()
    
    if (result.status === 'success') {
      return result.data
    } else {
      throw new Error(result.error || 'Request failed')
    }
  },

  get(endpoint: string): Promise<any> {
    return this.request(endpoint)
  },

  post(endpoint: string, data: any): Promise<any> {
    return this.request(endpoint, {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },

  put(endpoint: string, data: any): Promise<any> {
    return this.request(endpoint, {
      method: 'PUT', 
      body: JSON.stringify(data),
    })
  },

  patch(endpoint: string, data: any): Promise<any> {
    return this.request(endpoint, {
      method: 'PATCH',
      body: JSON.stringify(data),
    })
  },

  delete(endpoint: string): Promise<any> {
    return this.request(endpoint, {
      method: 'DELETE',
    })
  }
}