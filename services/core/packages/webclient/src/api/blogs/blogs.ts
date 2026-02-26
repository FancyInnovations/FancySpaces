import type {BlogArticle} from "@/api/blogs/types.ts";
import {useUserStore} from "@/stores/user.ts";

export async function getBlogArticlesForSpace(spaceId: string): Promise<BlogArticle[]> {
  const userStore = useUserStore();

  const response = await fetch(
    `/api/v1/spaces/${spaceId}/blog-articles`,
    {
      method: "GET",
      headers: {
        "Accept": "application/json",
        "Authorization": `Bearer ${userStore.token}`,
      }
    },
  );

  if (!response.ok) {
    throw new Error("Failed to fetch blog articles: " + await response.text());
  }

  const articles = await response.json();
  if (!Array.isArray(articles)) {
    return [];
  }

  articles.forEach((article: BlogArticle) => {
    article.published_at = new Date(article.published_at);
  });

  return articles as BlogArticle[];
}

export async function getBlogArticlesForUser(userId: string): Promise<BlogArticle[]> {
  const userStore = useUserStore();

  const response = await fetch(
    `/api/v1/users/${userId}/blog-articles`,
    {
      method: "GET",
      headers: {
        "Accept": "application/json",
        "Authorization": `Bearer ${userStore.token}`,
      }
    },
  );

  if (!response.ok) {
    throw new Error("Failed to fetch blog articles: " + await response.text());
  }

  const articles = await response.json();
  articles.forEach((article: BlogArticle) => {
    article.published_at = new Date(article.published_at);
  });

  return articles as BlogArticle[];
}

export async function getBlogArticle(articleId: string): Promise<BlogArticle> {
  const userStore = useUserStore();

  const response = await fetch(
    `/api/v1/blog-articles/${articleId}`,
    {
      method: "GET",
      headers: {
        "Accept": "application/json",
        "Authorization": `Bearer ${userStore.token}`,
      }
    },
  );

  if (!response.ok) {
    throw new Error("Failed to fetch blog article: " + await response.text());
  }

  const article = await response.json();
  article.published_at = new Date(article.published_at);

  return article as BlogArticle;
}

export async function createBlogArticle(spaceId: string, title: string, summary: string, content: string): Promise<BlogArticle> {
  const userStore = useUserStore();

  const response = await fetch(
    `/api/v1/blog-articles`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Accept": "application/json",
        "Authorization": `Bearer ${userStore.token}`,
      },
      body: JSON.stringify({
        space_id: spaceId,
        title: title,
        summary: summary,
        content: content,
      }),
    },
  );

  if (!response.ok) {
    throw new Error("Failed to create blog article: " + await response.text());
  }

  const article = await response.json();
  article.published_at = new Date(article.published_at);

  return article as BlogArticle;
}

export async function getBlogArticleContent(articleId: string): Promise<string> {
  const userStore = useUserStore();

  const response = await fetch(
    `/api/v1/blog-articles/${articleId}/content`,
    {
      method: "GET",
      headers: {
        "Accept": "text/plain",
        "Authorization": `Bearer ${userStore.token}`,
      }
    },
  );

  if (!response.ok) {
    throw new Error("Failed to fetch blog article content: " + await response.text());
  }

  return await response.text();
}

export async function updateBlogArticle(articleId: string, title: string, summary: string, content: string): Promise<BlogArticle> {
  const userStore = useUserStore();

  const response = await fetch(
    `/api/v1/blog-articles/${articleId}`,
    {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        "Accept": "application/json",
        "Authorization": `Bearer ${userStore.token}`,
      },
      body: JSON.stringify({
        title,
        summary,
        content,
      }),
    },
  );

  if (!response.ok) {
    throw new Error("Failed to update blog article: " + await response.text());
  }

  const article = await response.json();
  article.published_at = new Date(article.published_at);

  return article as BlogArticle;
}

export async function deleteBlogArticle(articleId: string): Promise<void> {
  const userStore = useUserStore();

  const response = await fetch(
    `/api/v1/blog-articles/${articleId}`,
    {
      method: "DELETE",
      headers: {
        "Authorization": `Bearer ${userStore.token}`,
      }
    },
  );

  if (!response.ok) {
    throw new Error("Failed to delete blog article: " + await response.text());
  }
}
