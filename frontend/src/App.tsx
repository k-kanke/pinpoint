import React from 'react'
import { Outlet, Link } from 'react-router-dom'

export default function App() {
  return (
    <div style={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      <header style={{ padding: '8px 12px', borderBottom: '1px solid #eee' }}>
        <Link to="/" style={{ textDecoration: 'none', fontWeight: 700 }}>Pinpoint</Link>
      </header>
      <main style={{ flex: 1 }}>
        <Outlet />
      </main>
    </div>
  )
}

