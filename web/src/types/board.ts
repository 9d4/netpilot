export interface BoardOfBoards {
  uuid: string
  name: string
  host: string
  port: string
}

export interface Board {
  created_at: string
  updated_at: string
  uuid: string
  name: string
  host: string
  port: string
  insecure_skip_verify: boolean
  user: string
}

export interface BoardInfo {
  'architecture-name': string
  'bad-blocks': string
  'board-name': string
  'build-time': string
  cpu: string
  'cpu-count': string
  'cpu-frequency': string
  'cpu-load': string
  'factory-software': string
  'free-hdd-space': string
  'free-memory': string
  platform: string
  'total-hdd-space': string
  'total-memory': string
  uptime: string
  version: string
  'write-sect-since-reboot': string
  'write-sect-total': string
}
