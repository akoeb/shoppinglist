<template>
  <div class="item-list">
    <h3>ItemList Component</h3>
    <draggable class="list-group" id="listgroup" v-model="shoppingItems" @end="syncOrder">
      <li
        class="list-group-item"
        v-for="item in shoppingItems"
        :id="getItemId(item.id)"
        :key="item.id"
        v-on:click="toggleStatus(item.id)"
      >
        <span v-bind:class="{'status-checked': isClosed(item) }">{{ item.title }}</span>
        <span class="pull-right">
          <button class="btn btn-xs btn-danger" v-on:click="deleteItem(item.id)">
            <i class="fas fa-trash fa-1x" aria-hidden="true"></i>
          </button>
        </span>
      </li>
    </draggable>
  </div>
</template>

<style scoped lang="scss">
.item-list {
  h3 {
    color: black;
  }
}
</style>

<script>
export default {
  name: 'ItemList',
  data() {
    return {
      shoppingItems,
      activeCategory
    }
  },
  created: function() {
    this.loadAll()
  },
  methods: {
    loadAll: function() {
      fetch('/api/categories/' + this.activeCategory + '/items', this.httpOptions())
        .then((res)  => res.json())
        .then((json) => {
            this.shoppingItems = json.items ? json.items : []
        })
        .catch((err) => console.error(err))
    },
  }
}
</script>
