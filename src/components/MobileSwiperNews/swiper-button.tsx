import { useSwiper } from "swiper/react";

export default function SlideButtons() {
  const swiper = useSwiper();

  return (
    <div className="swiper-btns">
      <button className="swiper-btn-prev" onClick={() => swiper.slidePrev()}>
        {" "}
        <img src="/swiper-arrow.png" alt="arrow-prev" />{" "}
      </button>
      <button className="swiper-btn-next" onClick={() => swiper.slideNext()}>
        {" "}
        <img id="next-arrow-img" src="/swiper-arrow.png" alt="arrow-next" />{" "}
      </button>
    </div>
  );
}
