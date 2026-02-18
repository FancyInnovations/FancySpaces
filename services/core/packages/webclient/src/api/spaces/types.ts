export interface Space {
  id: string;
  slug: string;
  title: string;
  summary: string;
  description: string;
  categories: string[];
  links: SpaceLink[];
  icon_url: string;
  status: string;
  created_at: Date;

  creator: string;
  members: SpaceMember[];

  issue_settings: SpaceIssueSettings;
  release_settings: SpaceReleaseSettings;
  maven_repository_settings: MavenRepositorySettings;
  storage_settings: StorageSettings;
  analytics_settings: AnalyticsSettings;
  secrets_settings: SecretsSettings;
}

export interface SpaceIssueSettings {
  enabled: boolean;
}

export interface SpaceReleaseSettings {
  enabled: boolean;
}

export interface MavenRepositorySettings {
  enabled: boolean;
}

export interface StorageSettings {
  enabled: boolean;
}

export interface AnalyticsSettings {
  enabled: boolean;
  require_write_key: boolean;
}

export interface SecretsSettings {
  enabled: boolean;
}

export interface SpaceLink {
  name: string;
  url: string;
}

export interface SpaceMember {
  user_id: string;
  role: string;
}

export function mapCategoryToDisplayname(name?: string): string {
  if (!name) return 'Unknown';

  switch (name.toLowerCase()) {
    case 'minecraft_plugin':
      return 'Minecraft Plugin';
    case 'minecraft_server':
      return 'Minecraft Server';
    case 'minecraft_mod':
      return 'Minecraft Mod';
    case 'hytale_plugin':
      return 'Hytale Plugin';
    case 'web_app':
      return 'Web App';
    case 'mobile_app':
      return 'Mobile App';
    case 'other':
      return 'Other';

      default:
      return name;
  }
}

export function mapLinkToDisplayname(name: string): string {
  switch (name.toLowerCase()) {
    case 'source_code':
      return 'Source Code';
    case 'documentation':
      return 'Documentation';
    case 'wiki':
      return 'Wiki';
    case 'discord':
      return 'Discord';
    case 'website':
      return 'Website';
    case 'issues':
      return 'Issues';
    default:
      return name;
  }
}

export function mapLinkToIcon(name: string): string {
  switch (name.toLowerCase()) {
    case 'source_code':
      return 'mdi-github';
    case 'documentation':
      return 'mdi-book-open-page-variant';
    case 'wiki':
      return 'mdi-book-open-page-variant';
    case 'discord':
      return 'mdi-chat';
    case 'website':
      return 'mdi-web';
    case 'issues':
      return 'mdi-bug-outline';
    default:
      return 'mdi-link';
  }
}
