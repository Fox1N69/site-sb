import { useState, useRef, useCallback } from 'react';
import type { UploadResponse } from '@/types/post';

interface ImageUploadProps {
    onImagesUploaded: (urls: string[]) => void;
    maxImages?: number;
    existingImages?: string[];
}

interface ImagePreview {
    id: string;
    url: string;
    file?: File;
    uploaded: boolean;
}

export default function ImageUpload({
    onImagesUploaded,
    maxImages = 10,
    existingImages = []
}: ImageUploadProps) {
    const [images, setImages] = useState<ImagePreview[]>(
        existingImages.map((url, index) => ({
            id: `existing-${index}`,
            url,
            uploaded: true
        }))
    );
    const [isDragging, setIsDragging] = useState(false);
    const [uploading, setUploading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const fileInputRef = useRef<HTMLInputElement>(null);

    const handleDragOver = useCallback((e: React.DragEvent) => {
        e.preventDefault();
        setIsDragging(true);
    }, []);

    const handleDragLeave = useCallback((e: React.DragEvent) => {
        e.preventDefault();
        setIsDragging(false);
    }, []);

    const handleDrop = useCallback((e: React.DragEvent) => {
        e.preventDefault();
        setIsDragging(false);

        const files = Array.from(e.dataTransfer.files);
        handleFiles(files);
    }, [images]);

    const handleFileInput = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.files) {
            const files = Array.from(e.target.files);
            handleFiles(files);
        }
    }, [images]);

    const handleFiles = (files: File[]) => {
        setError(null);

        // Validate file count
        if (images.length + files.length > maxImages) {
            setError(`Максимум ${maxImages} изображений`);
            return;
        }

        // Validate file types and sizes
        const validFiles: File[] = [];
        const errors: string[] = [];

        files.forEach(file => {
            if (!file.type.match(/image\/(jpg|jpeg|png|gif|webp)/)) {
                errors.push(`${file.name}: неподдерживаемый формат`);
                return;
            }
            if (file.size > 10 * 1024 * 1024) {
                errors.push(`${file.name}: файл слишком большой (max 10MB)`);
                return;
            }
            validFiles.push(file);
        });

        if (errors.length > 0) {
            setError(errors.join(', '));
        }

        // Create previews
        const newPreviews: ImagePreview[] = validFiles.map((file, index) => ({
            id: `new-${Date.now()}-${index}`,
            url: URL.createObjectURL(file),
            file,
            uploaded: false
        }));

        setImages(prev => [...prev, ...newPreviews]);
    };

    const uploadImages = async () => {
        const filesToUpload = images.filter(img => !img.uploaded && img.file);

        if (filesToUpload.length === 0) {
            const uploadedUrls = images.filter(img => img.uploaded).map(img => img.url);
            onImagesUploaded(uploadedUrls);
            return;
        }

        setUploading(true);
        setError(null);

        try {
            const formData = new FormData();
            filesToUpload.forEach(img => {
                if (img.file) {
                    formData.append('images', img.file);
                }
            });

            const token = localStorage.getItem('token');
            const apiUrl = import.meta.env.API_URL || 'http://localhost:4000';

            const response = await fetch(`${apiUrl}/admin/upload/images`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`
                },
                body: formData
            });

            if (!response.ok) {
                throw new Error('Ошибка загрузки изображений');
            }

            const data: UploadResponse = await response.json();

            // Update images with uploaded URLs
            setImages(prev => {
                const uploaded = prev.filter(img => img.uploaded);
                const newUploaded = data.urls.map((url, index) => ({
                    id: `uploaded-${Date.now()}-${index}`,
                    url: `${apiUrl}${url}`,
                    uploaded: true
                }));
                return [...uploaded, ...newUploaded];
            });

            // Notify parent component
            const allUrls = [
                ...images.filter(img => img.uploaded).map(img => img.url),
                ...data.urls.map(url => `${apiUrl}${url}`)
            ];
            onImagesUploaded(allUrls);

            if (data.errors && data.errors.length > 0) {
                setError(data.errors.join(', '));
            }
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Ошибка загрузки');
        } finally {
            setUploading(false);
        }
    };

    const removeImage = (id: string) => {
        setImages(prev => prev.filter(img => img.id !== id));

        // Notify parent
        const remainingUrls = images
            .filter(img => img.id !== id && img.uploaded)
            .map(img => img.url);
        onImagesUploaded(remainingUrls);
    };

    const moveImage = (fromIndex: number, toIndex: number) => {
        setImages(prev => {
            const newImages = [...prev];
            const [removed] = newImages.splice(fromIndex, 1);
            newImages.splice(toIndex, 0, removed);

            // Notify parent with new order
            const uploadedUrls = newImages
                .filter(img => img.uploaded)
                .map(img => img.url);
            onImagesUploaded(uploadedUrls);

            return newImages;
        });
    };

    return (
        <div className="image-upload">
            <div className="upload-header">
                <h3>Изображения новости</h3>
                <p className="upload-hint">
                    Первое изображение будет использовано как обложка
                </p>
            </div>

            <div
                className={`drop-zone ${isDragging ? 'dragging' : ''}`}
                onDragOver={handleDragOver}
                onDragLeave={handleDragLeave}
                onDrop={handleDrop}
                onClick={() => fileInputRef.current?.click()}
            >
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
                </svg>
                <p className="drop-zone-text">
                    Перетащите изображения сюда или нажмите для выбора
                </p>
                <p className="drop-zone-subtext">
                    JPG, PNG, GIF, WEBP до 10MB (макс. {maxImages} изображений)
                </p>
                <input
                    ref={fileInputRef}
                    type="file"
                    accept="image/jpeg,image/jpg,image/png,image/gif,image/webp"
                    multiple
                    onChange={handleFileInput}
                    style={{ display: 'none' }}
                />
            </div>

            {error && (
                <div className="upload-error">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    {error}
                </div>
            )}

            {images.length > 0 && (
                <div className="images-preview">
                    {images.map((image, index) => (
                        <div key={image.id} className="image-preview-item">
                            {index === 0 && (
                                <div className="cover-badge">Обложка</div>
                            )}
                            <img src={image.url} alt={`Preview ${index}`} />
                            <div className="image-overlay">
                                <div className="image-actions">
                                    {index > 0 && (
                                        <button
                                            type="button"
                                            className="action-btn"
                                            onClick={() => moveImage(index, index - 1)}
                                            title="Переместить влево"
                                        >
                                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                                            </svg>
                                        </button>
                                    )}
                                    {index < images.length - 1 && (
                                        <button
                                            type="button"
                                            className="action-btn"
                                            onClick={() => moveImage(index, index + 1)}
                                            title="Переместить вправо"
                                        >
                                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                                            </svg>
                                        </button>
                                    )}
                                    <button
                                        type="button"
                                        className="action-btn delete-btn"
                                        onClick={() => removeImage(image.id)}
                                        title="Удалить"
                                    >
                                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                                        </svg>
                                    </button>
                                </div>
                            </div>
                            {!image.uploaded && (
                                <div className="upload-pending-badge">
                                    Не загружено
                                </div>
                            )}
                        </div>
                    ))}
                </div>
            )}

            {images.some(img => !img.uploaded) && (
                <button
                    type="button"
                    className="upload-button"
                    onClick={uploadImages}
                    disabled={uploading}
                >
                    {uploading ? (
                        <>
                            <svg className="spinner" viewBox="0 0 50 50">
                                <circle className="path" cx="25" cy="25" r="20" fill="none" strokeWidth="5"></circle>
                            </svg>
                            Загрузка...
                        </>
                    ) : (
                        <>
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
                            </svg>
                            Загрузить изображения ({images.filter(img => !img.uploaded).length})
                        </>
                    )}
                </button>
            )}
        </div>
    );
}
