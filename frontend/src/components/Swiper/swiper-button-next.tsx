import { useSwiper } from "swiper/react";

export default function SlideButtonsNext() {
  const swiper = useSwiper();

  return (
    <button className="swiper-btn-next" onClick={() => swiper.slideNext()}>
      {" "}
      <img id="next-arrow-img" src="/swiper-arrow.png" alt="arrow-next" />{" "}
    </button>
  );
}
