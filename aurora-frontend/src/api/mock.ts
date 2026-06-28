import type {
  DashboardStats,
  Inbound,
  Node,
  Subscription,
  User,
} from '../types';

const NOW = new Date();
const DAY = 86400000;

/* ─── Inbounds ─── */
export const mockInbounds: Inbound[] = [
  {
    id: 'inb-001', tag: 'VLESS + XTLS', protocol: 'vless', port: 443,
    listen: '0.0.0.0', transport: 'tcp', security: 'xtls-vision',
    nodeId: 'node-001', nodeName: 'Frankfurt-1', enable: true,
    userCount: 142, upload: 1_200_000_000_000, download: 8_500_000_000_000,
    settings: {}, streamSettings: {},
    createdAt: new Date(NOW.getTime() - 30 * DAY).toISOString(),
    updatedAt: new Date(NOW.getTime() - 2 * DAY).toISOString(),
  },
  {
    id: 'inb-002', tag: 'VMess + WS', protocol: 'vmess', port: 8080,
    listen: '0.0.0.0', transport: 'ws', security: 'tls',
    nodeId: 'node-001', nodeName: 'Frankfurt-1', enable: true,
    userCount: 89, upload: 650_000_000_000, download: 3_200_000_000_000,
    settings: {}, streamSettings: {},
    createdAt: new Date(NOW.getTime() - 25 * DAY).toISOString(),
    updatedAt: new Date(NOW.getTime() - 5 * DAY).toISOString(),
  },
  {
    id: 'inb-003', tag: 'Trojan + gRPC', protocol: 'trojan', port: 443,
    listen: '0.0.0.0', transport: 'grpc', security: 'tls',
    nodeId: 'node-002', nodeName: 'Amsterdam-1', enable: true,
    userCount: 56, upload: 340_000_000_000, download: 1_800_000_000_000,
    settings: {}, streamSettings: {},
    createdAt: new Date(NOW.getTime() - 20 * DAY).toISOString(),
    updatedAt: new Date(NOW.getTime() - 1 * DAY).toISOString(),
  },
  {
    id: 'inb-004', tag: 'Shadowsocks-2022', protocol: 'shadowsocks-2022', port: 8388,
    listen: '0.0.0.0', transport: 'tcp', security: 'none',
    nodeId: 'node-002', nodeName: 'Amsterdam-1', enable: true,
    userCount: 34, upload: 120_000_000_000, download: 590_000_000_000,
    settings: {}, streamSettings: {},
    createdAt: new Date(NOW.getTime() - 15 * DAY).toISOString(),
    updatedAt: new Date(NOW.getTime() - 10 * DAY).toISOString(),
  },
  {
    id: 'inb-005', tag: 'Hysteria2', protocol: 'hysteria2', port: 443,
    listen: '0.0.0.0', transport: 'quic', security: 'tls',
    nodeId: 'node-003', nodeName: 'Singapore-1', enable: true,
    userCount: 78, upload: 890_000_000_000, download: 4_100_000_000_000,
    settings: {}, streamSettings: {},
    createdAt: new Date(NOW.getTime() - 10 * DAY).toISOString(),
    updatedAt: new Date(NOW.getTime() - 3 * DAY).toISOString(),
  },
  {
    id: 'inb-006', tag: 'VLESS + Reality', protocol: 'vless', port: 443,
    listen: '0.0.0.0', transport: 'tcp', security: 'reality',
    nodeId: 'node-003', nodeName: 'Singapore-1', enable: false,
    userCount: 0, upload: 0, download: 0,
    settings: {}, streamSettings: {},
    createdAt: new Date(NOW.getTime() - 5 * DAY).toISOString(),
    updatedAt: new Date(NOW.getTime() - 5 * DAY).toISOString(),
  },
];

/* ─── Users ─── */
export const mockUsers: User[] = Array.from({ length: 25 }, (_, i) => {
  const protocols = [
    ['vless', 'vmess'],
    ['trojan'],
    ['shadowsocks-2022', 'hysteria2'],
    ['vless', 'vmess', 'trojan'],
    ['hysteria2'],
    ['vless', 'tuic-v5'],
    ['vmess', 'shadowsocks'],
    ['vless'],
  ][i % 8] as User['protocols'];

  const statuses: User['status'][] = ['active', 'active', 'active', 'active', 'disabled', 'expired'];
  const status = statuses[i % statuses.length];
  const trafficLimit = [50, 100, 200, 500, 1000, 2000][i % 6] * 1_000_000_000;
  const trafficUsed = Math.floor(Math.random() * trafficLimit * 0.95);

  return {
    id: `usr-${String(i + 1).padStart(3, '0')}`,
    username: `user_${i + 1}`,
    email: `user${i + 1}@example.com`,
    status,
    protocols,
    inboundIds: [mockInbounds[i % mockInbounds.length].id],
    trafficLimit,
    trafficUsed,
    expireAt: new Date(NOW.getTime() + (i % 4 === 0 ? -5 * DAY : (i + 1) * 30 * DAY)).toISOString(),
    maxIps: [1, 2, 3, 5][i % 4],
    concurrentIps: Math.floor(Math.random() * 4),
    subscriptionToken: `sub-${crypto.randomUUID().slice(0, 8)}`,
    notes: i % 7 === 0 ? 'Priority customer' : '',
    createdAt: new Date(NOW.getTime() - (i + 1) * 15 * DAY).toISOString(),
    updatedAt: new Date(NOW.getTime() - i * 2 * DAY).toISOString(),
    lastSeenAt: status === 'active'
      ? new Date(NOW.getTime() - Math.random() * 3600000).toISOString()
      : null,
  };
});

