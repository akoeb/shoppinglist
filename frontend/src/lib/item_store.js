import { writable } from "svelte/store";
import { httpOptions, backend, notifyWhenChanges } from "../util";

// Get the value out of storage on load.
const data = JSON.parse(localStorage.getItem("itemStore"));

/*

FIXME: 
the stores are called either from local change, in this case we want to set a new 
version in localstore and push it to remote
or they are called from remote change, in this case we want to keep local 
version and do *NOT* want to trigger all those things in the subscribe method


// https://stackoverflow.com/questions/71570197/how-to-not-trigger-svelte-store-subscription-if-the-value-was-not-changed
// Set the stored value or a sane default.
export const itemStore = writable(data || []);
 function notifyWhenChanges(store, callback) {
   let init = false;
   return store.subscribe((value) => {
     if (init) callback(value);
     else init = true;
   });
 }
 */

// Set the stored value or a sane default.
export const itemStore = writable(data || { items: [], version: 0 });

// Anytime the store changes, update the local storage value.
itemStore.subscribe((value) => {
  if (value.local) {
    delete value.local;
    const version = Date.now();
    value.version = version;
    // trigger sync with backend
    pushToBackend(value);
  }
  localStorage.setItem("itemStore", JSON.stringify(value));
});

// function to sync changes to backend and local store:
async function pushToBackend(itemList) {
  // then sync to backend:
  await fetch(backend("api/items/sync"), httpOptions("POST", itemList)).catch(
    (err) => console.error(err)
  );
}

/*
// subscribe to store to keep backend in sync on changes;
var unsubscribe;

// load all Items from server
export const loadAllItems = (shopId = 0) => {
  if (unsubscribe) {
    unsubscribe();
  }
  fetch(backend("api/items"), httpOptions())
    .then((res) => res.json())
    .then((obj) => {
      let items = [];
      if (shopId > 0) {
        // filter items
        items = obj.items.filter((item) => {
          return item.shop && item.shop.id === shopId;
        });
      } else {
        items = obj.items ? obj.items : [];
      }
      itemStore.set(items);
      unsubscribe = itemStore.subscribe(sync);
    })
    .catch((err) => console.error(err));
};
*/

/*
export const addItemToStore = (newItem) => {
    let items = get(itemStore)
    // now send it to backend
    fetch(backend('api/items'), httpOptions("POST", newItem))
    .then((res)  => res.json())
    .then((obj) => {
        // send successful, add the item to the items array so it reflects on the page
        items.push(obj)
        itemStore.set(items)
    })
    .catch((err) => console.error(err))
    
}

// reorder the items after drag&drop
export const reorderItemsInStore = (reorder) => {
    fetch(backend('api/items/reorder'), httpOptions("POST", reorder))
    .then((res)  => res.json())
    .then((obj) => {
      // successful backend request returns a ordered items array, reflect it in page:
      itemStore.set(obj.items ? obj.items : [])
    })
    .catch((err) => console.error(err))
}

// delete an item
export const deleteItemFromStore = (id) => {
    // get the position of the item in the array:
    let items = get(itemStore)
    let index = items.findIndex(x => x.id == id)
    // send te delete request to backend, if successful remove the item from vue data array
    fetch(backend('api/items/') + id, httpOptions("DELETE"))
    .then((res) => {
        let new_items = items.splice(index, 1);
        itemStore.set(new_items ? new_items : [])
    })
    .catch((err) =>  console.error(err))
}

// update status
export const toggleStatusInStore = (id) => {
    // get a reference to the item with this id
    let item = get(itemStore).find(x => x.id === id)
    // copy the item by value ("clone")
    let toChange = Object.assign({}, item)
    // change status of the clone
    if (toChange.status === "CLOSED") {
        toChange.status = "OPEN"
    } else {
        toChange.status = "CLOSED"
    }
    // replace the object by the clone in the backend
    fetch(backend('api/items/' + id), httpOptions("PUT", toChange))
    .then((res) => {
        // after the call to the backend was successful, change status in object so it reflects on the page
        item.status = toChange.status
    })
    .catch((err) => console.error(err))
}
*/
