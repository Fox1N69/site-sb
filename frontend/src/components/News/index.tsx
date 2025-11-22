"use client";
import { useEffect, useState } from "react";
import axios from "axios";
import type { Post } from "@/types/post";

const useFetchNews = () => {
    const [news, setNews] = useState<Post[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchData = async () => {
            try {
                setLoading(true);
                const apiUrl = import.meta.env.API_URL || 'http://localhost:4000';
                const res = await axios.get<Post[]>(`${apiUrl}/post/`);
                setNews(res.data);
                setError(null);
            } catch (error) {
                setError("Не удалось загрузить новости");
                console.error("Error fetching news:", error);
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, []);

    return { news, loading, error };
};

export const NewsTSX = () => {
    const { news, loading, error } = useFetchNews();

    if (loading) {
        return (
            <div className="h-full w-full flex items-center justify-center">
                <p className="text-lg">Загрузка новостей...</p>
            </div>
        );
    }

    if (error) {
        return (
            <div className="h-full w-full flex items-center justify-center">
                <p className="text-lg text-red-500">{error}</p>
            </div>
        );
    }

    return (
        <div className="h-full w-full">
            {news.length === 0 ? (
                <p className="text-lg text-center">Новостей пока нет</p>
            ) : (
                news.map((item) => (
                    <div key={item.id} className="mb-6 p-4 border rounded-lg">
                        {item.preview_image && (
                            <img
                                src={item.preview_image}
                                alt={item.title}
                                className="w-full h-48 object-cover rounded mb-4"
                            />
                        )}
                        <h2 className="text-2xl font-bold mb-2">{item.title}</h2>
                        <div className="font-medium text-lg mb-2">{item.description}</div>
                        <p className="text-gray-700">{item.text}</p>
                        <a
                            href={`/news/${item.id}`}
                            className="inline-block mt-4 text-blue-600 hover:text-blue-800"
                        >
                            Читать полностью →
                        </a>
                    </div>
                ))
            )}
        </div>
    );
};
