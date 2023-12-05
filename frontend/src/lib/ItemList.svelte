<script>
  // @ts-nocheck
  import { flip } from "svelte/animate";
  import { onMount, createEventDispatcher } from "svelte";
  import { itemStore } from "./item_store.js";
  import Icon from "svelte-awesome";
  import trash from "svelte-awesome/icons/trash";
  import AddItem from "./AddItem.svelte";
  import { reorderStore, backend, httpOptions } from "../util.js";

  // show item as active on dragover
  let hovering = false;

  // filter items by shop
  let filterShopId = false;

  const dispatch = createEventDispatcher();

  onMount(() => {});

  /*
   ** functions:
   */

  // export function to reset the hovering flag
  // to be called from parent component
  export const resetHovering = () => {
    hovering = false;
  };

  // reset shop ID for filtering, to be called from parent
  export const toggleFilterShopId = (shopId) => {
    if (filterShopId && filterShopId === shopId) {
      filterShopId = false;
    } else {
      filterShopId = shopId;
    }
  };

  // this function loads list from backend and replaces store if backend is newer
  export const getFromBackend = () => {
    const localVersion = localStorage.getItem("itemVersion");
    fetch(backend("api/items"), httpOptions("GET"))
      .then((res) => res.json())
      .then((obj) => {
        if (obj.version >= localVersion) {
          obj.local = false;
          $itemStore = obj;
        }
      })
      .catch((err) => console.error(err));
  };

  // ********************************* //
  //     LOCAL  FUNCTIONS
  // ********************************* //

  // function for filtering the items array
  $: itemsList = () => {
    let retval = [];
    // show only items from one particular shop
    if (filterShopId) {
      $itemStore.items.forEach((item) => {
        if (item.shop && item.shop.uid === filterShopId) {
          retval.push(item);
        }
      });
    } else {
      // show all items
      retval = $itemStore.items;
    }
    return retval;
  };

  // delete an item from the store
  const deleteItem = (index) => {
    let newStore = $itemStore;
    newStore.items.splice(index, 1);
    newStore.local = true;
    $itemStore = newStore;
  };

  // toggle the strike through flag for items already bought
  const toggleStatus = (index) => {
    let newStore = $itemStore;
    let status = newStore.items[index].status;
    if (status === "CLOSED") {
      newStore.items[index].status = "OPEN";
    } else {
      newStore.items[index].status = "CLOSED";
    }
    newStore.local = true;
    $itemStore = newStore;
  };

  // get element styling depending on item properties:
  const getStyle = (item) => {
    let style = "";
    if (item.shop) {
      style = "background-color: " + item.shop.color + ";";
    }
    // if we want to show the item, use shops bg color
    return style;
  };

  // ********************************* //
  // DRAG  AND  DROP  FUNCTIONS
  // ********************************* //

  // drag & drop support: drop
  const drop = (event, target) => {
    event.dataTransfer.dropEffect = "move";
    const start = parseInt(event.dataTransfer.getData("index"));
    dispatch("resetHovering");
    $itemStore = reorderStore($itemStore, start, target);
  };

  // drag & drop support: dragstart
  const dragstart = (event, index, itemid) => {
    event.dataTransfer.effectAllowed = "move";
    event.dataTransfer.setData("index", index);
    event.dataTransfer.setData("itemid", itemid);
  };
</script>

<div>
  <div class="container">
    <ul class="list is-pulled-right">
      {#each itemsList() as item, index (index)}
        <li
          class="list-item columns"
          animate:flip
          style={getStyle(item)}
          on:click={toggleStatus(index)}
          draggable={true}
          on:dragstart={(event) => dragstart(event, index, item.uid)}
          on:drop|preventDefault={(event) => drop(event, index)}
          ondragover="return false"
          on:dragenter={() => (hovering = index)}
          class:is-active={hovering === index}
        >
          <span class="column is-11" class:checked={item.status === "CLOSED"}
            >{item.title}</span
          >
          <button class="column is-1" on:click={() => deleteItem(index)}
            ><Icon data={trash} class="no-pad" /></button
          >
        </li>
      {:else}
        <p class="list-item columns">The list is empty</p>
      {/each}

      <AddItem />
    </ul>
  </div>
</div>

<style>
  .list {
    border-radius: 4px;
    box-shadow:
      0 2px 3px rgba(10, 10, 10, 0.1),
      0 0 0 1px rgba(10, 10, 10, 0.1);
    width: 80%;
    padding-bottom: 0.8em;
  }

  .list-item {
    padding: 0.5em;
    margin: 0 1em;
    text-align: left;
  }

  .list-item:not(:last-child) {
    border-bottom: 1px solid #dbdbdb;
  }

  .list-item.is-active {
    background-color: #3273dc;
    color: #fff;
  }
  span.checked {
    text-decoration: line-through;
  }
</style>
