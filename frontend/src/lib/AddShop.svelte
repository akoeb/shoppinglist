<div id="showAddShopForm">
    <input id="newShopName" type="text" class="input" bind:value={newShop.name} placeholder="New Shop" />
    <ColorPicker bind:activeColor={newShop.color} />
    <button class="button is-link column is-2" on:click="{addShop}">Submit</button>
</div>

<script>
    import ColorPicker from './ColorPicker.svelte'
    import {shopStore, addShopToStore} from './shop_store'
    let newShop = {name: "", color: '#DAAFE9'}
    
    // create an item
    const addShop = () => {
        // if title is empty we do not do anything
        if (!newShop.name.trim()) {
            newShop = {}
            return
        }
        newShop.orderno = $shopStore.length + 1
        addShopToStore(newShop)
        newShop = {name: "", color: '#DAAFE9'}

  }

</script>

<style>

</style>