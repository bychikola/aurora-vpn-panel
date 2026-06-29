import {
  AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer,
  PieChart, Pie, Cell, BarChart, Bar, Legend,
} from 'recharts';
import {
  Users, Server, TrendingUp, TrendingDown, Activity,
} from 'lucide-react';
import { useDashboardStats } from '../api/hooks';
import { formatBytes, formatDate } from '../lib/utils';

const AURORA_COLORS = ['#7bf2a8', '#5eeadb', '#b98eff', '#60a5fa', '#f472b6', '#fbbf24', '#94a3b8'];

function StatCard({
  icon: Icon, label, value, sub,
}: {
  icon: React.ComponentType<{ className?: string }>;
  label: string;
  value: string;
  sub?: string | React.ReactNode;
}) {
  return (
    <div
      className="card-aurora rounded-xl p-5 flex flex-col gap-1"
      style={{ background: 'var(--color-polar-850)' }}
    >
      <div className="flex items-center gap-2 mb-1">
        <Icon className="w-4 h-4" />
        <span className="text-xs font-medium uppercase tracking-wider" style={{ color: 'var(--color-frost)' }}>
          {label}
        </span>
      </div>
      <span className="stat-value" style={{ color: 'var(--color-starlight)' }}>
        {value}
      </span>
      {sub && (
        <div className="text-xs mt-0.5" style={{ color: 'var(--color-ice)' }}>
          {sub}
        </div>
      )}
    </div>
  );
}

function NodeStatusBar({ nodes }: { nodes: { status: string }[] }) {
  const online = nodes.filter((n) => n.status === 'online').length;
  const degraded = nodes.filter((n) => n.status === 'degraded').length;
  const offline = nodes.filter((n) => n.status === 'offline').length;
  const total = nodes.length || 1;

  return (
    <div className="flex gap-0.5 h-2 rounded-full overflow-hidden mt-2">
      {online > 0 && (
        <div style={{ width: `${(online / total) * 100}%`, background: 'var(--color-success)' }} />
      )}
      {degraded > 0 && (
        <div style={{ width: `${(degraded / total) * 100}%`, background: 'var(--color-warning)' }} />
      )}
      {offline > 0 && (
        <div style={{ width: `${(offline / total) * 100}%`, background: 'var(--color-danger)' }} />
      )}
    </div>
  );
}

