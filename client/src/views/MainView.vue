<template>
  <div class="main">
    <section class="home">
      <Navigation/>
      <div class="home__container">
        <div class="title">
           <img src="../assets/img/H1.png" alt="">
        </div>
        
        <div class="content">
            <div class="content__text">
              <div class="content__name">Ситниковский отдел Военно Исторического Клуба "Сибирский рубеж"</div>
              <div class="content__description">Военно Исторический Клуб "Сибирский рубеж" занимается углубленным изучением истории Великой Отечественной Войны</div>
            </div>
            <div class="content__solider">
              <img class="solider" src="../assets/img/solider.png" alt="">
            </div>
            <div class="content__video">
              <img class="play__icon" src="../assets/img/play.svg" alt="" style="width: 30px; height: 35px; margin: 148px 0 17px 25px;">
              <div class="video__text">Смотреть видео</div> 
            </div>
        </div>
      </div>
    </section>
    <AboutView/>
  </div>
</template>




<script lang="ts">
import Footer from '@/components/Footer/Footer.vue';
import Navigation from '@/components/Navigate/Navigate.vue'
import AboutView from './AboutView.vue';


export default{
    components: { Footer, Navigation, AboutView }
}

</script>


<style lang="scss" src="@/assets/sass/_mainview.scss">
</style>