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
      return 'Key-Value Store';
    case 'document':
      return 'Document Store';
    case 'object':
      return 'Object Store';
    case 'analytical':
      return 'Analytical Store';
    case 'broker':
      return 'Message Broker';
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