export default function Dashboard() {
  const { data: stats, isLoading } = useDashboardStats();

  if (isLoading || !stats) {
    return (
      <div>
        <h2 className="text-lg font-semibold mb-6" style={{ color: 'var(--color-starlight)' }}>Dashboard</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
          {Array.from({ length: 4 }).map((_, i) => (
            <div key={i} className="skeleton h-28 rounded-xl" />
          ))}
        </div>
        <div className="skeleton h-80 rounded-xl" />
      </div>
    );
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-lg font-semibold m-0" style={{ color: 'var(--color-starlight)' }}>
            Dashboard
          </h2>
          <p className="text-sm mt-1" style={{ color: 'var(--color-frost)' }}>
            Real-time network overview
          </p>
        </div>
        <div className="flex items-center gap-2 text-xs" style={{ color: 'var(--color-ice)' }}>
          <Activity className="w-3 h-3" style={{ color: 'var(--color-success)' }} />
          Live &middot; updates every 15s
        </div>
      </div>

      {/* Stat Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
        <StatCard
          icon={Users}
          label="Total Users"
          value={String(stats.totalUsers)}
          sub={`${stats.activeUsers} active`}
        />
        <StatCard
          icon={Server}
          label="Nodes"
          value={`${stats.onlineNodes}/${stats.totalNodes} online`}
          sub={<NodeStatusBar nodes={[{ status: 'online' }, { status: 'online' }, { status: 'degraded' }, { status: 'offline' }]} />}
          accent="var(--color-aurora-violet)"
        />
        <StatCard
          icon={TrendingDown}
          label="Download"
          value={formatBytes(stats.totalTrafficDown)}
          sub="Total inbound"
          accent="var(--color-aurora-cyan)"
        />
        <StatCard
          icon={TrendingUp}
          label="Upload"
          value={formatBytes(stats.totalTrafficUp)}
          sub="Total outbound"
          accent="var(--color-aurora-blue)"
        />
      </div>

      {/* Traffic Chart */}
      <div
        className="card-aurora rounded-xl p-6 mb-6"
        style={{ background: 'var(--color-polar-850)' }}
      >
        <h3 className="text-sm font-semibold mb-4 uppercase tracking-wider" style={{ color: 'var(--color-frost)' }}>
          Traffic — Last 24 Hours
        </h3>
        <ResponsiveContainer width="100%" height={300}>
          <AreaChart data={stats.trafficHistory}>
            <defs>
              <linearGradient id="gradDownload" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="#5eeadb" stopOpacity={0.2} />
                <stop offset="95%" stopColor="#5eeadb" stopOpacity={0} />
              </linearGradient>
              <linearGradient id="gradUpload" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="#b98eff" stopOpacity={0.2} />
                <stop offset="95%" stopColor="#b98eff" stopOpacity={0} />
              </linearGradient>
            </defs>
            <CartesianGrid stroke="var(--color-polar-700)" strokeDasharray="3 3" />
            <XAxis
              dataKey="timestamp"
              tickFormatter={(t: string) => new Date(t).toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' })}
              tick={{ fill: 'var(--color-ice)', fontSize: 11 }}
              axisLine={{ stroke: 'var(--color-polar-600)' }}
              tickLine={false}
            />
            <YAxis
              tickFormatter={(v: number) => formatBytes(v)}
              tick={{ fill: 'var(--color-ice)', fontSize: 11 }}
              axisLine={{ stroke: 'var(--color-polar-600)' }}
              tickLine={false}
            />
            <Tooltip
              contentStyle={{
                background: 'var(--color-polar-800)',
                border: '1px solid var(--color-polar-600)',
                borderRadius: '0.5rem',
                color: 'var(--color-starlight)',
                fontSize: '0.8rem',
              }}
              labelFormatter={(t: string) => formatDate(t)}
              formatter={(value: number) => [formatBytes(value), '']}
            />
            <Area
              type="monotone"
              dataKey="download"
              stroke="#5eeadb"
              strokeWidth={2}
              fill="url(#gradDownload)"
              name="Download"
            />
            <Area
              type="monotone"
              dataKey="upload"
              stroke="#b98eff"
              strokeWidth={2}
              fill="url(#gradUpload)"
              name="Upload"
            />
          </AreaChart>
        </ResponsiveContainer>
      </div>

      {/* Bottom row: Protocol Pie + User Growth Bar */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
        {/* Protocol Distribution */}
        <div
          className="card-aurora rounded-xl p-6"
          style={{ background: 'var(--color-polar-850)' }}
        >
          <h3 className="text-sm font-semibold mb-4 uppercase tracking-wider" style={{ color: 'var(--color-frost)' }}>
            Protocol Distribution
          </h3>
          <ResponsiveContainer width="100%" height={260}>
            <PieChart>
              <Pie
                data={stats.protocolDistribution}
                dataKey="count"
                nameKey="protocol"
                cx="50%"
                cy="50%"
                outerRadius={100}
                innerRadius={55}
                paddingAngle={2}
              >
                {stats.protocolDistribution.map((_, i) => (
                  <Cell key={i} fill={AURORA_COLORS[i % AURORA_COLORS.length]} />
                ))}
              </Pie>
              <Tooltip
                contentStyle={{
                  background: 'var(--color-polar-800)',
                  border: '1px solid var(--color-polar-600)',
                  borderRadius: '0.5rem',
                  color: 'var(--color-starlight)',
                  fontSize: '0.8rem',
                }}
              />
              <Legend
                formatter={(value: string) => (
                  <span style={{ color: 'var(--color-frost)', fontSize: '0.8rem', textTransform: 'uppercase' }}>
                    {value}
                  </span>
                )}
              />
            </PieChart>
          </ResponsiveContainer>
        </div>

        {/* User Growth */}
        <div
          className="card-aurora rounded-xl p-6"
          style={{ background: 'var(--color-polar-850)' }}
        >
          <h3 className="text-sm font-semibold mb-4 uppercase tracking-wider" style={{ color: 'var(--color-frost)' }}>
            User Growth — Last 30 Days
          </h3>
          <ResponsiveContainer width="100%" height={260}>
            <BarChart data={stats.userGrowth}>
              <CartesianGrid stroke="var(--color-polar-700)" strokeDasharray="3 3" />
              <XAxis
                dataKey="date"
                tickFormatter={(d: string) => new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })}
                tick={{ fill: 'var(--color-ice)', fontSize: 10 }}
                axisLine={{ stroke: 'var(--color-polar-600)' }}
                tickLine={false}
                interval={4}
              />
              <YAxis
                tick={{ fill: 'var(--color-ice)', fontSize: 11 }}
                axisLine={{ stroke: 'var(--color-polar-600)' }}
                tickLine={false}
              />
              <Tooltip
                contentStyle={{
                  background: 'var(--color-polar-800)',
                  border: '1px solid var(--color-polar-600)',
                  borderRadius: '0.5rem',
                  color: 'var(--color-starlight)',
                  fontSize: '0.8rem',
                }}
                labelFormatter={(d: string) => formatDate(d)}
              />
              <Bar dataKey="count" fill="#7bf2a8" radius={[4, 4, 0, 0]} name="Users" />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </div>

      {/* Active Connections */}
      <div
        className="card-aurora rounded-xl p-6"
        style={{ background: 'var(--color-polar-850)' }}
      >
        <div className="flex items-center justify-between">
          <div>
            <h3 className="text-sm font-semibold uppercase tracking-wider mb-1" style={{ color: 'var(--color-frost)' }}>
              Active Connections
            </h3>
            <div className="stat-value" style={{ color: 'var(--color-starlight)' }}>
              {stats.activeConnections.toLocaleString()}
            </div>
          </div>
          <div className="flex items-center gap-2">
            <span className="w-2.5 h-2.5 rounded-full animate-pulse" style={{ background: 'var(--color-success)' }} />
            <span className="text-xs" style={{ color: 'var(--color-ice)' }}>Live</span>
          </div>
        </div>
      </div>
    </div>
  );
}
