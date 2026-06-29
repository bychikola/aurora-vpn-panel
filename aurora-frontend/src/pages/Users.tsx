import { useState } from 'react';
import { Plus, Search, Trash2, Edit3, QrCode } from 'lucide-react';
import { useUsers, useDeleteUser } from '../api/hooks';
import { formatBytes, formatDate } from '../lib/utils';
import type { User, UserFilters, Protocol } from '../types';
import { UserForm } from './UserForm';

const PROTOCOL_LABELS: Record<Protocol, string> = {
  vless: 'VLESS',
  vmess: 'VMess',
  trojan: 'Trojan',
  shadowsocks: 'SS',
  'shadowsocks-2022': 'SS-2022',
  hysteria2: 'Hysteria2',
  'tuic-v5': 'TUIC v5',
};

export default function Users() {
  const [filters, setFilters] = useState<UserFilters>({ page: 1, pageSize: 15 });
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState<User['status'] | 'all'>('all');
  const [protocolFilter, setProtocolFilter] = useState<Protocol | 'all'>('all');
  const [showForm, setShowForm] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);

  const appliedFilters = { ...filters, search: search || undefined, status: statusFilter, protocol: protocolFilter };
  const { data, isLoading } = useUsers(appliedFilters);
  const deleteUser = useDeleteUser();

  const handleSearch = () => {
    setFilters((f) => ({ ...f, page: 1 }));
  };

  const handleDelete = async (id: string) => {
    if (!window.confirm('Delete this user? All associated data will be lost.')) return;
    deleteUser.mutate(id);
  };

  const handleEdit = (user: User) => {
    setEditingUser(user);
    setShowForm(true);
  };

  const handleCreate = () => {
    setEditingUser(null);
    setShowForm(true);
  };

  const closeForm = () => {
    setShowForm(false);
    setEditingUser(null);
  };

  const trafficPercent = (used: number, limit: number) => {
    if (limit === 0) return 0;
    return Math.min((used / limit) * 100, 100);
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-lg font-semibold m-0" style={{ color: 'var(--color-starlight)' }}>Users</h2>
          <p className="text-sm mt-1" style={{ color: 'var(--color-frost)' }}>
            {data ? `${data.total} total` : 'Loading...'}
          </p>
        </div>
        <button className="btn-primary" onClick={handleCreate}>
          <Plus className="w-4 h-4" />
          New User
        </button>
      </div>

      {/* Filters */}
      <div className="flex flex-wrap gap-3 mb-4">
        <div className="relative flex-1 min-w-[200px]">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4" style={{ color: 'var(--color-ice)' }} />
          <input
            className="input-aurora pl-9"
            placeholder="Search by username or email..."
            value={search}
            onChange={(e) => { setSearch(e.target.value); }}
            onKeyDown={(e) => { if (e.key === 'Enter') handleSearch(); }}
          />
        </div>
        <select
          className="select-aurora"
          value={statusFilter}
          onChange={(e) => { setStatusFilter(e.target.value as User['status'] | 'all'); setFilters((f) => ({ ...f, page: 1 })); }}
        >
          <option value="all">All Status</option>
          <option value="active">Active</option>
          <option value="disabled">Disabled</option>
          <option value="expired">Expired</option>
        </select>
        <select
          className="select-aurora"
          value={protocolFilter}
          onChange={(e) => { setProtocolFilter(e.target.value as Protocol | 'all'); setFilters((f) => ({ ...f, page: 1 })); }}
        >
          <option value="all">All Protocols</option>
          {Object.entries(PROTOCOL_LABELS).map(([k, v]) => (
            <option key={k} value={k}>{v}</option>
          ))}
        </select>
      </div>

      {/* Table */}
      <div
        className="rounded-xl overflow-hidden"
        style={{ background: 'var(--color-polar-850)', border: '1px solid var(--color-polar-700)' }}
      >
        {isLoading ? (
          <div className="p-8 space-y-3">
            {Array.from({ length: 8 }).map((_, i) => (
              <div key={i} className="skeleton h-10 rounded" />
            ))}
          </div>
        ) : (
          <table className="table-aurora">
            <thead>
              <tr>
                <th>User</th>
                <th>Protocols</th>
                <th>Traffic</th>
                <th>Expires</th>
                <th>IPs</th>
                <th>Status</th>
                <th style={{ width: 100 }}>Actions</th>
              </tr>
            </thead>
            <tbody>
              {data?.data.map((user) => (
                <tr key={user.id}>
                  <td>
                    <div className="font-medium" style={{ color: 'var(--color-starlight)' }}>{user.username}</div>
                    <div className="text-xs" style={{ color: 'var(--color-ice)' }}>{user.email}</div>
                  </td>
                  <td>
                    <div className="flex flex-wrap gap-1">
                      {user.protocols.map((p) => (
                        <span key={p} className="badge badge-info">{PROTOCOL_LABELS[p]}</span>
                      ))}
                    </div>
                  </td>
                  <td>
                    <div className="flex items-center gap-2">
                      <div className="flex-1 min-w-[80px] h-1.5 rounded-full" style={{ background: 'var(--color-polar-700)' }}>
                        <div
                          className="h-full rounded-full transition-all"
                          style={{
                            width: `${trafficPercent(user.trafficUsed, user.trafficLimit)}%`,
                            background:
                              trafficPercent(user.trafficUsed, user.trafficLimit) > 90
                                ? 'var(--color-danger)'
                                : trafficPercent(user.trafficUsed, user.trafficLimit) > 70
                                  ? 'var(--color-warning)'
                                  : 'var(--color-aurora-green)',
                          }}
                        />
                      </div>
                      <span className="text-xs font-mono-num whitespace-nowrap" style={{ color: 'var(--color-frost)' }}>
                        {formatBytes(user.trafficUsed)} / {formatBytes(user.trafficLimit)}
                      </span>
                    </div>
                  </td>
                  <td>
                    <span className="text-sm font-mono-num" style={{ color: new Date(user.expireAt) < new Date() ? 'var(--color-danger)' : 'var(--color-starlight)' }}>
                      {formatDate(user.expireAt)}
                    </span>
                  </td>
                  <td>
                    <span className="font-mono-num text-sm" style={{ color: 'var(--color-frost)' }}>
                      {user.concurrentIps}/{user.maxIps}
                    </span>
                  </td>
                  <td>
                    <span className={`badge badge-${user.status === 'active' ? 'success' : user.status === 'disabled' ? 'warning' : 'danger'}`}>
                      {user.status}
                    </span>
                  </td>
                  <td>
                    <div className="flex items-center gap-1">
                      <button
                        className="p-1.5 rounded-md hover:bg-white/5 transition-colors"
                        onClick={() => handleEdit(user)}
                        title="Edit"
                      >
                        <Edit3 className="w-3.5 h-3.5" style={{ color: 'var(--color-frost)' }} />
                      </button>
                      <button
                        className="p-1.5 rounded-md hover:bg-white/5 transition-colors"
                        onClick={() => window.alert(`QR Code for ${user.username} — will open in a dialog`)}
                        title="QR Code"
                      >
                        <QrCode className="w-3.5 h-3.5" style={{ color: 'var(--color-frost)' }} />
                      </button>
                      <button
                        className="p-1.5 rounded-md hover:bg-white/5 transition-colors"
                        onClick={() => handleDelete(user.id)}
                        title="Delete"
                      >
                        <Trash2 className="w-3.5 h-3.5" style={{ color: 'var(--color-danger)' }} />
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
              {data?.data.length === 0 && (
                <tr>
                  <td colSpan={7} className="text-center py-12" style={{ color: 'var(--color-ice)' }}>
                    No users found
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        )}
      </div>

      {/* Pagination */}
      {data && data.total > filters.pageSize && (
        <div className="flex items-center justify-between mt-4">
          <span className="text-xs" style={{ color: 'var(--color-ice)' }}>
            Page {data.page} of {Math.ceil(data.total / data.pageSize)}
          </span>
          <div className="flex gap-2">
            <button
              className="btn-ghost text-xs"
              disabled={filters.page <= 1}
              onClick={() => setFilters((f) => ({ ...f, page: f.page - 1 }))}
            >
              Previous
            </button>
            <button
              className="btn-ghost text-xs"
              disabled={filters.page >= Math.ceil(data.total / data.pageSize)}
              onClick={() => setFilters((f) => ({ ...f, page: f.page + 1 }))}
            >
              Next
            </button>
          </div>
        </div>
      )}

      {/* User Form Modal */}
      {showForm && <UserForm user={editingUser} onClose={closeForm} />}
    </div>
  );
}
