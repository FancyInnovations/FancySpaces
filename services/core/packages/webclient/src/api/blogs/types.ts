export interface BlogArticle {
  id: string;
  space_id?: string;
  author: string;
  title: string;
  summary: string;
  published_at: Date;
}
