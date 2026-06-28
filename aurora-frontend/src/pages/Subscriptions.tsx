import { Copy, Key, Link2, QrCode } from 'lucide-react';
import { useSubscriptions } from '../api/hooks';
import { formatDateTime } from '../lib/utils';

export default function Subscriptions() {
  const { data: subscriptions, isLoading } = useSubscriptions();

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-lg font-semibold m-0" style={{ color: 'var(--color-starlight)' }}>Subscriptions</h2>
          <p className="text-sm mt-1" style={{ color: 'var(--color-frost)' }}>
            {subscriptions ? `${subscriptions.length} active tokens` : 'Loading...'}
          </p>
        </div>
      </div>

      {isLoading ? (
        <div className="space-y-3">
          {Array.from({ length: 5 }).map((_, i) => (
            <div key={i} className="skeleton h-16 rounded-xl" />
          ))}
        </div>
      ) : (
        <div className="space-y-3">
          {subscriptions?.map((sub) => (
            <div
              key={sub.id}
              className="card-aurora rounded-xl p-4 flex items-center justify-between"
              style={{ background: 'var(--color-polar-850)' }}
            >
              <div className="flex items-center gap-4">
                <Key className="w-5 h-5" style={{ color: sub.enabled ? 'var(--color-aurora-green)' : 'var(--color-ice)' }} />
                <div>
                  <div className="flex items-center gap-2">
                    <span className="font-medium text-sm" style={{ color: 'var(--color-starlight)' }}>
                      {sub.username}
                    </span>
                    <span className={`badge ${sub.enabled ? 'badge-success' : 'badge-warning'}`}>
                      {sub.enabled ? 'Active' : 'Disabled'}
                    </span>
                  </div>
                  <div className="text-xs font-mono-num mt-0.5" style={{ color: 'var(--color-ice)' }}>
                    {sub.url}
                  </div>
                </div>
              </div>

              <div className="flex items-center gap-3">
                <span
                  className="badge text-[10px] px-2 py-0.5"
                  style={{ background: 'var(--color-polar-750)', color: 'var(--color-aurora-cyan)' }}
                >
                  {sub.format}
                </span>
                {sub.lastRequestAt && (
                  <span className="text-[10px]" style={{ color: 'var(--color-ice)' }}>
                    Last: {formatDateTime(sub.lastRequestAt)}
                  </span>
                )}
                <div className="flex items-center gap-1">
                  <button
                    className="p-1.5 rounded-md hover:bg-white/5 transition-colors"
                    onClick={() => navigator.clipboard.writeText(sub.url)}
                    title="Copy link"
                  >
                    <Copy className="w-3.5 h-3.5" style={{ color: 'var(--color-frost)' }} />
                  </button>
                  <button
                    className="p-1.5 rounded-md hover:bg-white/5 transition-colors"
                    title="QR Code"
                    onClick={() => window.alert(`QR Code for ${sub.username}`)}
                  >
                    <QrCode className="w-3.5 h-3.5" style={{ color: 'var(--color-frost)' }} />
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
