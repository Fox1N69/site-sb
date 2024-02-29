import { Swiper, SwiperSlide } from "swiper/react";
import "swiper/css";
import SlideButtons from "./swiper-button";
import SlideButtonsNext from "./swiper-button-next";
import anime from "animejs";
import "@/assets/styles/swiper.scss";

export default function MySwiper() {
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
        <div className="fullname">Максимович Вячеслав</div>
        <div className="slide__left">
          <img src="/Vacheclav.png" alt="" />
        </div>
        <div className="slide__right">
          <div className="team__info">
            <div className="team__info_rank">Должность</div>
            <div className="team__info_discription">
              Lorem ipsum dolor sit amet consectetur. Morbi sit ut id feugiat
              nisl amet aliquet nibh at. Quam fringilla orci pellentesque orci.
              Arcu in nisi neque nunc. A auctor elit orci sem risus pellentesque
              facilisi ullamcorper pharetra. Eu ante amet massa vitae quis.
              Scelerisque faucibus ipsum nibh augue.
            </div>
          </div>
        </div>
      </SwiperSlide>

      {/*Slide 2*/}

      <SwiperSlide className="swiper__slide">
        <div className="fullname">Пинигин Станислав</div>
        <div className="slide__left">
          <img src="" alt="" />
        </div>
        <div className="slide__right">
          <div className="team__info">
            <div className="team__info_rank">Должность</div>
            <div className="team__info_discription">
              Lorem ipsum dolor sit amet consectetur. Morbi sit ut id feugiat
              nisl amet aliquet nibh at. Quam fringilla orci pellentesque orci.
              Arcu in nisi neque nunc. A auctor elit orci sem risus pellentesque
              facilisi ullamcorper pharetra. Eu ante amet massa vitae quis.
              Scelerisque faucibus ipsum nibh augue.
            </div>
          </div>
        </div>
      </SwiperSlide>

      {/*Slide 3*/}

      <SwiperSlide className="swiper__slide">
        <div className="fullname">Максимович Павел</div>
        <div className="slide__left">
          <img src="" alt="" />
        </div>
        <div className="slide__right">
          <div className="team__info">
            <div className="team__info_rank">Должность</div>
            <div className="team__info_discription">
              Lorem ipsum dolor sit amet consectetur. Morbi sit ut id feugiat
              nisl amet aliquet nibh at. Quam fringilla orci pellentesque orci.
              Arcu in nisi neque nunc. A auctor elit orci sem risus pellentesque
              facilisi ullamcorper pharetra. Eu ante amet massa vitae quis.
              Scelerisque faucibus ipsum nibh augue.
            </div>
          </div>
        </div>
      </SwiperSlide>

      {/*Slide 4*/}

      <SwiperSlide className="swiper__slide">
        <div className="slide__left">
          <img src="/video.png" alt="" />
        </div>
        <div className="slide__right">
          <div className="team__info">
            <div className="team__info_rank">Должность</div>
            <div className="team__info_discription">
              Lorem ipsum dolor sit amet consectetur. Morbi sit ut id feugiat
              nisl amet aliquet nibh at. Quam fringilla orci pellentesque orci.
              Arcu in nisi neque nunc. A auctor elit orci sem risus pellentesque
              facilisi ullamcorper pharetra. Eu ante amet massa vitae quis.
              Scelerisque faucibus ipsum nibh augue.
            </div>
          </div>
        </div>
      </SwiperSlide>
      <SlideButtons />
    </Swiper>
  );
}
