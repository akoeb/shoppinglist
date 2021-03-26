<template>
  <div class="category">
    <h3>Category Component</h3>
    <ul>
      <li
        v-for="item in categoryList"
        :key="item.id"
        :data-id="item.id"
      >
      <div>{{ item.name }}</div>
      </li>
    </ul>
  </div>
</template>

<style scoped lang="scss">
.category {
  h3 {
    color: green;
  }
}
</style>

<script>
export default {
  name: 'Category',

  data() {
    return {
      categoryList
    }
  },
  created: function() {
    this.loadAllCategories()
  },
  methods: {
    loadAllCategories: function() {
      fetch('//localhost:8080/api/categories/', this.httpOptions())
        .then((res)  => res.json())
        .then((json) => {
            this.categoryList = json.items ? json.items : []
        })
        .catch((err) => console.error(err))
    },
  }
}
</script>
