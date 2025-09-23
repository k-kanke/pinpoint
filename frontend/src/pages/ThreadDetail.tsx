import React, { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { addComment, getPin, PinDetail, uploadPhoto } from '../api/client'

export default function ThreadDetail() {
  const { id } = useParams()
  const [pin, setPin] = useState<PinDetail | null>(null)
  const [comment, setComment] = useState('')
  const [file, setFile] = useState<File | null>(null)

  useEffect(() => {
    if (!id) return
    getPin(id).then(setPin).catch(console.error)
  }, [id])

  const onPost = async () => {
    if (!id || !comment) return
    await addComment(id, comment)
    setComment('')
    const p = await getPin(id)
    setPin(p)
  }

  const onUpload = async () => {
    if (!id || !file) return
    await uploadPhoto(id, file)
    setFile(null)
    const p = await getPin(id)
    setPin(p)
  }

  if (!pin) return <div style={{ padding: 12 }}>Loading...</div>

  return (
    <div style={{ padding: 12 }}>
      <Link to="/">← Back</Link>
      <h2>{pin.title}</h2>
      <p>{pin.body}</p>
      <div style={{ color: '#666', fontSize: 12 }}>{new Date(pin.created_at).toLocaleString()}</div>

      {pin.photos && pin.photos.length > 0 && (
        <div style={{ marginTop: 16 }}>
          <h3>Photos</h3>
          <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap' }}>
            {pin.photos.map(ph => (
              <img key={ph.id} src={ph.url} alt="" style={{ maxWidth: 200, maxHeight: 200, objectFit: 'cover' }} />
            ))}
          </div>
        </div>
      )}

      <div style={{ marginTop: 16 }}>
        <h3>Upload Photo</h3>
        <input type="file" accept="image/*" onChange={e => setFile(e.target.files?.[0] || null)} />
        <button onClick={onUpload} disabled={!file} style={{ marginLeft: 8 }}>Upload</button>
      </div>

      <div style={{ marginTop: 16 }}>
        <h3>Comments</h3>
        <ul>
          {(pin.comments || []).map(c => (
            <li key={c.id}>
              {c.body}
              <div style={{ color: '#666', fontSize: 12 }}>{new Date(c.created_at).toLocaleString()}</div>
            </li>
          ))}
        </ul>
        <textarea rows={3} value={comment} onChange={e => setComment(e.target.value)} style={{ width: 400, display: 'block', marginTop: 8 }} />
        <button onClick={onPost} disabled={!comment}>Post Comment</button>
      </div>
    </div>
  )
}

