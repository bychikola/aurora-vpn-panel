import { Cpu, HardDrive, MemoryStick, Wifi, WifiOff, AlertTriangle } from 'lucide-react';
import { useNodes } from '../api/hooks';
import { formatBytes } from '../lib/utils';
import type { Node } from '../types';

function NodeCard({ node }: { node: Node }) {
  const statusIcon = {
    online: <Wifi className="w-4 h-4" style={{ color: 'var(--color-success)' }} />,
    offline: <WifiOff className="w-4 h-4" style={{ color: 'var(--color-danger)' }} />,
    degraded: <AlertTriangle className="w-4 h-4" style={{ color: 'var(--color-warning)' }} />,
  }[node.status];

  const statusBadge = {
    online: 'badge-success',
    offline: 'badge-danger',
    degraded: 'badge-warning',
  }[node.status];

  return (
    <div
      className="card-aurora rounded-xl p-5"
      style={{ background: 'var(--color-polar-850)' }}
    >
      {/* Header */}
      <div className="flex items-center justify-between mb-4">
        <div>
          <div className="flex items-center gap-2">
            {statusIcon}
            <span className="font-semibold text-sm" style={{ color: 'var(--color-starlight)' }}>
              {node.name}
            </span>
          </div>
          <div className="text-xs mt-1 font-mono-num" style={{ color: 'var(--color-ice)' }}>
            {node.host}:{node.port} &middot; API :{node.apiPort}
          </div>
        </div>
        <span className={`badge ${statusBadge}`}>{node.status}</span>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-3 gap-3 mb-4">
        <div>
          <div className="flex items-center gap-1.5 mb-1">
            <Cpu className="w-3 h-3" style={{ color: 'var(--color-ice)' }} />
            <span className="text-[10px] uppercase tracking-wider" style={{ color: 'var(--color-ice)' }}>CPU</span>
          </div>
          <div className="flex items-center gap-2">
            <div className="flex-1 h-1.5 rounded-full" style={{ background: 'var(--color-polar-700)' }}>
              <div
                className="h-full rounded-full"
                style={{
                  width: `${node.cpuPercent}%`,
                  background: node.cpuPercent > 80 ? 'var(--color-danger)' : node.cpuPercent > 60 ? 'var(--color-warning)' : 'var(--color-aurora-green)',
                }}
              />
            </div>
            <span className="text-xs font-mono-num" style={{ color: 'var(--color-frost)' }}>{node.cpuPercent}%</span>
          </div>
        </div>
        <div>
          <div className="flex items-center gap-1.5 mb-1">
            <MemoryStick className="w-3 h-3" style={{ color: 'var(--color-ice)' }} />
            <span className="text-[10px] uppercase tracking-wider" style={{ color: 'var(--color-ice)' }}>RAM</span>
          </div>
          <div className="flex items-center gap-2">
            <div className="flex-1 h-1.5 rounded-full" style={{ background: 'var(--color-polar-700)' }}>
              <div
                className="h-full rounded-full"
                style={{
                  width: `${node.memoryPercent}%`,
                  background: node.memoryPercent > 80 ? 'var(--color-danger)' : node.memoryPercent > 60 ? 'var(--color-warning)' : 'var(--color-aurora-cyan)',
                }}
              />
            </div>
            <span className="text-xs font-mono-num" style={{ color: 'var(--color-frost)' }}>{node.memoryPercent}%</span>
          </div>
        </div>
        <div>
          <div className="flex items-center gap-1.5 mb-1">
            <HardDrive className="w-3 h-3" style={{ color: 'var(--color-ice)' }} />
            <span className="text-[10px] uppercase tracking-wider" style={{ color: 'var(--color-ice)' }}>Disk</span>
          </div>
          <div className="flex items-center gap-2">
            <div className="flex-1 h-1.5 rounded-full" style={{ background: 'var(--color-polar-700)' }}>
              <div
                className="h-full rounded-full"
                style={{
                  width: `${node.diskPercent}%`,
                  background: 'var(--color-aurora-violet)',
                }}
              />
            </div>
            <span className="text-xs font-mono-num" style={{ color: 'var(--color-frost)' }}>{node.diskPercent}%</span>
          </div>
        </div>
      </div>

      {/* Traffic */}
      <div className="grid grid-cols-2 gap-3 pt-3 border-t" style={{ borderColor: 'var(--color-polar-700)' }}>
        <div>
          <span className="text-[10px] uppercase tracking-wider" style={{ color: 'var(--color-ice)' }}>↓ Downlink</span>
          <div className="text-xs font-mono-num mt-0.5" style={{ color: 'var(--color-aurora-cyan)' }}>
            {formatBytes(node.downlinkSpeed)}/s
          </div>
          <div className="text-[10px] mt-0.5" style={{ color: 'var(--color-ice)' }}>
            {formatBytes(node.downlinkTotal)} total
          </div>
        </div>
        <div>
          <span className="text-[10px] uppercase tracking-wider" style={{ color: 'var(--color-ice)' }}>↑ Uplink</span>
          <div className="text-xs font-mono-num mt-0.5" style={{ color: 'var(--color-aurora-violet)' }}>
            {formatBytes(node.uplinkSpeed)}/s
          </div>
          <div className="text-[10px] mt-0.5" style={{ color: 'var(--color-ice)' }}>
            {formatBytes(node.uplinkTotal)} total
          </div>
        </div>
      </div>

      {/* Meta */}
      <div className="flex items-center justify-between mt-3 pt-3 border-t text-[10px]" style={{ borderColor: 'var(--color-polar-700)', color: 'var(--color-ice)' }}>
        <span>{node.version}</span>
        <span>{node.userCount} users &middot; {node.inboundCount} inbounds</span>
        <span>{node.location}</span>
      </div>
    </div>
  );
}

export default function Nodes() {
  const { data: nodes, isLoading } = useNodes();

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-lg font-semibold m-0" style={{ color: 'var(--color-starlight)' }}>Nodes</h2>
          <p className="text-sm mt-1" style={{ color: 'var(--color-frost)' }}>
            {nodes ? `${nodes.length} servers` : 'Loading...'}
          </p>
        </div>
        <div className="flex items-center gap-3 text-xs" style={{ color: 'var(--color-frost)' }}>
          <span className="flex items-center gap-1"><Wifi className="w-3 h-3" style={{ color: 'var(--color-success)' }} /> Online</span>
          <span className="flex items-center gap-1"><AlertTriangle className="w-3 h-3" style={{ color: 'var(--color-warning)' }} /> Degraded</span>
          <span className="flex items-center gap-1"><WifiOff className="w-3 h-3" style={{ color: 'var(--color-danger)' }} /> Offline</span>
        </div>
      </div>

      {isLoading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {Array.from({ length: 4 }).map((_, i) => (
            <div key={i} className="skeleton h-60 rounded-xl" />
          ))}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {nodes?.map((node) => (
            <NodeCard key={node.id} node={node} />
          ))}
        </div>
      )}
    </div>
  );
}
