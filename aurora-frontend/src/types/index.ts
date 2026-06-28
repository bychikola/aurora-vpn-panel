/* ─── Protocol & Transport ─── */
export type Protocol = 'vless' | 'vmess' | 'trojan' | 'shadowsocks' | 'shadowsocks-2022' | 'hysteria2' | 'tuic-v5';
export type Transport = 'tcp' | 'http' | 'ws' | 'grpc' | 'quic';
export type Security = 'tls' | 'reality' | 'xtls-vision' | 'none';

export interface Inbound {
  id: string;
  tag: string;
  protocol: Protocol;
  port: number;
  listen: string;
  transport: Transport;
  security: Security;
  nodeId: string;
  nodeName: string;
  enable: boolean;
  userCount: number;
  upload: number;
  download: number;
  settings: Record<string, unknown>;
  streamSettings: Record<string, unknown>;
  createdAt: string;
  updatedAt: string;
}

/* ─── User ─── */
export interface User {
  id: string;
  username: string;
  email: string;
  status: 'active' | 'disabled' | 'expired';
  protocols: Protocol[];
  inboundIds: string[];
  trafficLimit: number; // bytes
  trafficUsed: number; // bytes
  expireAt: string;
  maxIps: number;
  concurrentIps: number;
  subscriptionToken: string;
  notes: string;
  createdAt: string;
  updatedAt: string;
  lastSeenAt: string | null;
}

export interface UserFormData {
  username: string;
  email: string;
  status: 'active' | 'disabled' | 'expired';
  protocols: Protocol[];
  inboundIds: string[];
  trafficLimit: number; // GB input, converted to bytes
  expireAt: string;
  maxIps: number;
  notes: string;
}

/* ─── Node ─── */
export interface Node {
  id: string;
  name: string;
  host: string;
  port: number;
  apiPort: number;
  status: 'online' | 'offline' | 'degraded';
  version: string;
  cpuPercent: number;
  memoryPercent: number;
  diskPercent: number;
  uplinkSpeed: number; // bytes/s
  downlinkSpeed: number; // bytes/s
  uplinkTotal: number; // bytes
  downlinkTotal: number; // bytes
  userCount: number;
  inboundCount: number;
  location: string;
  lastPing: string;
  createdAt: string;
}

/* ─── Subscription ─── */
export interface Subscription {
  id: string;
  userId: string;
  username: string;
  token: string;
  url: string;
  format: 'base64' | 'clash' | 'sing-box';
  enabled: boolean;
  lastRequestAt: string | null;
  userAgent: string | null;
  createdAt: string;
}

/* ─── Dashboard Stats ─── */
export interface DashboardStats {
  totalUsers: number;
  activeUsers: number;
  totalNodes: number;
  onlineNodes: number;
  totalTrafficUp: number;
  totalTrafficDown: number;
  activeConnections: number;
  protocolDistribution: { protocol: Protocol; count: number }[];
  trafficHistory: { timestamp: string; upload: number; download: number }[];
  userGrowth: { date: string; count: number }[];
}

/* ─── API Responses ─── */
export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  pageSize: number;
}

export interface ApiError {
  code: string;
  message: string;
}

/* ─── Filters ─── */
export interface UserFilters {
  search?: string;
  status?: User['status'] | 'all';
  protocol?: Protocol | 'all';
  page: number;
  pageSize: number;
}
