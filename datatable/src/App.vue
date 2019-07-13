<template>
  <div id="app">
    <vue-data-table></vue-data-table>
  </div>
</template>

<script>
import VueDataTable from './components/VueDataTable'

import '../node_modules/bootstrap/dist/css/bootstrap.min.css'

export default {
  name: 'app',
  components: {
    VueDataTable
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
