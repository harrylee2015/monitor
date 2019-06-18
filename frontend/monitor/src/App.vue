<template>
  <div id="app">
    <router-view v-if="isRouterAlive" />
  </div>
</template>

<script>
export default {
  name: "App",
  provide() {
    return {
      reload: this.reload
    };
  },
  data() {
    return {
      isRouterAlive: true
    };
  },
  methods: {
    reload() {
      this.isRouterAlive = false;
      this.$nextTick(function() {
        this.isRouterAlive = true;
      });
    }
  },
  // mounted(){
    
  //   function checkIE(){
  //      return '-ms-scroll-limit' in document.documentElement.style && '-ms-ime-align' in document.documentElement.style;
  //   }
  //   if (checkIE()) {      
  //     window.addEventListener('hashchange', () => {        
  //       var currentPath = window.location.hash.slice(1);
  //       if (this.$route.path !== currentPath) {        
  //         this.$router.push(currentPath)      
  //       }    
  //     }, false)    
  //   }

  // }
};
</script>

<style>
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  min-width: 1150px;
  background: #020F34;
  display: flex;
  justify-content: center;
  /* overflow: auto */
}
/* @media screen and (max-width: 1500px) {
  .manager_content {
    width: 1500px;
  }
} */
</style>
