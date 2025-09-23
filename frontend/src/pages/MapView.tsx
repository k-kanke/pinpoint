import React, { useEffect, useMemo, useRef, useState } from 'react'
import maplibregl, { Map, Marker } from 'maplibre-gl'
import { createPin, getNearbyPins, Pin } from '../api/client'
import { useNavigate } from 'react-router-dom'

export default function MapView() {
  const mapRef = useRef<HTMLDivElement | null>(null)
  const mapObj = useRef<Map | null>(null)
  const [center, setCenter] = useState<[number, number]>([139.767, 35.681]) // Tokyo default
  const [pins, setPins] = useState<Pin[]>([])
  const markers = useRef<Marker[]>([])
  const navigate = useNavigate()

  // init map
  useEffect(() => {
    if (!mapRef.current || mapObj.current) return
    const map = new maplibregl.Map({
      container: mapRef.current,
      style: 'https://demotiles.maplibre.org/style.json',
      center,
      zoom: 13,
    })
    mapObj.current = map

    map.on('moveend', () => {
      const c = map.getCenter()
      loadPins(c.lat, c.lng)
    })

    // get geolocation once
    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition((pos) => {
        const lng = pos.coords.longitude
        const lat = pos.coords.latitude
        setCenter([lng, lat])
        map.setCenter([lng, lat])
        loadPins(lat, lng)
      }, () => loadPins(35.681, 139.767))
    } else {
      loadPins(35.681, 139.767)
    }

    return () => map.remove()
  }, [])

  // render markers
  useEffect(() => {
    const map = mapObj.current
    if (!map) return
    // clear old markers
    markers.current.forEach(m => m.remove())
    markers.current = []
    pins.forEach((p) => {
      const el = document.createElement('div')
      el.style.width = '12px'
      el.style.height = '12px'
      el.style.background = '#e11'
      el.style.borderRadius = '50%'
      el.style.border = '2px solid #fff'
      el.style.boxShadow = '0 0 2px rgba(0,0,0,0.4)'
      el.style.cursor = 'pointer'
      el.title = p.title
      el.addEventListener('click', () => navigate(`/pins/${p.id}`))
      const mk = new Marker({ element: el }).setLngLat([p.lon, p.lat]).addTo(map)
      markers.current.push(mk)
    })
  }, [pins])

  async function loadPins(lat: number, lon: number) {
    try {
      const list = await getNearbyPins(lat, lon, 1000, 100)
      setPins(list)
    } catch (e) {
      console.error(e)
    }
  }

  // simple create form
  const [title, setTitle] = useState('')
  const [body, setBody] = useState('')
  const canSubmit = useMemo(() => title.length > 0 && body.length > 0, [title, body])

  const onCreate = async () => {
    const map = mapObj.current
    if (!map) return
    const c = map.getCenter()
    try {
      const res = await createPin({ title, body, lat: c.lat, lon: c.lng })
      setTitle(''); setBody('')
      // reload pins
      loadPins(c.lat, c.lng)
      navigate(`/pins/${res.id}`)
    } catch (e) {
      console.error(e)
      alert('Failed to create pin')
    }
  }

  return (
    <div style={{ display: 'flex', height: '100%' }}>
      <div style={{ flex: 1, position: 'relative' }}>
        <div ref={mapRef} style={{ position: 'absolute', inset: 0 }} />
      </div>
      <aside style={{ width: 360, borderLeft: '1px solid #eee', padding: 12, overflow: 'auto' }}>
        <h3>Create Pin</h3>
        <input value={title} onChange={e => setTitle(e.target.value)} placeholder="Title" style={{ width: '100%', marginBottom: 8 }} />
        <textarea value={body} onChange={e => setBody(e.target.value)} placeholder="Body" rows={4} style={{ width: '100%', marginBottom: 8 }} />
        <button onClick={onCreate} disabled={!canSubmit}>Post</button>

        <h3 style={{ marginTop: 24 }}>Nearby Pins</h3>
        <ul style={{ listStyle: 'none', padding: 0 }}>
          {pins.map(p => (
            <li key={p.id} style={{ padding: '8px 0', borderBottom: '1px solid #f0f0f0' }}>
              <a href={`#/pins/${p.id}`} onClick={(e) => { e.preventDefault(); navigate(`/pins/${p.id}`) }}>{p.title}</a>
              <div style={{ color: '#666', fontSize: 12 }}>{new Date(p.created_at).toLocaleString()}</div>
            </li>
          ))}
        </ul>
      </aside>
    </div>
  )
}
