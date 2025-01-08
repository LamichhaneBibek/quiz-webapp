<script lang="ts">
    import Button from "./lib/Button.svelte";
import QuizCard from "./lib/QuizCard.svelte";
  let quizzes :{_id: string, name: string}[] = [];
  async function getQuizzes() {
    let response = await fetch('http://localhost:8000/api/quizzes');
    if (!response.ok) {
      throw new Error('Failed to fetch quizzes');
    }
    let json = await response.json();
    quizzes = json;
  }

  let code = "";

  function connect (){
    let websocket = new WebSocket('ws://localhost:8000/ws');
    websocket.onopen = () => {
      console.log('open connection');
      websocket.send('Hello');
    };

    websocket.onmessage = (event) => {
      console.log(event.data);
    };
  }

  function hostQuiz(quiz){
    let websocket = new WebSocket('ws://localhost:8000/ws');
    websocket.onopen = () => {
      console.log('open connection');
      websocket.send(JSON.stringify({type: 'host', quizId: quiz._id}));
    };

    websocket.onmessage = (event) => {
      console.log(event.data);
    };
  }

</script>

<button on:click={getQuizzes}> Get Quizzes</button>

<div>
  {#each quizzes as quiz}
    <QuizCard on:host={() => hostQuiz(quiz)} quiz = {quiz} />
  {/each}
</div>
<input type="text" class="border" placeholder="Game code" />
<Button on:click={connect}>Join game</Button>



