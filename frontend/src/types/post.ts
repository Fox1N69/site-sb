/**
 * Post model interface matching backend structure
 * Backend model: backend/internal/models/post_model.go
 */
export interface Post {
    id: number;
    preview_image: string;      // Cover/thumbnail image
    images?: string[];           // Array of additional images
    title: string;
    description: string;
    text: string;
    created_at: string;
    updated_at: string;
    deleted_at?: string | null;
}

/**
 * API response types
 */
export interface PostResponse {
    message?: string;
    data?: Post;
}

export interface PostsListResponse extends Array<Post> { }

/**
 * Image upload response
 */
export interface UploadResponse {
    message: string;
    urls: string[];
    errors?: string[];
    count: number;
}
