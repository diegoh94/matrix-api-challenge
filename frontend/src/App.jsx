import { useState } from 'react';
import { factorizeMatrix, parseMatrixInput, rotateMatrix } from './api.js';

const DEFAULT_MATRIX = `[
  [1, 2],
  [3, 4],
  [5, 6]
]`;

const OPERATIONS = {
  qr: {
    label: 'Calcular QR y estadísticas',
    loadingLabel: 'Factorizando…',
    title: 'Matrix QR',
    description: 'Matriz → factorización QR (Q, R) → estadísticas sobre Q y R.',
  },
  rotate: {
    label: 'Rotar 90° y estadísticas',
    loadingLabel: 'Rotando…',
    title: 'Matrix Rotate',
    description: 'Matriz → rotación 90° horario → estadísticas sobre la matriz rotada.',
  },
};

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

function StatisticsPanel({ statistics, hint }) {
  return (
    <Section step={3} title="Estadísticas" hint={hint}>
      <dl className="grid grid-cols-2 gap-3 text-sm sm:grid-cols-4">
        <div>
          <dt className="text-stone-500">Máximo</dt>
          <dd className="font-mono font-medium">{statistics.max}</dd>
        </div>
        <div>
          <dt className="text-stone-500">Mínimo</dt>
          <dd className="font-mono font-medium">{statistics.min}</dd>
        </div>
        <div>
          <dt className="text-stone-500">Promedio</dt>
          <dd className="font-mono font-medium">{statistics.average}</dd>
        </div>
        <div>
          <dt className="text-stone-500">Suma</dt>
          <dd className="font-mono font-medium">{statistics.sum}</dd>
        </div>
      </dl>
      <p className="mt-3 text-sm text-stone-600">
        ¿Matriz diagonal? {statistics.hasDiagonalMatrix ? 'Sí' : 'No'}
        {statistics.diagonalMatrices.length > 0 &&
          ` — ${statistics.diagonalMatrices.join(', ')}`}
      </p>
    </Section>
  );
}

export default function App() {
  const [matrixText, setMatrixText] = useState(DEFAULT_MATRIX);
  const [submittedMatrix, setSubmittedMatrix] = useState(null);
  const [result, setResult] = useState(null);
  const [operation, setOperation] = useState(null);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  async function handleSubmit(event, selectedOperation) {
    event.preventDefault();
    setError('');
    setResult(null);
    setSubmittedMatrix(null);
    setOperation(selectedOperation);
    setLoading(true);

    try {
      const matrix = parseMatrixInput(matrixText);
      const data =
        selectedOperation === 'qr'
          ? await factorizeMatrix(matrix)
          : await rotateMatrix(matrix, 90);

      setSubmittedMatrix(matrix);
      setResult(data);
    } catch (submitError) {
      setError(submitError.message);
    } finally {
      setLoading(false);
    }
  }

  const activeOperation = operation ? OPERATIONS[operation] : OPERATIONS.qr;

  return (
    <main className="mx-auto max-w-3xl px-4 py-8">
      <header className="mb-8">
        <h1 className="text-2xl font-bold">Matrix API</h1>
        <p className="mt-1 text-sm text-stone-600">
          Factorización QR o rotación 90° con estadísticas desde node-api.
        </p>
      </header>

      <form className="space-y-4 rounded-lg border border-stone-300 bg-white p-4">
        <label className="block">
          <span className="mb-1 block text-sm font-medium">Matriz de entrada (JSON)</span>
          <textarea
            value={matrixText}
            onChange={(event) => setMatrixText(event.target.value)}
            rows={8}
            className="w-full rounded border border-stone-300 px-3 py-2 font-mono text-sm"
          />
        </label>

        <div className="flex flex-wrap gap-2">
          <button
            type="button"
            disabled={loading}
            onClick={(event) => handleSubmit(event, 'qr')}
            className="rounded bg-stone-800 px-4 py-2 text-sm text-white hover:bg-stone-700 disabled:opacity-50"
          >
            {loading && operation === 'qr' ? OPERATIONS.qr.loadingLabel : OPERATIONS.qr.label}
          </button>
          <button
            type="button"
            disabled={loading}
            onClick={(event) => handleSubmit(event, 'rotate')}
            className="rounded border border-stone-800 px-4 py-2 text-sm text-stone-800 hover:bg-stone-100 disabled:opacity-50"
          >
            {loading && operation === 'rotate'
              ? OPERATIONS.rotate.loadingLabel
              : OPERATIONS.rotate.label}
          </button>
        </div>
      </form>

      {error && (
        <p className="mt-4 rounded border border-red-300 bg-red-50 px-3 py-2 text-sm text-red-800">
          {error}
        </p>
      )}

      {result && submittedMatrix && operation === 'qr' && (
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

          <StatisticsPanel
            statistics={result.statistics}
            hint="node-api calcula sobre todos los valores de Q y R."
          />
        </div>
      )}

      {result && submittedMatrix && operation === 'rotate' && (
        <div className="mt-6 space-y-4">
          <Section step={1} title="Matriz de entrada" hint="Datos enviados a go-api.">
            <MatrixTable data={submittedMatrix} />
          </Section>

          <Section
            step={2}
            title={`Rotación ${result.degrees}° horario`}
            hint="Reordenamiento por índices — O(m×n), sin trigonometría."
          >
            <MatrixTable data={result.rotated.data} />
          </Section>

          <StatisticsPanel
            statistics={result.statistics}
            hint="node-api calcula sobre todos los valores de la matriz rotada."
          />
        </div>
      )}

      {!result && (
        <p className="mt-4 text-sm text-stone-500">
          Operación activa: {activeOperation.description}
        </p>
      )}
    </main>
  );
}
