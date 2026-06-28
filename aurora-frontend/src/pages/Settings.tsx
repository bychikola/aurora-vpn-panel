import { useState } from 'react';
import { Save, Shield, Bell, Globe, Database, Server } from 'lucide-react';

export default function Settings() {
  const [saved, setSaved] = useState(false);

  const handleSave = (e: React.FormEvent) => {
    e.preventDefault();
    setSaved(true);
    setTimeout(() => setSaved(false), 2000);
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-lg font-semibold m-0" style={{ color: 'var(--color-starlight)' }}>Settings</h2>
          <p className="text-sm mt-1" style={{ color: 'var(--color-frost)' }}>
            Panel configuration
          </p>
        </div>
        <button className="btn-primary" onClick={handleSave}>
          <Save className="w-4 h-4" />
          {saved ? 'Saved ✓' : 'Save Changes'}
        </button>
      </div>

      <div className="space-y-6 max-w-2xl">
        {/* General */}
        <section
          className="card-aurora rounded-xl p-6"
          style={{ background: 'var(--color-polar-850)' }}
        >
          <div className="flex items-center gap-2 mb-5">
            <Globe className="w-4 h-4" style={{ color: 'var(--color-aurora-cyan)' }} />
            <h3 className="text-sm font-semibold uppercase tracking-wider m-0" style={{ color: 'var(--color-frost)' }}>
              General
            </h3>
          </div>
          <div className="space-y-4">
            <div>
              <label className="block text-xs font-medium mb-1.5" style={{ color: 'var(--color-frost)' }}>Panel Name</label>
              <input className="input-aurora" defaultValue="AURORA VPN Panel" />
            </div>
            <div>
              <label className="block text-xs font-medium mb-1.5" style={{ color: 'var(--color-frost)' }}>Public Domain</label>
              <input className="input-aurora" defaultValue="aurora.example.com" placeholder="panel.yourdomain.com" />
            </div>
            <div>
              <label className="block text-xs font-medium mb-1.5" style={{ color: 'var(--color-frost)' }}>Language</label>
              <select className="select-aurora" defaultValue="en">
                <option value="en">English</option>
                <option value="ru">Русский</option>
                <option value="fa">فارسی</option>
                <option value="zh">中文</option>
              </select>
            </div>
          </div>
        </section>

        {/* Security */}
        <section
          className="card-aurora rounded-xl p-6"
          style={{ background: 'var(--color-polar-850)' }}
        >
          <div className="flex items-center gap-2 mb-5">
            <Shield className="w-4 h-4" style={{ color: 'var(--color-aurora-violet)' }} />
            <h3 className="text-sm font-semibold uppercase tracking-wider m-0" style={{ color: 'var(--color-frost)' }}>
              Security
            </h3>
          </div>
          <div className="space-y-4">
            <div>
              <label className="block text-xs font-medium mb-1.5" style={{ color: 'var(--color-frost)' }}>JWT Secret</label>
              <input className="input-aurora" type="password" defaultValue="••••••••••••••••" />
            </div>
            <div>
              <label className="block text-xs font-medium mb-1.5" style={{ color: 'var(--color-frost)' }}>Session Timeout (minutes)</label>
              <input className="input-aurora" type="number" defaultValue={60} />
            </div>
            <div className="flex items-center justify-between">
              <div>
                <span className="text-sm" style={{ color: 'var(--color-starlight)' }}>Fail2Ban</span>
                <p className="text-xs mt-0.5" style={{ color: 'var(--color-ice)' }}>Block IPs after repeated failed logins</p>
              </div>
              <button
                className="relative w-10 h-5 rounded-full transition-colors"
                style={{ background: 'var(--color-aurora-green)' }}
                onClick={() => {}}
              >
                <span className="absolute right-0.5 top-0.5 w-4 h-4 bg-white rounded-full shadow" />
              </button>
            </div>
          </div>
        </section>

        {/* Xray Core */}
        <section
          className="card-aurora rounded-xl p-6"
          style={{ background: 'var(--color-polar-850)' }}
        >
          <div className="flex items-center gap-2 mb-5">
            <Server className="w-4 h-4" style={{ color: 'var(--color-aurora-green)' }} />
            <h3 className="text-sm font-semibold uppercase tracking-wider m-0" style={{ color: 'var(--color-frost)' }}>
              Xray Core
            </h3>
          </div>
          <div className="space-y-4">
            <div>
              <label className="block text-xs font-medium mb-1.5" style={{ color: 'var(--color-frost)' }}>gRPC API Host</label>
              <input className="input-aurora" defaultValue="127.0.0.1:10085" />
            </div>
            <div>
              <label className="block text-xs font-medium mb-1.5" style={{ color: 'var(--color-frost)' }}>Config Path</label>
              <input className="input-aurora" defaultValue="/etc/xray/config.json" />
            </div>
          </div>
        </section>

        {/* Logging */}
        <section
          className="card-aurora rounded-xl p-6"
          style={{ background: 'var(--color-polar-850)' }}
        >
          <div className="flex items-center gap-2 mb-5">
            <Database className="w-4 h-4" style={{ color: 'var(--color-ice)' }} />
            <h3 className="text-sm font-semibold uppercase tracking-wider m-0" style={{ color: 'var(--color-frost)' }}>
              Logging & Retention
            </h3>
          </div>
          <div className="space-y-4">
            <div>
              <label className="block text-xs font-medium mb-1.5" style={{ color: 'var(--color-frost)' }}>Log Retention (days)</label>
              <input className="input-aurora" type="number" defaultValue={30} />
            </div>
            <div className="flex items-center justify-between">
              <div>
                <span className="text-sm" style={{ color: 'var(--color-starlight)' }}>Auto-clean expired logs</span>
                <p className="text-xs mt-0.5" style={{ color: 'var(--color-ice)' }}>Purge logs older than retention period</p>
              </div>
              <button
                className="relative w-10 h-5 rounded-full transition-colors"
                style={{ background: 'var(--color-aurora-green)' }}
                onClick={() => {}}
              >
                <span className="absolute right-0.5 top-0.5 w-4 h-4 bg-white rounded-full shadow" />
              </button>
            </div>
          </div>
        </section>

        {/* Bell / Notifications */}
        <section
          className="card-aurora rounded-xl p-6"
          style={{ background: 'var(--color-polar-850)' }}
        >
          <div className="flex items-center gap-2 mb-5">
            <Bell className="w-4 h-4" style={{ color: 'var(--color-warning)' }} />
            <h3 className="text-sm font-semibold uppercase tracking-wider m-0" style={{ color: 'var(--color-frost)' }}>
              Notifications
            </h3>
          </div>
          <div className="space-y-4">
            {['Traffic limit reached', 'User expired', 'Node goes offline', 'New user registration'].map((item) => (
              <div key={item} className="flex items-center justify-between">
                <span className="text-sm" style={{ color: 'var(--color-starlight)' }}>{item}</span>
                <button
                  className="relative w-10 h-5 rounded-full transition-colors"
                  style={{ background: 'var(--color-aurora-green)' }}
                  onClick={() => {}}
                >
                  <span className="absolute right-0.5 top-0.5 w-4 h-4 bg-white rounded-full shadow" />
                </button>
              </div>
            ))}
          </div>
        </section>
      </div>
    </div>
  );
}
