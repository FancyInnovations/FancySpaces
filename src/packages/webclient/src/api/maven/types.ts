export interface SpaceMavenRepository {
  space_id: string;
  name: string;
  public: boolean;
  created_at: Date;
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
