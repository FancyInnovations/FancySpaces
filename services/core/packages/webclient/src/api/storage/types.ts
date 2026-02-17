export interface SpaceDatabase {
  name: string;
  created_at: Date;
}

export interface SpaceDatabaseCollection {
  database: string;
  name: string;
  created_at: Date;
  engine: string;
}

export function mapEngineKeyToName(engine?: string): string {
  switch (engine) {
    case 'kv':
      return 'Key-Value store';
    case 'document':
      return 'Document store';
    case 'object':
      return 'Object store';
    case 'analytical':
      return 'Analytical store';
    case 'broker':
      return 'Message broker';
    default:
      return engine || 'Unknown';
  }
}


export interface KVValue {
  key: string;
  value: any;
  type: string;
  ttl?: number;
}
