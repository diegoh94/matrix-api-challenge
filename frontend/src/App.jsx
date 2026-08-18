import { useState } from 'react';
import { factorizeMatrix, parseMatrixInput } from './api.js';

const DEFAULT_MATRIX = `[
  [1, 2],
  [3, 4],
  [5, 6]
]`;

function MatrixTable({ title, data }) {
  if (!data?.length) {
    return null;
  }

  return (
    <div className="rounded border border-stone-300 bg-white p-3">
      <h3 className="mb-2 text-sm font-semibold text-stone-700">{title}</h3>
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
  const [result, setResult] = useState(null);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  async function handleSubmit(event) {
    event.preventDefault();
    setError('');
    setResult(null);
    setLoading(true);

    try {
      const matrix = parseMatrixInput(matrixText);
      const data = await factorizeMatrix(matrix);
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
          Envía una matriz a go-api y muestra Q, R y estadísticas.
        </p>
      </header>

      <form onSubmit={handleSubmit} className="space-y-4 rounded-lg border border-stone-300 bg-white p-4">
        <label className="block">
          <span className="mb-1 block text-sm font-medium">Matriz (JSON)</span>
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
          {loading ? 'Calculando…' : 'Factorizar QR'}
        </button>
      </form>

      {error && (
        <p className="mt-4 rounded border border-red-300 bg-red-50 px-3 py-2 text-sm text-red-800">
          {error}
        </p>
      )}

      {result && (
        <section className="mt-6 space-y-4">
          <div className="rounded border border-stone-300 bg-white p-3 text-sm">
            <p>
              Matriz entrada: {result.input.rows}×{result.input.cols}
            </p>
            <p className="mt-2 grid grid-cols-2 gap-2 sm:grid-cols-4">
              <span>max: {result.statistics.max}</span>
              <span>min: {result.statistics.min}</span>
              <span>avg: {result.statistics.average}</span>
              <span>sum: {result.statistics.sum}</span>
            </p>
            <p className="mt-2 text-stone-600">
              Diagonal: {result.statistics.hasDiagonalMatrix ? 'sí' : 'no'}
              {result.statistics.diagonalMatrices.length > 0 &&
                ` (${result.statistics.diagonalMatrices.join(', ')})`}
            </p>
          </div>

          <MatrixTable title="Q" data={result.qr.Q} />
          <MatrixTable title="R" data={result.qr.R} />
        </section>
      )}
    </main>
  );
}
