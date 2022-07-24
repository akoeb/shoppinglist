<script>
    import { onMount } from 'svelte';
    import { backend } from './util';
    import { loadAllItems } from './lib/item_store.js'
    import ItemList from './lib/ItemList.svelte'
    import AddItem from './lib/AddItem.svelte';
    import ShopList from './lib/ShopList.svelte';
    import Icon from "svelte-awesome";
    import basket from 'svelte-awesome/icons/shoppingBasket'

    
    onMount( () => {
      setupStream()
	  });


    // Server-Sent Events:
    // setup event stream for listening on updates
    const setupStream = () => {
        let es = new EventSource(backend('events'));
        es.onmessage = function(event) {
            let data = JSON.parse(event.data);
            if (data.cmd == 'UPDATE') {
              // tell ItemList and ShopList components to reload their lists
              loadAllItems()
            }
        }

        es.addEventListener('error', event => {
            if (event.readyState == EventSource.CLOSED) {
                console.log('Event was closed');
                console.log(EventSource);
            }
        }, false);
    }



</script>


<main>

  <h1>Shopping List&nbsp;<Icon data={basket} scale="2"/></h1>

  <ItemList />
  <p class="columns column is-12"></p>
  <AddItem />
  <p class="columns column is-12"></p>
  <ShopList />
  
  

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
