<script>
  import { itemStore } from "./item_store.js";
  import { uid } from "../util.js";

  // one temporary item to be used for creation
  let newItem = {};

  // create an item
  const addItem = () => {
    // if title is empty we do not do anything
    if (!newItem.title.trim()) {
      newItem = {};
      return;
    }
    // add item status and ordering
    newItem.status = "OPEN";
    // orderno must be the last of the shopping list plus 1
    newItem.orderno = $itemStore.length + 1;
    newItem.uid = uid();

    let store = $itemStore;
    store.items.push(newItem);
    store.local = true;
    $itemStore = store;
    newItem = {};
  };
</script>

<div>
  <form on:submit|preventDefault={addItem}>
    <div class="columns mt-5 mr-2 ml-2">
      <input
        id="newItemTitle"
        type="text"
        class="input column is-10 mr-1"
        bind:value={newItem.title}
        placeholder="New Item"
      />
      <button
        class="button is-link column is-2"
        on:click|preventDefault={addItem}>Submit</button
      >
    </div>
  </form>
</div>

<style></style>
