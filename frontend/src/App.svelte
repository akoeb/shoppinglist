<script>
  import { onMount } from "svelte";
  import { backend } from "./util";
  import ItemList from "./lib/ItemList.svelte";
  import ShopList from "./lib/ShopList.svelte";
  import Icon from "svelte-awesome";
  import basket from "svelte-awesome/icons/shoppingBasket";

  let itemListComponent;
  let shopListComponent;
  let resetHovering;
  let filterItemsByShop;

  onMount(() => {
    // something has been dropped somewhere, reset all active classes:
    resetHovering = function () {
      itemListComponent.resetHovering();
      shopListComponent.resetHovering();
    };
    // a shop has been clicked, show only items
    // from this shop (or reset to show all)
    filterItemsByShop = function (ev) {
      itemListComponent.toggleFilterShopId(ev.detail);
    };
    // we initialize the sync from backend to both stores
    itemListComponent.getFromBackend();
    shopListComponent.getFromBackend();

    // and we setup an event stream to get notified if someone else changes something
    setupStream();
  });

  // Server-Sent Events:
  // setup event stream for listening on updates
  const setupStream = () => {
    let es = new EventSource(backend("events"));
    es.onmessage = function (event) {
      let data = JSON.parse(event.data);
      if (data.cmd == "UPDATE") {
        // tell ItemList and ShopList components to reload their lists
        itemListComponent.getFromBackend();
        shopListComponent.getFromBackend();
      }
    };

    es.addEventListener(
      "error",
      (event) => {
        if (event.readyState == EventSource.CLOSED) {
          console.log("Event was closed");
          console.log(EventSource);
        }
      },
      false
    );
  };
</script>

<svelte:body
  on:drop={resetHovering}
  on:dragover={(ev) => {
    ev.preventDefault();
  }}
/>
<main>
  <h1>Shopping List&nbsp;<Icon data={basket} scale="2" /></h1>

  <ItemList bind:this={itemListComponent} on:resetHovering={resetHovering} />
  <p class="columns column is-12"></p>
  <p class="columns column is-12"></p>
  <ShopList
    bind:this={shopListComponent}
    on:resetHovering={resetHovering}
    on:filterItemsByShop={filterItemsByShop}
  />
</main>

<style>
  main {
    text-align: center;
    padding: 1em;
    margin: 0 auto;
  }

  h1 {
    color: #ff3e00;
    text-transform: uppercase;
    font-size: 2rem;
    font-weight: 100;
    line-height: 1.1;
    margin: 2rem auto;
    max-width: 14rem;
  }

  @media (min-width: 480px) {
    h1 {
      max-width: none;
    }
  }
</style>
