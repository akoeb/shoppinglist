<script>
  import { onMount } from "svelte";
  import ColorPicker from "./ColorPicker.svelte";
  import { shopStore } from "./shop_store";
  import { uid } from "../util";
  let newShop = { name: "" };
  export let showModal = false;
  let dialog; // HTML Dialog element
  $: if (dialog && showModal) dialog.showModal();

  // focus on input field
  onMount(() => inputFocus());

  // create an item
  const addShop = () => {
    // if title is empty we do not do anything
    if (!newShop.name.trim()) {
      newShop = {};
      return;
    }
    let store = $shopStore;
    newShop.orderno = store.length + 1;
    newShop.uid = uid();
    store.items.push(newShop);
    store.local = true;
    $shopStore = store;
    newShop = { name: "", color: newShop.color };
    inputFocus();
  };

  const changeColor = (ev) => {
    newShop.color = ev.detail;
    inputFocus();
  };

  // set focus on input field upon creation
  function inputFocus() {
    document.getElementById("newShopName").focus();
  }
</script>

<dialog
  bind:this={dialog}
  on:close={() => (showModal = false)}
  on:click|self={() => dialog.close()}
>
  <br />
  <form on:submit|preventDefault={addShop}>
    <input
      id="newShopName"
      type="text"
      class="input"
      bind:value={newShop.name}
      placeholder="New Shop"
      style="background-color: {newShop.color}"
    />
  </form>
  <ColorPicker on:activeColor={changeColor} />
  <button class="button is-link column is-4" on:click|preventDefault={addShop}
    >Submit</button
  >
  <button class="close" on:click={() => dialog.close()}>&times;</button>
</dialog>

<style>
  dialog::backdrop {
    background: rgba(0, 0, 0, 0.3);
  }
  .close {
    color: #aaaaaa;
    position: absolute;
    top: 0;
    right: 0;
    /* float: right;*/
    font-size: 28px;
    font-weight: bold;
  }

  .close:hover,
  .close:focus {
    color: #000;
    text-decoration: none;
    cursor: pointer;
  }
</style>
