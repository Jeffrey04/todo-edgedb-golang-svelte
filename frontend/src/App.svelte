<script>
  import "./app.css";

  import Todo from "./lib/Todo.svelte";
  import { onMount } from "svelte";

  //let todos = [
  //    {
  //      id: 1,
  //      item: "wash the car",
  //      done_at: Date.now(),
  //      created_at: Date.now(),
  //      updated_at: Date.now(),
  //    },
  //    {
  //      id: 2,
  //      item: "feed the cat",
  //      done_at: undefined,
  //      created_at: undefined,
  //      updated_at: undefined,
  //    },
  //    {
  //      id: 3,
  //      item: "clean the desk",
  //      done_at: undefined,
  //      created_at: undefined,
  //      updated_at: undefined,
  //    },
  //  ],
  let todos = [],
    todo_new = {
      id: -1,
      item: "",
      done_at: undefined,
      created_at: undefined,
      updated_at: undefined,
    };

  let active = undefined;

  let active_handler = (event) => {
    active = event.detail;
  };

  let add_handler = (event) => {
    todos = todos.concat([event.detail]);
    todo_new = todo_new;
  };

  let update_handler = (event) => {
    todos = todos.reduce(
      (current, incoming) =>
        current.concat([
          incoming.id == event.detail.id ? event.detail : incoming,
        ]),
      []
    );

    active = undefined;
  };

  let delete_handler = (event) => {
    todos = todos.filter((incoming) => incoming.id != event.detail.id);

    active = undefined;
  };

  let reset_handler = () => {
    todo_new = todo_new;
  };

  onMount(async () => {
    fetch("/api/")
      .then((data) => {
        return data.json();
      })
      .then((data) => {
        todos = data;
      });
  });
</script>

<main>
  <div class="container">
    <h1>Todo Example App</h1>

    <ul class="collection">
      {#key todos}
        {#each todos as todo}
          <Todo
            on:active={active_handler}
            on:update={update_handler}
            on:delete={delete_handler}
            is_edit={active?.id === todo.id && active?.is_edit}
            is_delete={active?.id === todo.id && active?.is_delete}
            {...todo}
          />
        {/each}
      {/key}
      {#if !active}
        {#key todo_new}
          <Todo
            on:reset={reset_handler}
            on:add={add_handler}
            is_new={true}
            {...todo_new}
          />
        {/key}
      {/if}
    </ul>
  </div>
</main>
