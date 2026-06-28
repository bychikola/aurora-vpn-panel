import { NavLink } from 'react-router-dom';
import {
  LayoutDashboard,
  Users,
  Server,
  Radio,
  Key,
  Settings,
} from 'lucide-react';

const navItems = [
  { to: '/', icon: LayoutDashboard, label: 'Dashboard', end: true },
  { to: '/users', icon: Users, label: 'Users' },
  { to: '/nodes', icon: Server, label: 'Nodes' },
  { to: '/inbounds', icon: Radio, label: 'Inbounds' },
  { to: '/subscriptions', icon: Key, label: 'Subscriptions' },
  { to: '/settings', icon: Settings, label: 'Settings' },
];

export function Sidebar() {
  return (
    <aside
      className="fixed left-0 top-[2px] bottom-0 w-56 flex flex-col border-r z-30"
      style={{
        background: 'var(--color-polar-900)',
        borderColor: 'var(--color-polar-700)',
      }}
    >
      {/* Brand */}
      <div className="px-4 py-5 flex items-center gap-3">
        <div
          className="w-8 h-8 rounded-lg flex items-center justify-center text-sm font-bold"
          style={{
            background: 'linear-gradient(135deg, var(--color-aurora-green), var(--color-aurora-cyan))',
            color: 'var(--color-polar-900)',
            fontFamily: 'var(--font-mono)',
          }}
        >
          A
        </div>
        <div>
          <div
            className="text-sm font-semibold tracking-wide"
            style={{ color: 'var(--color-starlight)' }}
          >
            AURORA
          </div>
          <div
            className="text-[10px] uppercase tracking-wider font-medium"
            style={{ color: 'var(--color-ice)' }}
          >
            VPN Panel
          </div>
        </div>
      </div>

      {/* Nav */}
      <nav className="flex-1 px-3 space-y-0.5">
        {navItems.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            end={item.end}
            className={({ isActive }) =>
              `sidebar-link ${isActive ? 'active' : ''}`
            }
          >
            <item.icon />
            <span>{item.label}</span>
          </NavLink>
        ))}
      </nav>

      {/* Footer */}
      <div
        className="px-4 py-4 border-t text-[11px]"
        style={{
          borderColor: 'var(--color-polar-700)',
          color: 'var(--color-ice)',
        }}
      >
        <div className="flex items-center justify-between">
          <span>v0.1.0-alpha</span>
          <span
            className="inline-flex items-center gap-1"
            style={{ color: 'var(--color-success)' }}
          >
            <span className="w-1.5 h-1.5 rounded-full bg-current inline-block" />
            API Online
          </span>
        </div>
      </div>
    </aside>
  );
}
