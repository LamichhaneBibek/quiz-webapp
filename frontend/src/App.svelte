<script lang="ts">
    import Button from "./lib/Button.svelte";
    import QuizCard from "./lib/QuizCard.svelte";
    import { NetService } from "./service/net";
    import type { Quiz, QuizQuestion } from "./model/quiz";

  let quizzes :{_id: string, name: string}[] = [];
  let currentQuestion: QuizQuestion | null = null;

  let netService = new NetService();
  netService.connect();
  netService.onPacket((packet) => {
    console.log(packet);
    switch (packet.id) {
      case 2:{
        currentQuestion = packet.question;
        break;
      }
    }
  });

  async function getQuizzes() {
    let response = await fetch('http://localhost:8000/api/quizzes');
    if (!response.ok) {
      alert('Failed to fetch quizzes');
      return;
    }
    let json = await response.json();
    quizzes = json;
  }

  let code = "";
  let msg = "";

  function connect (){
    netService.sendPacket({
      id: 0,
      code: "1234",
      name: "kebib"
    })
  }

  function hostQuiz(quiz: Quiz){
    netService.sendPacket({
      id: 1,
      quizId: quiz.id
    })
  }

</script>

<Button on:click={getQuizzes}> Get Quizzes</Button>
Message: {msg}

<div>
  {#each quizzes as quiz}
    <QuizCard on:host={() => hostQuiz(quiz)} quiz = {quiz} />
  {/each}
</div>
<input bind:value={code} type="text" class="border" placeholder="Game code" />
<Button on:click={connect}>Join game</Button>

{#if currentQuestion != null}
 <h2 class="text-4xl font-bold mt-8">{currentQuestion.name}</h2>
 <div class="fle">
    {#each currentQuestion.choices as choice}
     <div class="flex-1 bg-blue-400 text-center font-bold text-2xl text-white justify-venter items-center p-8">
        {choice.name}
     </div>
    {/each}
 </div>
{/if}


