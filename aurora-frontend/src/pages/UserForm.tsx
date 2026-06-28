import { useState, useEffect } from 'react';
import { X, Save } from 'lucide-react';
import { useCreateUser, useUpdateUser } from '../api/hooks';
import type { User, UserFormData, Protocol } from '../types';

const PROTOCOLS: { value: Protocol; label: string }[] = [
  { value: 'vless', label: 'VLESS' },
  { value: 'vmess', label: 'VMess' },
  { value: 'trojan', label: 'Trojan' },
  { value: 'shadowsocks', label: 'Shadowsocks' },
  { value: 'shadowsocks-2022', label: 'Shadowsocks-2022' },
  { value: 'hysteria2', label: 'Hysteria2' },
  { value: 'tuic-v5', label: 'TUIC v5' },
];

const emptyForm: UserFormData = {
  username: '',
  email: '',
  status: 'active',
  protocols: ['vless'],
  inboundIds: [],
  trafficLimit: 100,
  expireAt: new Date(Date.now() + 30 * 86400000).toISOString().slice(0, 10),
  maxIps: 3,
  notes: '',
};

interface Props {
  user: User | null;
  onClose: () => void;
}

export function UserForm({ user, onClose }: Props) {
  const [form, setForm] = useState<UserFormData>(emptyForm);
  const createUser = useCreateUser();
  const updateUser = useUpdateUser();
  const isEdit = !!user;

  useEffect(() => {
    if (user) {
      setForm({
        username: user.username,
        email: user.email,
        status: user.status,
        protocols: user.protocols,
        inboundIds: user.inboundIds,
        trafficLimit: Math.round(user.trafficLimit / 1_000_000_000),
        expireAt: user.expireAt.slice(0, 10),
        maxIps: user.maxIps,
        notes: user.notes,
      });
    }
  }, [user]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (isEdit) {
      updateUser.mutate({ id: user!.id, data: form }, { onSuccess: onClose });
    } else {
      createUser.mutate(form, { onSuccess: onClose });
    }
  };

  const toggleProtocol = (p: Protocol) => {
    setForm((f) => ({
      ...f,
      protocols: f.protocols.includes(p)
        ? f.protocols.filter((x) => x !== p)
        : [...f.protocols, p],
    }));
  };

  const isPending = createUser.isPending || updateUser.isPending;

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <div className="flex items-center justify-between mb-5">
          <h3 className="text-base font-semibold m-0" style={{ color: 'var(--color-starlight)' }}>
            {isEdit ? 'Edit User' : 'Create User'}
          </h3>
          <button onClick={onClose} className="p-1 rounded-md hover:bg-white/5 transition-colors">
            <X className="w-4 h-4" style={{ color: 'var(--color-frost)' }} />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-xs font-medium mb-1.5" style={{ color: 'var(--color-frost)' }}>
                Username
              </label>
              <input
                className="input-aurora"
                value={form.username}
                onChange={(e) => setForm((f) => ({ ...f, username: e.target.value }))}
                required
                placeholder="johndoe"
              />
            </div>
            <div>
              <label className="block text-xs font-medium mb-1.5" style={{ color: 'var(--color-frost)' }}>
                Email
              </label>
              <input
                className="input-aurora"
                type="email"
                value={form.email}
                onChange={(e) => setForm((f) => ({ ...f, email: e.target.value }))}
                placeholder="john@example.com"
              />
            </div>
          </div>

          <div className="grid grid-cols-3 gap-4">
            <div>
              <label className="block text-xs font-medium mb-1.5" style={{ color: 'var(--color-frost)' }}>
                Traffic Limit (GB)
              </label>
              <input
                className="input-aurora"
                type="number"
                min={1}
                value={form.trafficLimit}
                onChange={(e) => setForm((f) => ({ ...f, trafficLimit: Number(e.target.value) }))}
              />
            </div>
            <div>
              <label className="block text-xs font-medium mb-1.5" style={{ color: 'var(--color-frost)' }}>
                Expires
              </label>
              <input
                className="input-aurora"
                type="date"
                value={form.expireAt}
                onChange={(e) => setForm((f) => ({ ...f, expireAt: e.target.value }))}
              />
            </div>
            <div>
              <label className="block text-xs font-medium mb-1.5" style={{ color: 'var(--color-frost)' }}>
                Max IPs
              </label>
              <input
                className="input-aurora"
                type="number"
                min={1}
                max={10}
                value={form.maxIps}
                onChange={(e) => setForm((f) => ({ ...f, maxIps: Number(e.target.value) }))}
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-xs font-medium mb-1.5" style={{ color: 'var(--color-frost)' }}>
                Status
              </label>
              <select
                className="select-aurora"
                value={form.status}
                onChange={(e) => setForm((f) => ({ ...f, status: e.target.value as UserFormData['status'] }))}
              >
                <option value="active">Active</option>
                <option value="disabled">Disabled</option>
                <option value="expired">Expired</option>
              </select>
            </div>
          </div>

          <div>
            <label className="block text-xs font-medium mb-2" style={{ color: 'var(--color-frost)' }}>
              Protocols
            </label>
            <div className="flex flex-wrap gap-2">
              {PROTOCOLS.map(({ value, label }) => {
                const active = form.protocols.includes(value);
                return (
                  <button
                    key={value}
                    type="button"
                    className={`text-xs px-3 py-1.5 rounded-full font-medium transition-all ${
                      active
                        ? 'border-transparent'
                        : 'border'
                    }`}
                    style={
                      active
                        ? {
                            background: 'linear-gradient(135deg, var(--color-aurora-green), var(--color-aurora-cyan))',
                            color: 'var(--color-polar-900)',
                          }
                        : {
                            borderColor: 'var(--color-polar-600)',
                            color: 'var(--color-frost)',
                          }
                    }
                    onClick={() => toggleProtocol(value)}
                  >
                    {label}
                  </button>
                );
              })}
            </div>
          </div>

          <div>
            <label className="block text-xs font-medium mb-1.5" style={{ color: 'var(--color-frost)' }}>
              Notes
            </label>
            <textarea
              className="input-aurora"
              rows={2}
              value={form.notes}
              onChange={(e) => setForm((f) => ({ ...f, notes: e.target.value }))}
              placeholder="Optional notes..."
            />
          </div>

          <div className="flex justify-end gap-3 pt-2">
            <button type="button" className="btn-ghost" onClick={onClose}>
              Cancel
            </button>
            <button type="submit" className="btn-primary" disabled={isPending}>
              <Save className="w-4 h-4" />
              {isPending ? 'Saving...' : isEdit ? 'Update User' : 'Create User'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
