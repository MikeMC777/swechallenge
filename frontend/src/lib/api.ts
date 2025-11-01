const API = import.meta.env.VITE_API_URL || 'http://localhost:8081'

export type TickerRow = {
  ticker: string
  company?: string
  last_action?: string
  last_brokerage?: string
  last_rating?: string
  last_target?: number
  last_time: string
  raised_30d: number
  lowered_30d: number
  reiterated_30d: number
  initiated_30d: number
  score: number
  rationale: string
}

export async function listTickers(
  q = '',
  sort = 'score',
  dir: 'ASC' | 'DESC' = 'DESC',
  offset = 0,
  limit = 20
) {
  const u = new URL(API + '/api/stocks')
  if (q) u.searchParams.set('q', q)
  u.searchParams.set('sort', sort)
  u.searchParams.set('dir', dir)
  u.searchParams.set('offset', String(offset))
  u.searchParams.set('limit', String(limit))
  const res = await fetch(u)
  if (!res.ok) throw new Error('failed: listTickers')
  return res.json() as Promise<{ items: TickerRow[]; limit: number; offset: number }>
}

export async function getTicker(symbol: string) {
  const res = await fetch(`${API}/api/stocks/${symbol}`)
  if (!res.ok) throw new Error('not found')
  return res.json() as Promise<{
    ticker: string
    events: Array<{
      company?: string
      action?: string
      brokerage?: string
      rating_from?: string
      rating_to?: string
      target_from?: number
      target_to?: number
      target_delta?: number
      time: string
    }>
  }>
}

export async function recommend() {
  const res = await fetch(API + '/api/recommendations')
  if (!res.ok) throw new Error('failed: recommend')
  return res.json() as Promise<{
    items: { symbol: string; name?: string; score: number; rationale: string }[]
  }>
}
