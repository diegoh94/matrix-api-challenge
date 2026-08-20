const API_BASE = import.meta.env.VITE_API_URL ?? '';
const API_KEY = import.meta.env.VITE_API_KEY ?? '';

const TOKEN_REFRESH_MARGIN_MS = 60_000;

let cachedToken = null;
let expiresAtMs = 0;

function isCachedTokenValid() {
  return cachedToken !== null && Date.now() < expiresAtMs - TOKEN_REFRESH_MARGIN_MS;
}

function clearTokenCache() {
  cachedToken = null;
  expiresAtMs = 0;
}

async function fetchTokenFromApi() {
  if (!API_KEY) {
    throw new Error('Frontend sin VITE_API_KEY configurada');
  }

  const response = await fetch(`${API_BASE}/auth/token`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ apiKey: API_KEY }),
  });

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.error ?? 'No se pudo obtener el token');
  }

  cachedToken = data.token;
  expiresAtMs = Date.now() + data.expiresIn * 1000;

  return cachedToken;
}

async function getAccessToken() {
  if (isCachedTokenValid()) {
    return cachedToken;
  }

  return fetchTokenFromApi();
}

async function postMatrixQr(matrix, token) {
  return fetch(`${API_BASE}/api/v1/matrix/qr`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ matrix }),
  });
}

export async function factorizeMatrix(matrix) {
  let token = await getAccessToken();
  let response = await postMatrixQr(matrix, token);

  if (response.status === 401) {
    clearTokenCache();
    token = await fetchTokenFromApi();
    response = await postMatrixQr(matrix, token);
  }

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.error ?? 'Error al factorizar la matriz');
  }

  return data;
}

export function parseMatrixInput(rawValue) {
  const parsed = JSON.parse(rawValue);

  if (!Array.isArray(parsed) || parsed.length === 0) {
    throw new Error('La matriz debe ser un array de filas');
  }

  return parsed;
}
