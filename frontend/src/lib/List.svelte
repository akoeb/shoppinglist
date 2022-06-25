
<div>
  <form on:submit|preventDefault={addItem}>
    <label for="name">Add an item</label>
    <input id="name" type="text" bind:value={newItem.title} />
  </form>

  <ul class="list">
    {#each items as item, index (item.title)}
      <li 
        class="list-item"
        on:click="{toggleStatus(item.id)}" 
        draggable={true}
        on:dragstart={event => dragstart(event, index)}
        on:drop|preventDefault={event => drop(event, index)}
        ondragover="return false"
        on:dragenter={() => hovering = index}
        class:is-active={hovering === index}>

        <!-- input type="checkbox" bind:checked={item.done} / -->
        <span class:checked="{item.status === 'CLOSED'}">{item.title}</span>
        <button on:click={() => deleteItem(item.id)}>&times;</button>
      </li>
    {:else}
	    <p>The list is empty</p>
    {/each}
  </ul>
</div>


<script>
// @ts-nocheck
import { onMount } from 'svelte';
import {flip} from 'svelte/animate';
import  {httpOptions, backend} from '../util';


  // the list of items to be bought:
  let items = [];
  // one temporary item to be used for creation
  let newItem = {};
  // some variable for drg&drop
  let hovering = false;

  /*
   * Lifecycle: on Mount of the app
   */

  onMount( () => {
		loadAll()
	});


  /*
  ** functions:
  */

  // drag & drop support: drop
  // TODO: if dropped outside the ul, show more precisely that this was outside the range
  const drop = (event, target) => {
    event.dataTransfer.dropEffect = 'move'; 
    const start = parseInt(event.dataTransfer.getData("text/plain"));
    const newTracklist = items

    if (start < target) {
      newTracklist.splice(target + 1, 0, newTracklist[start]);
      newTracklist.splice(start, 1);
    } else {
      newTracklist.splice(target, 0, newTracklist[start]);
      newTracklist.splice(start + 1, 1);
    }
    items = newTracklist
    hovering = null

    // sync to server
    let reorder = {}
    let count = 0;
    items.forEach(function(x, idx) {
      count ++
      reorder[x.id] = count
    })

    fetch(backend('items/reorder'), httpOptions("POST", reorder))
    .then((res)  => res.json())
    .then((json) => {
      // successful backend request returns a orordered items array, reflect it in page:
      //items = json.items ? json.items : []
    })
    .catch((err) => console.error(err))

  }

  // drag & drop support: dragstart
  const dragstart = (event, i) => {
    event.dataTransfer.effectAllowed = 'move';
    event.dataTransfer.dropEffect = 'move';
    const start = i;
    event.dataTransfer.setData('text/plain', start);
  }

  // load the complete list from server, called by the onmount lofecycle hook
  const loadAll = () => {
    // @ts-ignore
    fetch(backend('items'), httpOptions())
        .then((res)  => res.json())
        .then((json) => {
            items = json.items ? json.items : []
        })
        .catch((err) => console.error(err))
  }

  // create an item
  const addItem = () => {
    // if title is empty we do not do anything
    if (!newItem.title.trim()) {
        newItem = {}
        return
    }
    // add item status and ordering
    newItem.status = "OPEN"
    // orderno must be the last of the shopping list plus 1
    newItem.orderno = items.length

    // now send it to backend
    fetch(backend('items'), httpOptions("POST", newItem))
    .then((res)  => res.json())
    .then((json) => {
        // send successful, add the item to the items array so it reflects on the page
        items.push(json)
        // assign array to itself to trigger svelte reactivity:
        items = items;
        // this is not needed any more
        newItem = {}
    })
    .catch((err) => console.error(err))
  }
  
  // delete an item
  const deleteItem = (id) => {
      // get the position of the item in the array:
      let index = items.findIndex(x => x.id == id)
      // send te delete request to backend, if successful remove the item from vue data array
      fetch(backend('items/') + id, httpOptions("DELETE"))
      .then((res) => {
        items.splice(index, 1);
        // assign array to itself to trigger svelte reactivity:
        items = items;
      })
      .catch((err) =>  console.log(err))
  }

  // update status
  const toggleStatus = (id) => {
      // get a reference to the item with this id
      let item = items.find(x => x.id === id)
      // copy the item by value ("clone")
      let toChange = Object.assign({}, item)
      // change status of the clone
      if (toChange.status === "CLOSED") {
          toChange.status = "OPEN"
      } else {
          toChange.status = "CLOSED"
      }
      // replace the object by the clone in the backend
      fetch(backend('items/' + id), httpOptions("PUT", toChange))
      .then((res) => {
          // after the call to the backend was successful, change status in object so it reflects on the page
          item.status = toChange.status
          // assign array to itself to trigger svelte reactivity:
          items = items;
      })
      .catch((err) => console.log(err))
  };

</script>

<style>
.list {
  background-color: white;
  border-radius: 4px;
  box-shadow: 0 2px 3px rgba(10, 10, 10, 0.1), 0 0 0 1px rgba(10, 10, 10, 0.1);
}

.list-item {
  display: block;
  padding: 0.5em 1em;
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