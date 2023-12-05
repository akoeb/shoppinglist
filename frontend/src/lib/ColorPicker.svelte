<script>
  import { onMount, createEventDispatcher } from "svelte";
  const dispatch = createEventDispatcher();
  let activeColor = "#DAAFE9";
  let values = [
    [
      "#DAAFE9",
      "#C7DBF5",
      "#AAD5FB",
      "#ADE5DA",
      "#B0EDC3",
      "#FDF0A4",
      "#F8D6A2",
    ],
    [
      "#C47ADA",
      "#90BAEE",
      "#75BAFA",
      "#72D5BF",
      "#73DE8C",
      "#FBE66E",
      "#F5B969",
    ],
    [
      "#AE44B7",
      "#5E7ABC",
      "#5E7000",
      "#4DACA9",
      "#63B75A",
      "#EDBD4A",
      "#EC9740",
    ],
    [
      "#501B87",
      "#021B6B",
      "#0C2794",
      "#337277",
      "#2F6A52",
      "#AE802F",
      "#AD6127",
    ],
  ];
  onMount(() => dispatch("activeColor", activeColor));

  const changeActive = (value) => {
    console.log("changeActive color to " + value);
    activeColor = value;
    dispatch("activeColor", activeColor);
  };
</script>

<div>
  {#each values as val, oid}
    {#each val as innerValue, iid}
      <button
        style="background: {innerValue};"
        on:click|preventDefault={() => {
          changeActive(innerValue);
        }}
        class="button"
      >
      </button>
    {/each}
    <br />
  {/each}
</div>
