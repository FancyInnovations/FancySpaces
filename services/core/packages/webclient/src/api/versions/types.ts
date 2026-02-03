export interface SpaceVersion {
  space_id: string;
  id: string;
  name: string;
  platform: string;
  channel: string;
  published_at: Date;
  changelog: string;
  supported_platform_versions: string[];
  files: SpaceVersionFile[];
}

export interface SpaceVersionFile {
  name: string;
  url: string;
  size: number;
}

export function mapPlatformToDisplayname(name?: string): string {
  if (!name) return 'Unknown';

  switch (name.toLowerCase()) {
    case 'bukkit':
      return 'Bukkit';
    case 'spigot':
      return 'Spigot';
    case 'paper':
      return 'Paper';
    case 'purpur':
      return 'Purpur';
    case 'folia':
      return 'Folia';
    case 'bungeecord':
      return 'BungeeCord';
    case 'waterfall':
      return 'Waterfall';
    case 'velocity':
      return 'Velocity';
    case 'fabric':
      return 'Fabric';
    case 'forge':
      return 'Forge';
    case 'quilt':
      return 'Quilt';
    case 'liteloader':
      return 'LiteLoader';
    case 'hytale_plugin':
      return 'Hytale Plugin';
    case 'executable':
      return 'Executable';

    default:
      return name;
  }
}
