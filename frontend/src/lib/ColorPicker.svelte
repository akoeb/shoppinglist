<div>
<button class="button" style="color: {activeColor}" on:click="{toggleShowValues()}"></button>
<input class="input" value="{activeColor}" />
{#if showValues}
    {#each values as val, index}
        {#each val as innerValue, innerIndex}
            <button
                style="background: {innerValue};"
                on:click={() => { changeActive(innerValue) }}
                class="button">
            </button>
        {/each}
    {/each}
{/if}


</div>

<script>
    export let activeColor = '#DAAFE9'
    let showValues = false;
    let values = [
		['#DAAFE9', '#C7DBF5', '#AAD5FB', '#ADE5DA', '#B0EDC3', '#FDF0A4', '#F8D6A2'],
		['#C47ADA', '#90BAEE', '#75BAFA', '#72D5BF', '#73DE8C', '#FBE66E', '#F5B969'],
		['#AE44B7', '#5E7ABC', '#5E7ABC', '#4DACA9', '#63B75A', '#EDBD4A', '#EC9740'],
		['#501B87', '#021B6B', '#0C2794', '#337277', '#2F6A52', '#AE802F', '#AD6127']
    ]

    const toggleShowValues = () => {
        if (showValues) {
            showValues = false;
        } else {
            showValues = true;
        }
    }

    const changeActive = (value) => {
        activeColor = value;
    }

</script>

