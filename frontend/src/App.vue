<template>
  <div id="app">
    <div class="">
    <draggable class="list-group" id="listgroup" v-model="shoppingItems" @end="syncOrder">
            <li class="list-group-item" v-for="item in shoppingItems" :id="getItemId(item.id)" :key="item.id" v-on:click="toggleStatus(item.id)" >
                <span v-bind:class="{'status-checked': isClosed(item) }">{{ item.title }}</span>
            <span class="pull-right">
                <button class="btn btn-xs btn-danger" v-on:click="deleteItem(item.id)">
                    <i class="fas fa-trash fa-1x" aria-hidden="true"></i>
                </button>
            </span>
        </li>
    </draggable>

    <div class="input-group">
        <input type="text"
            class="form-control"
            placeholder="New Item" 
            autofocus
            v-on:keyup.enter="addItem"
            v-model="newItem.title">
        <span class="input-group-btn">
                <button class="btn mybtn ml-1" type="button" v-on:click="addItem">CREATE</button>
        </span>
    </div><!-- /input-group -->
    </div>

  </div>
</template>

<script>
import Vue from "vue";
import draggable from 'vuedraggable'

function httpOptions(m, obj) {
    // this object will be returned
    var options = {
        method: m !== undefined ? m : "GET",
        mode: "cors",
        cache: "no-cache",
        credentials: "same-origin",
        headers: {
            "Content-Type": "application/json"
        },
        redirect: "follow",
    }
    if (obj !== undefined) {
        options['body'] = JSON.stringify(obj)
    }
    return options
}


export default Vue.extend({
//    name: 'app',

    // this defines the data on which vue listens
    data() {
        return {
            shoppingItems: [],
            newItem: {},
        }
    },
    // this happens on load of the page
    created: function() {
        this.setupStream()
        this.loadAll()
    },
    // define the methods
    methods: {
        // setup event stream for listening on updates
        setupStream: function() {
            let es = new EventSource('events');
            let that = this
            es.onmessage = function(event) {
                let data = JSON.parse(event.data);
                if (data.cmd == 'UPDATE') {
                    that.loadAll()
                }
            }

            es.addEventListener('error', event => {
                if (event.readyState == EventSource.CLOSED) {
                    console.log('Event was closed');
                    console.log(EventSource);
                }
            }, false);

        },
        // load all items
        loadAll: function() {
            fetch('items', httpOptions())
                .then((res)  => res.json())
                .then((json) => {
                    this.shoppingItems = json.items ? json.items : []
                })
                .catch((err) => console.error(err))
        },

        // create an item
        addItem: function() {
            // if title is empty we do not do anything
            if (!this.newItem.title.trim()) {
                this.newItem = {}
                return
            }
            // add item status and ordering
            this.newItem.status = "OPEN"
            // orderno must be the last of the shopping list plus 1
            this.newItem.orderno = this.shoppingItems.length

            // now send it to backend
            fetch('items', httpOptions("POST", this.newItem))
            .then((res)  => res.json())
            .then((json) => {
                // send successful, add the item to vue data array so it reflects on the page
                this.newItem.id = json.id
                this.shoppingItems.push(this.newItem)
                // this is not needed any more
                this.newItem = {}
            })
            .catch((err) => console.error(err))
        },

        // delete an item
        deleteItem: function(id) {
            // get the position of the item in the array:
            let index = this.shoppingItems.findIndex(x => x.id == id)
            // send te delete request to backend, if successful remove the item from vue data array
            fetch('items/' + id, httpOptions("DELETE"))
            .then((res) => this.shoppingItems.splice(index, 1))
            .catch((err) =>  console.log(error))
        },
        // update status
        toggleStatus: function(id) {
            // get a reference to the item with this id
            let item = this.shoppingItems.find(x => x.id === id)
            // copy the item by value ("clone")
            let toChange = Object.assign({}, item)
            // change status of the clone
            if (toChange.status === "CLOSED") {
                toChange.status = "OPEN"
            } else {
                toChange.status = "CLOSED"
            }
            // replace the object by the clone in the backend
            fetch('items/' + id, httpOptions("PUT", toChange))
            .then((res) => {
                // after the call to the backend was successful, change status in object so it reflects on the page
                item.status = toChange.status
            })
            .catch((err) => console.log(error))
        },
        // return true if the item status is closed
        isClosed: function(item) {
            return item.status === "CLOSED"
        },
        syncOrder: function(event) {
            // after resorting, we calculate order numbers based on new position
            // and sync those to the backend
            var reorder = {}
            var count = 0
            var max = this.shoppingItems.length
            this.shoppingItems.forEach(function(x, idx) {
                count ++
                reorder[x.id] = count
                if (count === max) {
                    // last element in the loop, send reorder array to backend:
                    options = httpOptions("POST", reorder)
                    fetch('items/reorder', options)
                    .then((res)  => res.json())
                    .then((json) => {
                        // successful backend request returns a orordered items array, reflect it in page:
                        this.shoppingItems = json.items ? json.items : []
                    })
                    .catch((err) => console.error(err))
                }
            })
        },
        getItemId: function(id) {
            // all others get an html id derived from their item position:
            return "item_" + this.shoppingItems.findIndex(x => x.id == id)
        },

    },
    // this is the end of the methods definition
    // in this block we define dynamic ("computed") properties
    computed: {
    }



});
</script>

<style lang="scss">


</style>
