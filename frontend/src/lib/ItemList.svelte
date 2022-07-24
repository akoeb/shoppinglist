
<div>
    <div class="container">
      <ul class="list is-pulled-right">
        {#each $itemStore as item, index (item.title)}
           
          <li 
            class="list-item columns"
            animate:flip
            style="{getStyle(item)}"
            on:click="{toggleStatusInStore(item.id)}" 
            draggable={true}
            on:dragstart={event => dragstart(event, index, item.id)}
            on:drop|preventDefault={event => drop(event, index)}
            ondragover="return false"
            on:dragenter={() => hovering = index}
            class:is-active={hovering === index}>

            <span class="column is-11" class:checked="{item.status === 'CLOSED'}">{item.title}</span>
            <button class="column is-1" on:click={() => deleteItemFromStore(item.id)}><Icon data={trash} class="no-pad"/></button>
          </li>
        {:else}
          <p>The list is empty</p>
        {/each}
      </ul>
    </div>  
   
</div>


<script>
  // @ts-nocheck
  import {flip} from 'svelte/animate';
  import { onMount } from 'svelte';
  import {itemStore, loadAllItems, reorderItemsInStore, deleteItemFromStore, toggleStatusInStore} from './item_store.js'
  import Icon from "svelte-awesome";
  import trash from 'svelte-awesome/icons/trash'
  

  // some variable for drg&drop
  let hovering = false;

  onMount( () => {    
    loadAllItems()
  })

  /*
  ** functions:
  */

  // drag & drop support: drop
  // TODO: if dropped outside the ul, show more precisely that this was outside the range
  const drop = (event, target) => {
    event.dataTransfer.dropEffect = 'move'; 
    const start = parseInt(event.dataTransfer.getData("index"));
    const newList = $itemStore

    if (start < target) {
      newList.splice(target + 1, 0, newList[start]);
      newList.splice(start, 1);
    } else {
      newList.splice(target, 0, newList[start]);
      newList.splice(start + 1, 1);
    }
    // items = newList
    hovering = null

    // sync to server
    let reorder = {}
    let count = 0;
    newList.forEach(function(x, idx) {
      count ++
      reorder[x.id] = count
    })

    // call reorder from store
    reorderItemsInStore(reorder)

  }

  // drag & drop support: dragstart
  const dragstart = (event, index, itemid) => {
    event.dataTransfer.effectAllowed = 'move';
    //event.dataTransfer.dropEffect = 'move';
    event.dataTransfer.setData('index', index);
    event.dataTransfer.setData('itemid', itemid);
  }

  // get element styling edpending on item properties:
  const getStyle = (item) => {
    let style = "";
    if (item.shop) {
      style = "background-color: " + item.shop.color + ";"
    }
    // if we want to show the item, use shops bg color
    return style;
  }
</script>

<style>
.list {
  border-radius: 4px;
  box-shadow: 0 2px 3px rgba(10, 10, 10, 0.1), 0 0 0 1px rgba(10, 10, 10, 0.1);
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