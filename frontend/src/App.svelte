<script>
  import List from './lib/List.svelte'
  import Icon from "svelte-awesome";
  import basket from 'svelte-awesome/icons/shoppingBasket'
</script>


<main>

  <h1>Shopping List&nbsp;<Icon data={basket} scale="2"/></h1>

  <List />

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
