<div>
    <ul class="shops">
        {#each $shopStore as shop, index (shop.name)}
          {#if insertDropZoneBefore  === index}
          <li 
            class="shop-item" 
            style="background-color: gray;"
            on:drop|preventDefault={event => drop(event, index, shop.id)}
            >&nbsp;</li>
          {/if}
          <li 
            class="shop-item columns"
            style="background-color: {shop.color};"
            draggable={true}
            on:dragstart={event => dragstart(event, index, shop.id)}
            on:dragenter={(event) => dragenter(event, index)}
            on:dragleave={(event) => dragleave()}
            on:drop|preventDefault={event => drop(event, index, shop.id)}
            ondragover="return false"
            class:is-active={hovering === index}>

            <!-- input type="checkbox" bind:checked={item.done} / -->
            <div class="column is-10" on:click="{toggleFilterItems(shop.id)}">{shop.name}</div>
            <button class="column is-2" on:click={() => deleteShop(shop.id)}><Icon data={trash} class="no-pad"/></button>
          </li>
          {#if insertDropZoneAfter === index}
          <li 
            class="shop-item" 
            style="background-color: gray;"
            on:drop|preventDefault={event => drop(event, index, shop.id)}
            >&nbsp;</li>
          {/if}
        {:else}
        <li class="shop-item">No shops yet</li>
        {/each}
        <li 
          class="shop-item"
          on:click|preventDefault={event => toggleShowAddShopForm()}>
          <div class=""><Icon data={plus}/></div>
        </li>
            
    </ul>
    {#if showAddShopForm}
    <AddShop />
    {/if}

</div>

<script>

import { onMount } from 'svelte';
import {shopStore, loadAllShops, assignItemToShop, deleteShopFromStore} from './shop_store.js'
import {itemStore, loadAllItems} from './item_store'
import AddShop from './AddShop.svelte';
import Icon from "svelte-awesome";
import plus from 'svelte-awesome/icons/plus'
import trash from 'svelte-awesome/icons/trash'

let hovering = false;
let showAddShopForm = false;
let filterItemsByShop = false;
let draggingIndex = false;
let insertDropZoneBefore = false;
let insertDropZoneAfter = false;

onMount( () => {    
    loadAllShops()
  })

const toggleShowAddShopForm = () => {
  if (showAddShopForm) {
    showAddShopForm = false;
  } else {
    showAddShopForm = true;
  }
}

const toggleFilterItems = (shopId) => {
  if (filterItemsByShop && filterItemsByShop === shopId) {
    loadAllItems()
  }
 else {
   loadAllItems(shopId)
   filterItemsByShop = shopId
 }
}

// drag & drop support: dragstart
const dragstart = (event, index, shopid) => {
    event.dataTransfer.effectAllowed = 'move';
    //event.dataTransfer.dropEffect = 'move';
    event.dataTransfer.setData('index', index);
    event.dataTransfer.setData('shopid', shopid);
    draggingIndex = index
}

const dragenter = (event, index) => {
  // TODO: create dropzone
  // what is being dragged here?
  const dragitem = get_dragitem(event.dataTransfer.types)
  if (dragitem === 'item') {
    hovering = index
    return
  }
  else if (dragitem === 'shop') {
    if (draggingIndex < index) {
      insertDropZoneBefore = index
    } 
    else if (draggingIndex > index) {
      insertDropZoneAfter = index
    } 
   console.log("dragenter: " + index + ", start: "+draggingIndex)
  }
}
const dragleave = () => {
  console.log("dragleave")
  //insertDropZoneAfter = false;
  //insertDropZoneBefore = false;
}
const drop = (event, index, shopid) => {
  let dragitem = get_dragitem(event.dataTransfer.types)
  if (dragitem === 'item') {
    // a shopping cart item is assigned to shop
    assignItemToShop(event.dataTransfer.getData("itemid"), shopid)
    hovering = null
  } else if (dragitem === 'shop') {
    // reorder of shops:
    let start = draggingIndex;
    draggingIndex = false;
    event.dataTransfer.dropEffect = 'move'; 
    // TODO:
  }
}

const deleteShop = (shopId) => {

  $itemStore
    .filter(x => x.shop && x.shop.id === shopId)
    .forEach(x => assignItemToShop(x.id, 0))
  
  filterItemsByShop = false;

  deleteShopFromStore(shopId)
}

// find out, what is being dragged:
const get_dragitem = (array) => {
  if (array.includes('itemid')) {
    return 'item'
  } else if (array.includes('shopid')) {
    return 'shop'
  }
}
</script>

<style>
ul.shops {
  position: fixed;
  top: 122px;
  left: 5px;
  border: 1px solid #ccc;
  background-color: #f1f1f1;
  width: 20%;
}
li.shop-item {
  border-bottom: 1px solid #ccc;;

}
li.shop-item.is-active {
  background-color: #3273dc;
  color: #fff;
}
</style>