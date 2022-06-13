<template>
  <router-view/>
</template>

<script>
import Menu from "@/components/Menu";
import Result from "@/components/Result";
export default {
  components: {Result, Menu}
}
</script>

<style lang="scss">
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

nav {
  padding: 30px;

  a {
    font-weight: bold;
    color: #2c3e50;

    &.router-link-exact-active {
      color: #42b983;
    }
  }
}

.router-link-active {
  text-decoration: none;
  color: #fff;
}
a {
  text-decoration: none;
  color: inherit;
}

</style>
