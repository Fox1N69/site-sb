import { Swiper, SwiperSlide } from "swiper/react";
import { Pagination } from "swiper/modules";
import "swiper/css";
import "swiper/css/pagination";
import SlideButtons from "./swiper-button";
import anime from "animejs";
import "@/assets/styles/swiper.scss";

// Определение типа данных для карточки новостей
interface NewsCard {
    id?: number;
    preview_image: string;
    title: string;
    description: string;
}

// Пример карточек новостей
const newsData: NewsCard[] = [
    {
        id: 1,
        preview_image: "",
        title: "Тестовая новость",
        description: "сайт находится в разработке."
    },
    {
        id: 4,
        preview_image: "",
        title: "Четвертая карточка",
        description: "Lorem ipsum dolor sit amet consectetur."
    },
];

export default function MySwiper() {
    const props = {
        loop: true,
        effect: "fade",
        speed: 1000,
        pagination: { clickable: true },
        navigation: {
            nextEl: ".button-next-slide",
            prevEl: ".button-prev-slide",
        },
        spaceBetween: 3,
        breakpoints: {
            320: {
                slidesPerView: 1,
                spaceBetween: 0,
            },
            1350: {
                slidesPerView: 2,
            },
            1450: {
                slidesPerView: 3,
            },
        },
    };

    return (
        <Swiper
            modules={[Pagination]}
            className="swiper"
            {...props}
            onSlideChangeTransitionStart={() => {
                anime({
                    targets: ".swiper-slide-active .team__info_rank",
                    translateY: [100, 0],
                    opacity: [0, 1],
                    delay: anime.stagger(250, { start: 500 }),
                    easing: "easeInOutQuart",
                });
                anime({
                    targets: ".swiper-slide-active .team__info_discription",
                    translateY: [100, 0],
                    opacity: [0, 1],
                    delay: anime.stagger(250, { start: 1000 }),
                    easing: "easeInOutQuart",
                });
            }}
        >
            {newsData.map((item) => (
                <SwiperSlide key={item.id} className="swiper__slide-news">
                    <div className="new__card">
                        <div className="news__card-img">
                            <img src={item.preview_image} alt="" className="news__img" />
                            <div className="overlay"></div>
                            <a href={`/news/${item.id}`} className="a__write new__btn">
                                Читать
                            </a>
                        </div>
                        <div className="new-card__content">
                            <div className="new-card__title">{item.title}</div>
                            <div className="new-card__description">{item.description}</div>
                        </div>
                    </div>
                </SwiperSlide>
            ))}
            <SlideButtons />
        </Swiper>
    );
}
