<script>
  import { onMount, createEventDispatcher } from "svelte";
  import { shopStore } from "./shop_store.js";
  import { itemStore } from "./item_store.js";
  import { reorderStore, backend, httpOptions } from "../util.js";
  import AddShop from "./AddShop.svelte";
  import Icon from "svelte-awesome";
  import plus from "svelte-awesome/icons/plus";
  import trash from "svelte-awesome/icons/trash";

  let hovering = false;
  let filterItemsByShop = false;
  let showModal = false;

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

  // this function loads list from backend and replaces store if backend is newer
  export const getFromBackend = () => {
    const localVersion = localStorage.getItem("shopVersion");
    fetch(backend("api/shops"), httpOptions("GET"))
      .then((res) => res.json())
      .then((obj) => {
        if (obj.version >= localVersion) {
          obj.local = false;
          $shopStore = obj;
        }
      })
      .catch((err) => console.error(err));
  };

  // ********************************* //
  //     LOCAL  FUNCTIONS
  // ********************************* //

  // delete Shop from the list
  const deleteShop = (index) => {
    filterItemsByShop = false;

    // create copy of shop list
    let newStore = $shopStore;
    // item to delete
    let deletedShop = newStore.items[index];

    // remove item from  list
    newStore.items.splice(index, 1);

    // this is a local change
    newStore.local = true;

    // rewrite back to store
    $shopStore = newStore;

    // remove shop references in the items list
    let newItemStore = $itemStore;
    newItemStore.items.forEach((item) => {
      if (item.shop.uid === deletedShop.uid) {
        item.shop = {};
      }
    });
    newItemStore.local = true;
    $itemStore = newItemStore;
  };

  //assign an item to a shop
  const assignItemToShop = (itemid, shop) => {
    let newItemStore = $itemStore;
    let item = newItemStore.items.find((obj) => {
      return obj.uid === itemid;
    });
    if (!item) {
      console.error("item not found");
      return;
    }
    item.shop = shop;
    newItemStore.local = true;
    $itemStore = newItemStore;
  };

  // shop has been clicked, filter the item list
  // to only show items from this shop
  const toggleFilterItemsByShop = (shopId) => {
    // change local state:
    if (filterItemsByShop && filterItemsByShop === shopId) {
      filterItemsByShop = false;
    } else {
      filterItemsByShop = shopId;
    }
    // notify parent to filter item list:
    dispatch("filterItemsByShop", shopId);
  };

  // ********************************* //
  // DRAG  AND  DROP  FUNCTIONS
  // ********************************* //

  // drag & drop support: dragstart
  const dragstart = (event, index, shopid) => {
    event.dataTransfer.effectAllowed = "move";
    event.dataTransfer.setData("index", index);
    event.dataTransfer.setData("shopid", shopid);
  };

  const dragenter = (event, index) => {
    hovering = index;
  };

  const drop = (event, target, shop) => {
    const start = parseInt(event.dataTransfer.getData("index"));
    const item_type = get_item_type(event.dataTransfer.types);
    dispatch("resetHovering");

    if (item_type === "item") {
      // a shopping cart item is assigned to shop
      assignItemToShop(event.dataTransfer.getData("itemid"), shop);
    } else if (item_type === "shop") {
      // reorder of shops:
      event.dataTransfer.dropEffect = "move";
      $shopStore = reorderStore($shopStore, start, target);
    }
  };

  // find out, what is being dragged:
  const get_item_type = (array) => {
    if (array.includes("itemid")) {
      return "item";
    } else if (array.includes("shopid")) {
      return "shop";
    }
  };
</script>

<div class="shop-list is-pulled-left">
  <ul>
    {#each $shopStore.items as shop, index (index)}
      <li
        class="shop-item columns"
        style="background-color: {shop.color};"
        draggable={true}
        on:dragstart={(event) => dragstart(event, index, shop)}
        on:dragenter={(event) => dragenter(event, index)}
        on:drop|preventDefault={(event) => drop(event, index, shop)}
        ondragover="return false"
        class:is-active={hovering === index}
        class:is-filtered={filterItemsByShop === shop.uid}
      >
        <div class="column is-10" on:click={toggleFilterItemsByShop(shop.uid)}>
          {shop.name}
        </div>
        <button class="column is-2" on:click={() => deleteShop(index)}
          ><Icon data={trash} class="no-pad" /></button
        >
      </li>
    {:else}
      <li class="shop-item">No shops yet</li>
    {/each}
    <li class="shop-item" on:click|preventDefault={() => (showModal = true)}>
      <div class=""><Icon data={plus} /></div>
    </li>
  </ul>
  <AddShop bind:showModal />
</div>

<style>
  .shop-list {
    position: fixed;
    top: 122px;
    left: 5px;
    border: 1px solid #ccc;
    border-radius: 4px;
    background-color: #f1f1f1;
    width: 20%;
    box-shadow:
      0 2px 3px rgba(10, 10, 10, 0.1),
      0 0 0 1px rgba(10, 10, 10, 0.1);
    padding-bottom: 0.8em;
  }
  li.shop-item {
    border-bottom: 1px solid #ccc;
    padding: 1px;
    margin: 1px;
  }
  li.shop-item.is-active {
    background-color: #3273dc;
    color: #fff;
    border: 1px;
    border-color: black;
  }
  li.shop-item.is-filtered {
    border: 4px solid red;
    border-radius: 4px;
  }
</style>