/* ─── Nodes ─── */
export const mockNodes: Node[] = [
  {
    id: 'node-001', name: 'Frankfurt-1', host: '185.220.101.1', port: 443,
    apiPort: 10085, status: 'online', version: 'v25.5.1',
    cpuPercent: 34, memoryPercent: 62, diskPercent: 45,
    uplinkSpeed: 12_500_000, downlinkSpeed: 98_000_000,
    uplinkTotal: 8_500_000_000_000, downlinkTotal: 52_000_000_000_000,
    userCount: 231, inboundCount: 2, location: '🇩🇪 Frankfurt, DE',
    lastPing: new Date(NOW.getTime() - 3000).toISOString(),
    createdAt: new Date(NOW.getTime() - 60 * DAY).toISOString(),
  },
  {
    id: 'node-002', name: 'Amsterdam-1', host: '45.92.38.1', port: 443,
    apiPort: 10085, status: 'online', version: 'v25.5.1',
    cpuPercent: 22, memoryPercent: 45, diskPercent: 38,
    uplinkSpeed: 8_200_000, downlinkSpeed: 54_000_000,
    uplinkTotal: 4_700_000_000_000, downlinkTotal: 28_000_000_000_000,
    userCount: 90, inboundCount: 2, location: '🇳🇱 Amsterdam, NL',
    lastPing: new Date(NOW.getTime() - 5000).toISOString(),
    createdAt: new Date(NOW.getTime() - 45 * DAY).toISOString(),
  },
  {
    id: 'node-003', name: 'Singapore-1', host: '103.142.140.1', port: 443,
    apiPort: 10085, status: 'degraded', version: 'v25.4.0',
    cpuPercent: 78, memoryPercent: 88, diskPercent: 72,
    uplinkSpeed: 15_000_000, downlinkSpeed: 72_000_000,
    uplinkTotal: 6_200_000_000_000, downlinkTotal: 35_000_000_000_000,
    userCount: 78, inboundCount: 2, location: '🇸🇬 Singapore',
    lastPing: new Date(NOW.getTime() - 45000).toISOString(),
    createdAt: new Date(NOW.getTime() - 30 * DAY).toISOString(),
  },
  {
    id: 'node-004', name: 'Tokyo-1', host: '153.127.96.1', port: 443,
    apiPort: 10085, status: 'offline', version: 'v25.3.0',
    cpuPercent: 0, memoryPercent: 0, diskPercent: 0,
    uplinkSpeed: 0, downlinkSpeed: 0,
    uplinkTotal: 1_200_000_000_000, downlinkTotal: 8_000_000_000_000,
    userCount: 45, inboundCount: 1, location: '🇯🇵 Tokyo, JP',
    lastPing: new Date(NOW.getTime() - 3600000).toISOString(),
    createdAt: new Date(NOW.getTime() - 50 * DAY).toISOString(),
  },
];

/* ─── Subscriptions ─── */
export const mockSubscriptions: Subscription[] = mockUsers.slice(0, 10).map((u) => ({
  id: `sub-${u.id}`,
  userId: u.id,
  username: u.username,
  token: u.subscriptionToken,
  url: `https://aurora.example.com/sub/${u.subscriptionToken}`,
  format: (['base64', 'clash', 'sing-box'] as const)[Math.floor(Math.random() * 3)],
  enabled: u.status === 'active',
  lastRequestAt: u.status === 'active'
    ? new Date(NOW.getTime() - Math.random() * 86400000).toISOString()
    : null,
  userAgent: 'v2rayNG/1.8.5',
  createdAt: u.createdAt,
}));

/* ─── Dashboard Stats ─── */
export const mockDashboardStats: DashboardStats = {
  totalUsers: mockUsers.length,
  activeUsers: mockUsers.filter((u) => u.status === 'active').length,
  totalNodes: mockNodes.length,
  onlineNodes: mockNodes.filter((n) => n.status === 'online').length,
  totalTrafficUp: 21_570_000_000_000,
  totalTrafficDown: 127_800_000_000_000,
  activeConnections: 847,
  protocolDistribution: [
    { protocol: 'vless', count: 142 },
    { protocol: 'vmess', count: 98 },
    { protocol: 'trojan', count: 56 },
    { protocol: 'hysteria2', count: 78 },
    { protocol: 'shadowsocks-2022', count: 34 },
    { protocol: 'shadowsocks', count: 22 },
    { protocol: 'tuic-v5', count: 15 },
  ],
  trafficHistory: Array.from({ length: 24 }, (_, i) => ({
    timestamp: new Date(NOW.getTime() - (23 - i) * 3600000).toISOString(),
    upload: 200_000_000 + Math.random() * 800_000_000,
    download: 1_000_000_000 + Math.random() * 4_000_000_000,
  })),
  userGrowth: Array.from({ length: 30 }, (_, i) => ({
    date: new Date(NOW.getTime() - (29 - i) * DAY).toISOString(),
    count: 180 + Math.floor(i * 2.3) + Math.floor(Math.random() * 10),
  })),
};
