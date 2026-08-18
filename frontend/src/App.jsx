import { useState } from 'react';
import { factorizeMatrix, parseMatrixInput } from './api.js';

const DEFAULT_MATRIX = `[
  [1, 2],
  [3, 4],
  [5, 6]
]`;

function Section({ step, title, hint, children }) {
  return (
    <section className="rounded-lg border border-stone-300 bg-white p-4">
      <p className="text-xs font-medium uppercase tracking-wide text-stone-500">Paso {step}</p>
      <h2 className="mt-1 text-lg font-semibold">{title}</h2>
      {hint && <p className="mt-1 text-sm text-stone-600">{hint}</p>}
      <div className="mt-3">{children}</div>
    </section>
  );
}

function MatrixTable({ title, data }) {
  if (!data?.length) {
    return null;
  }

  return (
    <div>
      {title && <h3 className="mb-2 text-sm font-semibold text-stone-700">{title}</h3>}
      <div className="overflow-x-auto">
        <table className="text-sm">
          <tbody>
            {data.map((row, rowIndex) => (
              <tr key={rowIndex}>
                {row.map((value, columnIndex) => (
                  <td
                    key={columnIndex}
                    className="border border-stone-200 px-2 py-1 text-right font-mono"
                  >
                    {Number(value).toFixed(4)}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

export default function App() {
  const [matrixText, setMatrixText] = useState(DEFAULT_MATRIX);
  const [submittedMatrix, setSubmittedMatrix] = useState(null);
  const [result, setResult] = useState(null);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  async function handleSubmit(event) {
    event.preventDefault();
    setError('');
    setResult(null);
    setSubmittedMatrix(null);
    setLoading(true);

    try {
      const matrix = parseMatrixInput(matrixText);
      const data = await factorizeMatrix(matrix);
      setSubmittedMatrix(matrix);
      setResult(data);
    } catch (submitError) {
      setError(submitError.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="mx-auto max-w-3xl px-4 py-8">
      <header className="mb-8">
        <h1 className="text-2xl font-bold">Matrix QR</h1>
        <p className="mt-1 text-sm text-stone-600">
          Matriz → factorización QR (Q, R) → estadísticas sobre Q y R.
        </p>
      </header>

      <form onSubmit={handleSubmit} className="space-y-4 rounded-lg border border-stone-300 bg-white p-4">
        <label className="block">
          <span className="mb-1 block text-sm font-medium">Matriz de entrada (JSON)</span>
          <textarea
            value={matrixText}
            onChange={(event) => setMatrixText(event.target.value)}
            rows={8}
            className="w-full rounded border border-stone-300 px-3 py-2 font-mono text-sm"
          />
        </label>

        <button
          type="submit"
          disabled={loading}
          className="rounded bg-stone-800 px-4 py-2 text-sm text-white hover:bg-stone-700 disabled:opacity-50"
        >
          {loading ? 'Procesando…' : 'Calcular QR y estadísticas'}
        </button>
      </form>

      {error && (
        <p className="mt-4 rounded border border-red-300 bg-red-50 px-3 py-2 text-sm text-red-800">
          {error}
        </p>
      )}

      {result && submittedMatrix && (
        <div className="mt-6 space-y-4">
          <Section step={1} title="Matriz de entrada" hint="Datos enviados a go-api.">
            <MatrixTable data={submittedMatrix} />
          </Section>

          <Section
            step={2}
            title="Factorización QR"
            hint="go-api descompone A = Q × R con Gonum."
          >
            <div className="space-y-4">
              <MatrixTable title="Q" data={result.qr.Q} />
              <MatrixTable title="R" data={result.qr.R} />
            </div>
          </Section>

          <Section
            step={3}
            title="Estadísticas"
            hint="node-api calcula sobre todos los valores de Q y R."
          >
            <dl className="grid grid-cols-2 gap-3 text-sm sm:grid-cols-4">
              <div>
                <dt className="text-stone-500">Máximo</dt>
                <dd className="font-mono font-medium">{result.statistics.max}</dd>
              </div>
              <div>
                <dt className="text-stone-500">Mínimo</dt>
                <dd className="font-mono font-medium">{result.statistics.min}</dd>
              </div>
              <div>
                <dt className="text-stone-500">Promedio</dt>
                <dd className="font-mono font-medium">{result.statistics.average}</dd>
              </div>
              <div>
                <dt className="text-stone-500">Suma</dt>
                <dd className="font-mono font-medium">{result.statistics.sum}</dd>
              </div>
            </dl>
            <p className="mt-3 text-sm text-stone-600">
              ¿Matriz diagonal? {result.statistics.hasDiagonalMatrix ? 'Sí' : 'No'}
              {result.statistics.diagonalMatrices.length > 0 &&
                ` — ${result.statistics.diagonalMatrices.join(', ')}`}
            </p>
          </Section>
        </div>
      )}
    </main>
  );
}
