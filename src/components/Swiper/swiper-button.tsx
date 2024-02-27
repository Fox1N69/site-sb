import { useSwiper } from "swiper/react";

export default function SlideButtons() {
    const swiper = useSwiper();

    return(
        <div className="swiper-btns">
            <button onClick={() => swiper.slidePrev()}>Prev</button>
            <button onClick={() => swiper.slideNext()}>Next</button>
        </div>
    );
};
