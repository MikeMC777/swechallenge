CREATE DATABASE IF NOT EXISTS stocks;
SET DATABASE = stocks;

-- 1) Páginas crudas para auditoría
CREATE TABLE IF NOT EXISTS stocks_raw (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  source TEXT NOT NULL,
  fetched_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  page_key TEXT,
  payload JSONB NOT NULL
);

-- 2) Eventos de analistas normalizados (uno por ítem del API)
--    { ticker, target_from, target_to, company, action, brokerage, rating_from, rating_to, time }
CREATE TABLE IF NOT EXISTS analyst_events (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  ticker TEXT NOT NULL,
  company TEXT,
  action TEXT,                 -- "target raised by" | "reiterated by" | "target lowered by" | "initiated by" ...
  brokerage TEXT,
  rating_from TEXT,
  rating_to TEXT,
  target_from NUMERIC,         -- en USD, parseado de "${n}"
  target_to NUMERIC,           -- en USD, parseado de "${n}"
  target_delta NUMERIC,        -- target_to - target_from (NULL si falta from)
  event_time TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_analyst_events_ticker_time ON analyst_events (ticker, event_time DESC);

-- 3) Resumen por ticker (para la lista de la UI)
CREATE TABLE IF NOT EXISTS tickers_summary (
  ticker TEXT PRIMARY KEY,
  company TEXT,
  last_action TEXT,
  last_brokerage TEXT,
  last_rating TEXT,            -- rating_to del último evento
  last_target NUMERIC,         -- target_to del último evento
  last_time TIMESTAMPTZ,
  raised_30d INT NOT NULL DEFAULT 0,
  lowered_30d INT NOT NULL DEFAULT 0,
  reiterated_30d INT NOT NULL DEFAULT 0,
  initiated_30d INT NOT NULL DEFAULT 0,
  score NUMERIC NOT NULL DEFAULT 0, -- puntaje recomendado calculado
  rationale TEXT NOT NULL DEFAULT '',
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
