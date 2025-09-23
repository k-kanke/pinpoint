import axios from 'axios'

export const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || 'http://localhost:8080',
  timeout: 10000,
})

export type Pin = {
  id: string
  title: string
  body: string
  lat: number
  lon: number
  expiry_at?: string | null
  created_at: string
}

export type PinDetail = Pin & {
  photos?: { id: string; thread_id: string; url: string }[]
  comments?: { id: string; thread_id: string; body: string; created_at: string }[]
}

export async function getNearbyPins(lat: number, lon: number, radius = 1000, limit = 100) {
  const res = await api.get<Pin[]>('/api/pins', { params: { lat, lon, radius, limit } })
  return res.data
}

export async function createPin(data: { title: string; body: string; lat: number; lon: number; expiry_at?: string; photo_url?: string }) {
  const res = await api.post<{ id: string }>('/api/pins', data)
  return res.data
}

export async function getPin(id: string) {
  const res = await api.get<PinDetail>(`/api/pins/${id}`)
  return res.data
}

export async function addComment(id: string, body: string) {
  const res = await api.post<{ id: string }>(`/api/pins/${id}/comments`, { body })
  return res.data
}

export async function uploadPhoto(id: string, file: File) {
  const form = new FormData()
  form.append('file', file)
  const res = await api.post<{ url: string }>(`/api/pins/${id}/photos`, form, { headers: { 'Content-Type': 'multipart/form-data' } })
  return res.data
}
