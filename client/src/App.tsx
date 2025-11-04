import { useEffect, useState } from 'react'

type ApiResponse = { message: string }

export default function App() {
  const [msg, setMsg] = useState<string>('loading...')

  useEffect(() => {
    fetch('/api/message')
      .then(r => r.json())
      .then((data: ApiResponse) => setMsg(data.message))
      .catch(() => setMsg('failed to load'))
  }, [])

  return (
    <main>
      <h1>Go + Gin + Vite</h1>
      <p>API says: <strong>{msg}</strong></p>
    </main>
  )
}
