import { Outlet } from 'react-router-dom';
import { Sidebar } from './Sidebar';

export function Layout() {
  return (
    <div className="flex min-h-screen" style={{ background: 'var(--color-polar-900)' }}>
      <Sidebar />
      <main className="flex-1 ml-56 pt-[2px]">
        <div className="page-enter p-6 max-w-[1440px]">
          <Outlet />
        </div>
      </main>
    </div>
  );
}
