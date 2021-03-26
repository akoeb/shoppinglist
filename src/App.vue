<template>
<div class="home">
    <img alt="Vue logo" src="assets/logo.png">
    <Category />
    <Location />
    <ItemList />
    <AddItem />
</div>
</template>

<script>
// @ is an alias to /src
import Category from '@/components/Category.vue'
import Location from '@/components/Location.vue'
import ItemList from '@/components/ItemList.vue'
import AddItem from '@/components/AddItem.vue'

export default {
  name: 'home',
  components: {
    Category, Location, ItemList, AddItem
  }
}
</script>

<style lang="scss">
@import "app.scss";
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
#nav {
  padding: 30px;
  a {
    font-weight: bold;
    color: #2c3e50;
  }
}
</style>
