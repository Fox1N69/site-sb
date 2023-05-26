<template>
  <nav>
    <router-link to="/">Home</router-link> |
    <router-link to="/about">About</router-link>
  </nav>
  <router-view/>
</template>

<style lang="scss">
@font-face {
  font-family: 'Wadik';
  src: url('./assets/fonts/Wadik.ttf');
}
</style>
