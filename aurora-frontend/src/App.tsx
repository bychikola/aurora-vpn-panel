import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Layout } from './components/Layout';
import Dashboard from './pages/Dashboard';
import Users from './pages/Users';
import Nodes from './pages/Nodes';
import Inbounds from './pages/Inbounds';
import Subscriptions from './pages/Subscriptions';
import Settings from './pages/Settings';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 10_000,
      retry: 2,
    },
  },
});

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        {/* Aurora animated gradient strip */}
        <div className="aurora-strip" />

        <Routes>
          <Route element={<Layout />}>
            <Route index element={<Dashboard />} />
            <Route path="users" element={<Users />} />
            <Route path="nodes" element={<Nodes />} />
            <Route path="inbounds" element={<Inbounds />} />
            <Route path="subscriptions" element={<Subscriptions />} />
            <Route path="settings" element={<Settings />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  );
}
