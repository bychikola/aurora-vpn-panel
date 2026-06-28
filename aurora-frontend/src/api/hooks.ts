import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import type { UserFilters, UserFormData } from '../types';
import {
  fetchDashboardStats,
  fetchUsers,
  fetchUser,
  createUser,
  updateUser,
  deleteUser,
  fetchNodes,
  fetchInbounds,
  fetchSubscriptions,
} from './client';

/* ─── Dashboard ─── */
export function useDashboardStats() {
  return useQuery({
    queryKey: ['dashboard', 'stats'],
    queryFn: fetchDashboardStats,
    refetchInterval: 15000,
  });
}

/* ─── Users ─── */
export function useUsers(filters: UserFilters) {
  return useQuery({
    queryKey: ['users', filters],
    queryFn: () => fetchUsers(filters),
  });
}

export function useUser(id: string | null) {
  return useQuery({
    queryKey: ['users', id],
    queryFn: () => fetchUser(id!),
    enabled: !!id,
  });
}

export function useCreateUser() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: UserFormData) =>
      createUser({
        username: data.username,
        email: data.email,
        status: data.status,
        protocols: data.protocols,
        inboundIds: data.inboundIds,
        trafficLimit: data.trafficLimit * 1_000_000_000,
        trafficUsed: 0,
        expireAt: data.expireAt,
        maxIps: data.maxIps,
        concurrentIps: 0,
        subscriptionToken: `sub-${crypto.randomUUID().slice(0, 8)}`,
        notes: data.notes,
      }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['users'] });
      qc.invalidateQueries({ queryKey: ['dashboard'] });
    },
  });
}

export function useUpdateUser() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<UserFormData> }) =>
      updateUser(id, data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['users'] });
    },
  });
}

export function useDeleteUser() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => deleteUser(id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['users'] });
      qc.invalidateQueries({ queryKey: ['dashboard'] });
    },
  });
}

/* ─── Nodes ─── */
export function useNodes() {
  return useQuery({
    queryKey: ['nodes'],
    queryFn: fetchNodes,
    refetchInterval: 30000,
  });
}

/* ─── Inbounds ─── */
export function useInbounds() {
  return useQuery({
    queryKey: ['inbounds'],
    queryFn: fetchInbounds,
    refetchInterval: 30000,
  });
}

/* ─── Subscriptions ─── */
export function useSubscriptions() {
  return useQuery({
    queryKey: ['subscriptions'],
    queryFn: fetchSubscriptions,
    refetchInterval: 30000,
  });
}
