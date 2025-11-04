const base = import.meta.env.VITE_API_URL ?? "/api";
const url = (path: string) => `${base}${path}`;

export async function getJson<T>(path: string, init?: RequestInit) {
  const res = await fetch(url(path), {
    ...init,
    headers: { Accept: "application/json", ...(init?.headers || {}) },
  });
  const ct = res.headers.get("content-type") || "";
  if (!res.ok) throw new Error(`HTTP ${res.status}`);
  if (!ct.includes("application/json")) {
    const txt = await res.text().catch(() => "");
    throw new Error(
      `Expected JSON, got ${ct || "unknown"} — ${txt.slice(0, 80)}`
    );
  }
  return res.json() as Promise<T>;
}
