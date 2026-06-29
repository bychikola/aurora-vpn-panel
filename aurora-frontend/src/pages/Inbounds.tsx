import { Shield, Plug, ArrowUpDown } from 'lucide-react';
import { useInbounds } from '../api/hooks';
import { formatBytes } from '../lib/utils';
import type { Inbound } from '../types';

const PROTOCOL_COLORS: Record<string, string> = {
  vless: 'var(--color-aurora-green)',
  vmess: 'var(--color-aurora-cyan)',
  trojan: 'var(--color-aurora-violet)',
  shadowsocks: 'var(--color-aurora-blue)',
  'shadowsocks-2022': 'var(--color-aurora-pink)',
  hysteria2: 'var(--color-warning)',
  'tuic-v5': 'var(--color-frost)',
};

function InboundRow({ inbound }: { inbound: Inbound }) {
  return (
    <div
      className="card-aurora rounded-xl p-5 flex flex-col gap-3"
      style={{ background: 'var(--color-polar-850)' }}
    >
      {/* Top row */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div
            className="w-3 h-3 rounded-full"
            style={{ background: inbound.enable ? PROTOCOL_COLORS[inbound.protocol] ?? 'var(--color-ice)' : 'var(--color-polar-500)' }}
          />
          <div>
            <div className="flex items-center gap-2">
              <span className="font-semibold text-sm" style={{ color: 'var(--color-starlight)' }}>
                {inbound.tag}
              </span>
              <span className="badge text-[10px] px-2 py-0.5" style={{
                background: `${PROTOCOL_COLORS[inbound.protocol]}20`,
                color: PROTOCOL_COLORS[inbound.protocol],
              }}>
                {inbound.protocol.toUpperCase()}
              </span>
            </div>
            <div className="text-xs mt-0.5 font-mono-num" style={{ color: 'var(--color-ice)' }}>
              {inbound.listen}:{inbound.port} &middot; {inbound.nodeName}
            </div>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <span
            className="inline-flex items-center gap-1 px-2 py-0.5 rounded text-[10px] font-medium"
            style={{
              background: inbound.enable ? 'rgba(74,222,128,0.12)' : 'rgba(148,163,184,0.08)',
              color: inbound.enable ? 'var(--color-success)' : 'var(--color-ice)',
            }}
          >
            {inbound.enable ? 'Active' : 'Disabled'}
          </span>
        </div>
      </div>

      {/* Detail chips */}
      <div className="flex flex-wrap gap-2">
        <span
          className="inline-flex items-center gap-1 px-2 py-1 rounded-md text-[10px] font-medium"
          style={{ background: 'var(--color-polar-750)', color: 'var(--color-frost)' }}
        >
          <Shield className="w-3 h-3" />
          {inbound.security === 'none' ? 'No Security' : inbound.security.toUpperCase()}
        </span>
        <span
          className="inline-flex items-center gap-1 px-2 py-1 rounded-md text-[10px] font-medium"
          style={{ background: 'var(--color-polar-750)', color: 'var(--color-frost)' }}
        >
          <Plug className="w-3 h-3" />
          {inbound.transport.toUpperCase()}
        </span>
      </div>

      {/* Traffic bar */}
      <div className="grid grid-cols-2 gap-4 pt-2 border-t" style={{ borderColor: 'var(--color-polar-700)' }}>
        <div>
          <div className="flex items-center gap-1 mb-1">
            <ArrowUpDown className="w-3 h-3" style={{ color: 'var(--color-ice)' }} />
            <span className="text-[10px] uppercase tracking-wider" style={{ color: 'var(--color-ice)' }}>
              Traffic
            </span>
          </div>
          <div className="text-xs font-mono-num" style={{ color: 'var(--color-frost)' }}>
            ↓ {formatBytes(inbound.download)} &middot; ↑ {formatBytes(inbound.upload)}
          </div>
        </div>
        <div className="flex items-end justify-end">
          <span className="text-xs font-mono-num" style={{ color: 'var(--color-ice)' }}>
            {inbound.userCount} users
          </span>
        </div>
      </div>
    </div>
  );
}

export default function Inbounds() {
  const { data: inbounds, isLoading } = useInbounds();

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-lg font-semibold m-0" style={{ color: 'var(--color-starlight)' }}>Inbounds</h2>
          <p className="text-sm mt-1" style={{ color: 'var(--color-frost)' }}>
            {inbounds ? `${inbounds.length} configured` : 'Loading...'}
          </p>
        </div>
      </div>

      {isLoading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {Array.from({ length: 6 }).map((_, i) => (
            <div key={i} className="skeleton h-36 rounded-xl" />
          ))}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {inbounds?.map((inbound) => (
            <InboundRow key={inbound.id} inbound={inbound} />
          ))}
        </div>
      )}
    </div>
  );
}
