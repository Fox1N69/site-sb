<template>
  <div class="main">
    <section class="home">
      <Navigation/>
      <div class="title">
        <h1 class="title__h1">Сибирский  рубежъ</h1>
      </div>
      
      <div class="content">
          <div class="content__name">Ситниковский отдел Военно Исторического Клуба "Сибирский рубеж"</div>
          <div class="content__description">Военно Исторический Клуб "Сибирский рубеж" занимается углубленным изучением истории Великой Отечественной Войны</div>
          <div class="content__solider">
            <img class="solider" src="../assets/img/solider.png" alt="">
          </div>
          <div class="content__video"></div>
      </div>
    </section>
  </div>
</template>




<script lang="ts">
import Footer from '@/components/Footer/Footer.vue';
import Navigation from '@/components/Navigate/Navigate.vue'


export default{
    components: { Footer, Navigation}
}

</script>


<style lang="scss" src="@/assets/sass/_mainview.scss">
</style>