import { useEffect, useState } from "react";
import { getJson } from "./lib/api";

type Health = { status?: "ok" | "degraded" | "down" };

export default function App() {
  const [status, setStatus] = useState("…");
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    (async () => {
      try {
        const data = await getJson<Health>("/health");
        setStatus(data?.status ?? "unknown");
        setError(null);
      } catch (e: unknown) {
        setStatus("down");
        setError(e instanceof Error ? e.message : String(e));
      }
    })();
  }, []);

  return (
    <main style={{ padding: 20, fontFamily: "system-ui" }}>
      <h1>Go + Vite Starter</h1>
      <p>Health: {status}</p>
      {error && <p style={{ color: "crimson" }}>Error: {error}</p>}
    </main>
  );
}
