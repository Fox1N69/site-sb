import { useState, useEffect } from 'react';
import type { Post } from '@/types/post';
import ImageUpload from './ImageUpload';

export default function AdminPanel() {
    const [posts, setPosts] = useState<Post[]>([]);
    const [loading, setLoading] = useState(true);
    const [showCreateForm, setShowCreateForm] = useState(false);
    const [editingPost, setEditingPost] = useState<Post | null>(null);
    const [user, setUser] = useState<any>(null);

    // Form state
    const [formData, setFormData] = useState({
        title: '',
        description: '',
        text: '',
        preview_image: '',
        images: [] as string[]
    });

    useEffect(() => {
        checkAuth();
        fetchPosts();
    }, []);

    const checkAuth = () => {
        const token = localStorage.getItem('token');
        const userStr = localStorage.getItem('user');

        if (!token) {
            window.location.href = '/admin/login';
            return;
        }

        if (userStr) {
            setUser(JSON.parse(userStr));
        }
    };

    const fetchPosts = async () => {
        try {
            const apiUrl = import.meta.env.API_URL || 'http://localhost:4000';
            const response = await fetch(`${apiUrl}/post/`);
            const data = await response.json();
            setPosts(Array.isArray(data) ? data : []);
        } catch (error) {
            console.error('Error fetching posts:', error);
        } finally {
            setLoading(false);
        }
    };

    const handleLogout = () => {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        localStorage.removeItem('tokenExpiry');
        window.location.href = '/admin/login';
    };

    const handleImagesUploaded = (urls: string[]) => {
        setFormData(prev => ({
            ...prev,
            preview_image: urls[0] || '',
            images: urls.slice(1)
        }));
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        try {
            const token = localStorage.getItem('token');
            const apiUrl = import.meta.env.API_URL || 'http://localhost:4000';

            const endpoint = editingPost
                ? `${apiUrl}/admin/post/update/${editingPost.id}`
                : `${apiUrl}/admin/post/create`;

            const method = editingPost ? 'PATCH' : 'POST';

            const response = await fetch(endpoint, {
                method,
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify(formData)
            });

            if (!response.ok) {
                throw new Error('Ошибка при сохранении новости');
            }

            // Reset form and refresh posts
            setFormData({
                title: '',
                description: '',
                text: '',
                preview_image: '',
                images: []
            });
            setShowCreateForm(false);
            setEditingPost(null);
            await fetchPosts();

            alert(editingPost ? 'Новость обновлена!' : 'Новость создана!');
        } catch (error) {
            alert(error instanceof Error ? error.message : 'Ошибка при сохранении');
        }
    };

    const handleEdit = (post: Post) => {
        setEditingPost(post);
        setFormData({
            title: post.title,
            description: post.description,
            text: post.text,
            preview_image: post.preview_image,
            images: post.images || []
        });
        setShowCreateForm(true);
        window.scrollTo({ top: 0, behavior: 'smooth' });
    };

    const handleDelete = async (id: number) => {
        if (!confirm('Вы уверены, что хотите удалить эту новость?')) {
            return;
        }

        try {
            const token = localStorage.getItem('token');
            const apiUrl = import.meta.env.API_URL || 'http://localhost:4000';

            const response = await fetch(`${apiUrl}/admin/post/delete/${id}`, {
                method: 'DELETE',
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            if (!response.ok) {
                throw new Error('Ошибка при удалении новости');
            }

            await fetchPosts();
            alert('Новость удалена!');
        } catch (error) {
            alert(error instanceof Error ? error.message : 'Ошибка при удалении');
        }
    };

    const cancelEdit = () => {
        setEditingPost(null);
        setFormData({
            title: '',
            description: '',
            text: '',
            preview_image: '',
            images: []
        });
        setShowCreateForm(false);
    };

    if (loading) {
        return (
            <div className="admin-loading">
                <svg className="spinner" viewBox="0 0 50 50">
                    <circle className="path" cx="25" cy="25" r="20" fill="none" strokeWidth="5"></circle>
                </svg>
                <p>Загрузка...</p>
            </div>
        );
    }

    return (
        <div className="admin-panel">
            {/* Header */}
            <header className="admin-header">
                <div className="header-content">
                    <div className="header-left">
                        <h1>Админ-панель</h1>
                        <p>Управление новостями</p>
                    </div>
                    <div className="header-right">
                        {user && (
                            <div className="user-info">
                                <span className="user-name">{user.fio || user.username}</span>
                                <span className="user-role">{user.role}</span>
                            </div>
                        )}
                        <button onClick={handleLogout} className="logout-btn">
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                            </svg>
                            Выйти
                        </button>
                    </div>
                </div>
            </header>

            <div className="admin-container">
                {/* Create/Edit Form */}
                {showCreateForm ? (
                    <div className="form-section">
                        <div className="section-header">
                            <h2>{editingPost ? 'Редактировать новость' : 'Создать новость'}</h2>
                            <button onClick={cancelEdit} className="cancel-btn">
                                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                                </svg>
                                Отмена
                            </button>
                        </div>

                        <form onSubmit={handleSubmit} className="post-form">
                            <div className="form-group">
                                <label htmlFor="title">Заголовок *</label>
                                <input
                                    type="text"
                                    id="title"
                                    value={formData.title}
                                    onChange={(e) => setFormData(prev => ({ ...prev, title: e.target.value }))}
                                    required
                                    placeholder="Введите заголовок новости"
                                />
                            </div>

                            <div className="form-group">
                                <label htmlFor="description">Описание *</label>
                                <textarea
                                    id="description"
                                    value={formData.description}
                                    onChange={(e) => setFormData(prev => ({ ...prev, description: e.target.value }))}
                                    required
                                    rows={3}
                                    placeholder="Краткое описание новости"
                                />
                            </div>

                            <div className="form-group">
                                <label htmlFor="text">Полный текст *</label>
                                <textarea
                                    id="text"
                                    value={formData.text}
                                    onChange={(e) => setFormData(prev => ({ ...prev, text: e.target.value }))}
                                    required
                                    rows={10}
                                    placeholder="Полный текст новости"
                                />
                            </div>

                            <div className="form-group">
                                <ImageUpload
                                    onImagesUploaded={handleImagesUploaded}
                                    existingImages={editingPost ? [formData.preview_image, ...formData.images] : []}
                                    maxImages={10}
                                />
                            </div>

                            <div className="form-actions">
                                <button type="submit" className="submit-btn">
                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                                    </svg>
                                    {editingPost ? 'Обновить новость' : 'Создать новость'}
                                </button>
                            </div>
                        </form>
                    </div>
                ) : (
                    <>
                        {/* Stats */}
                        <div className="stats-grid">
                            <div className="stat-card">
                                <div className="stat-icon">
                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z" />
                                    </svg>
                                </div>
                                <div className="stat-content">
                                    <p className="stat-label">Всего новостей</p>
                                    <p className="stat-value">{posts.length}</p>
                                </div>
                            </div>

                            <button onClick={() => setShowCreateForm(true)} className="create-news-card">
                                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                                </svg>
                                <span>Создать новость</span>
                            </button>
                        </div>

                        {/* Posts List */}
                        <div className="posts-section">
                            <h2>Все новости</h2>

                            {posts.length === 0 ? (
                                <div className="empty-state">
                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9.a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z" />
                                    </svg>
                                    <p>Новостей пока нет</p>
                                    <button onClick={() => setShowCreateForm(true)} className="empty-action-btn">
                                        Создать первую новость
                                    </button>
                                </div>
                            ) : (
                                <div className="posts-grid">
                                    {posts.map(post => (
                                        <div key={post.id} className="post-card">
                                            {post.preview_image && (
                                                <div className="post-image">
                                                    <img src={post.preview_image} alt={post.title} />
                                                </div>
                                            )}
                                            <div className="post-content">
                                                <h3>{post.title}</h3>
                                                <p className="post-description">{post.description}</p>
                                                <div className="post-meta">
                                                    <span className="post-date">
                                                        {new Date(post.created_at).toLocaleDateString('ru-RU')}
                                                    </span>
                                                </div>
                                            </div>
                                            <div className="post-actions">
                                                <button onClick={() => handleEdit(post)} className="edit-btn">
                                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                                                    </svg>
                                                    Редактировать
                                                </button>
                                                <button onClick={() => handleDelete(post.id)} className="delete-btn">
                                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                                                    </svg>
                                                    Удалить
                                                </button>
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            )}
                        </div>
                    </>
                )}
            </div>
        </div>
    );
}
