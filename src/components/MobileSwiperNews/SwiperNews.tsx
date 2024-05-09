import { Swiper, SwiperSlide } from "swiper/react";
import { Pagination } from "swiper/modules";
import "swiper/css";
import SlideButtons from "./swiper-button";
import "swiper/css/pagination";
import anime from "animejs";
import "@/assets/styles/news_swiper.scss";

export default function NewsSwiper() {
  const props = {
    loop: true,
    effect: "fade",
    speed: 1000,
    navigation: {
      nextEl: ".button-next-slide",
      prevEL: ".button-prev.slide",
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
      pagination={true}
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
      {/*Slide 1*/}
      <SwiperSlide className="swiper__slide">
        <div className="news__card">
          <img src="/newcardnohover1.png" alt="" />
          <div className="card__content">
            <h3 className="card__title"></h3>
            <p className="card__text"></p>
          </div>
        </div>
      </SwiperSlide>

      {/*Slide 2*/}

      <SwiperSlide className="swiper__slide"></SwiperSlide>

      {/*Slide 3*/}

      <SwiperSlide className="swiper__slide"></SwiperSlide>
    </Swiper>
  );
}
