<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import Button from "./Button.svelte";

    const dispatch = createEventDispatcher();

    export let quiz: {_id: string, name: string};
    function host(){
        dispatch("host");
    }
</script>

<div class="flex justify-between items-center bg-white border-2 border-gray-300 p-4 m-2 rounded-xl">
   <p>
       {quiz.name}
   </p> 
   <Button on:click={host}>
       Host
    </Button>
</div>