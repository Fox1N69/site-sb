"use client";
import { useEffect, useState } from "react";
import axios from "axios";

interface News {
  id: number;
  title: string;
  description: string;
  text: string;
}

const useFetchNews = () => {
  const [news, setNews] = useState<News[]>([]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const res = await axios.get<News[]>("http://localhost:4000/admin/news");
        setNews(res.data);
      } catch (error) {}
    };

    fetchData();
  }, []);
  return news;
};

export const NewsTSX = () => {
  const news = useFetchNews();

  return (
    <div className=" h-full w-full">
      {news.map((item) => (
        <>
          <h2>{item.title}</h2>
          <div className="font-medium text-lg">{item.description}</div>
          <p>{item.text}</p>
        </>
      ))}
    </div>
  );
};
