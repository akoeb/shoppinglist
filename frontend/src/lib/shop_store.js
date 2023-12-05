import { writable, get } from "svelte/store";
import { backend, httpOptions } from "../util";

// Get the value out of storage on load.
const data = JSON.parse(localStorage.getItem("shopStore"));

export let shopStore = writable(data || { items: [], version: 0 });

// Anytime the store changes, update the local storage value.
shopStore.subscribe((value) => {
  if (value.local) {
    delete value.local;
    const version = Date.now();
    value.version = version;
    // trigger sync with backend
    pushToBackend(value);
  }
  localStorage.setItem("shopStore", JSON.stringify(value));
});

// function to sync changes to backend and local store:
async function pushToBackend(shopList) {
  // then sync to backend:
  await fetch(backend("api/shops/sync"), httpOptions("POST", shopList)).catch(
    (err) => console.error(err)
  );
}

/*
// initialize or overwrite the store with contents from backend
// can be called from outside on mount or on receiving of change messages
export const loadAllShops = () => {
       
    fetch(backend('api/shops'), httpOptions())
      .then((res)  => res.json())
      .then((obj) => {
        shopStore.set(obj.shops ? obj.shops : [])
      })
      .catch((err) => console.error(err))
}
*/

/*
export const addShopToStore = (newShop) => {
    let shops = get(shopStore)
    // now send it to backend
    fetch(backend('api/shops'), httpOptions("POST", newShop))
    .then((res)  => res.json())
    .then((obj) => {
        shops.push(obj)
        shopStore.set(shops)
    })
    .catch((err) => console.error(err))        
}

export const assignItemToShop = (itemid, shopid) => {
  let url = backend('api/items/') + itemid + "/shop/" + shopid
  fetch(url , httpOptions("POST", {}))
  .catch((err) => console.error(err))        

}

// delete an item
export const deleteShopFromStore = (shopId) => {
  // get the position of the item in the array:
  let newList = get(shopStore).filter(x => x.id !== shopId)
  // send te delete request to backend, if successful remove the item from vue data array
  fetch(backend('api/shops/') + shopId, httpOptions("DELETE"))
  .then((res) => {
      shopStore.set(newList)
  })
  .catch((err) =>  console.error(err))
}
*/
