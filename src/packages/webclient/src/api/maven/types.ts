export interface SpaceMavenRepository {
  space_id: string;
  name: string;
  public: boolean;
  created_at: Date;
  internal_mirror?: SpaceMavenRepositoryInternalMirror | null;
}

export interface SpaceMavenRepositoryInternalMirror {
  space_id: string;
  repository: string;
}

export interface SpaceMavenRepositoryArtifact {
  space_id: string;
  repository: string;
  group: string;
  id: string;
  versions: SpaceMavenRepositoryArtifactVersion[];
}

export interface SpaceMavenRepositoryArtifactVersion {
  version: string;
  published_at: Date;
  files: SpaceMavenRepositoryArtifactVersionFile[];
}

export interface SpaceMavenRepositoryArtifactVersionFile {
  name: string;
  size: number;
  url: string;
}
