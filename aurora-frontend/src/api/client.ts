import type {
  DashboardStats,
  Inbound,
  Node,
  PaginatedResponse,
  Subscription,
  User,
  UserFilters,
} from '../types';
import {
  mockDashboardStats,
  mockInbounds,
  mockNodes,
  mockSubscriptions,
  mockUsers,
} from './mock';

/* Simulated network delay: 200–600ms */
function delay<T>(data: T): Promise<T> {
  const ms = 200 + Math.random() * 400;
  return new Promise((resolve) => setTimeout(() => resolve(data), ms));
}

/* ─── Dashboard ─── */
export async function fetchDashboardStats(): Promise<DashboardStats> {
  return delay({ ...mockDashboardStats });
}

/* ─── Users ─── */
export async function fetchUsers(filters: UserFilters): Promise<PaginatedResponse<User>> {
  let filtered = [...mockUsers];

  if (filters.search) {
    const q = filters.search.toLowerCase();
    filtered = filtered.filter(
      (u) => u.username.toLowerCase().includes(q) || u.email.toLowerCase().includes(q),
    );
  }
  if (filters.status && filters.status !== 'all') {
    filtered = filtered.filter((u) => u.status === filters.status);
  }
  if (filters.protocol && filters.protocol !== 'all') {
    filtered = filtered.filter((u) => u.protocols.includes(filters.protocol));
  }

  const total = filtered.length;
  const start = (filters.page - 1) * filters.pageSize;
  const data = filtered.slice(start, start + filters.pageSize);

  return delay({ data, total, page: filters.page, pageSize: filters.pageSize });
}

export async function fetchUser(id: string): Promise<User | null> {
  return delay(mockUsers.find((u) => u.id === id) ?? null);
}

export async function createUser(user: Omit<User, 'id' | 'createdAt' | 'updatedAt' | 'lastSeenAt'>): Promise<User> {
  const newUser: User = {
    ...user,
    id: `usr-${Date.now()}`,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
    lastSeenAt: null,
  };
  mockUsers.push(newUser);
  return delay(newUser);
}

export async function updateUser(id: string, data: Partial<User>): Promise<User> {
  const idx = mockUsers.findIndex((u) => u.id === id);
  if (idx === -1) throw new Error('User not found');
  mockUsers[idx] = { ...mockUsers[idx], ...data, updatedAt: new Date().toISOString() };
  return delay(mockUsers[idx]);
}

export async function deleteUser(id: string): Promise<void> {
  const idx = mockUsers.findIndex((u) => u.id === id);
  if (idx !== -1) mockUsers.splice(idx, 1);
  return delay(undefined);
}

/* ─── Nodes ─── */
export async function fetchNodes(): Promise<Node[]> {
  return delay([...mockNodes]);
}

export async function fetchNode(id: string): Promise<Node | null> {
  return delay(mockNodes.find((n) => n.id === id) ?? null);
}

/* ─── Inbounds ─── */
export async function fetchInbounds(): Promise<Inbound[]> {
  return delay([...mockInbounds]);
}

export async function fetchInbound(id: string): Promise<Inbound | null> {
  return delay(mockInbounds.find((i) => i.id === id) ?? null);
}

/* ─── Subscriptions ─── */
export async function fetchSubscriptions(): Promise<Subscription[]> {
  return delay([...mockSubscriptions]);
}
