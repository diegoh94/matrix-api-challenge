const API_BASE = import.meta.env.VITE_API_URL ?? '';
const API_KEY = import.meta.env.VITE_API_KEY ?? '';

async function requestToken() {
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

  return data.token;
}

export async function factorizeMatrix(matrix) {
  const token = await requestToken();

  const response = await fetch(`${API_BASE}/api/v1/matrix/qr`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ matrix }),
  });

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
