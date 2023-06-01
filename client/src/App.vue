import router from './router/index';

<template>
  <nav>
    <router-link to="about"></router-link>
  </nav>
  <router-view/>
</template>

<style lang="scss">
@font-face {
  font-family: 'Wadik';
  src: url('./assets/fonts/Wadik.ttf');
  
}

@font-face {
  font-family: 'Roboto';
  src: url('./assets/fonts/Roboto.ttf'); 
}

*{
  padding: 0;
  margin: 0;
  box-sizing: border-box;
}
</style>
