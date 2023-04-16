<script>
  import { createEventDispatcher } from "svelte";

  const dispatch = createEventDispatcher();

  export let is_delete = false,
    is_edit = false,
    id,
    item = "",
    done_at,
    updated_at,
    created_at,
    is_new = false;

  let widget;

  $: get_widget_class = () =>
    ["collection-item", "avatar"]
      .concat(is_edit || is_delete ? ["active"] : [])
      .join(" ");

  $: is_editable = is_edit || is_new;

  let done_handler = () => {
    const data = {
      id: id,
      item: item,
      done_at: done_at ? undefined : new Date().toISOString(),
      updated_at: new Date().toISOString(),
      created_at: created_at,
    };

    fetch(`/api/${id}/`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    }).then(() => {
      dispatch("update", data);
    });
  };

  let confirm_handler = () => {
    if (is_new) {
      let data = {
        item: widget.value,
        done_at: undefined,
        updated_at: new Date().toISOString(),
        created_at: new Date().toISOString(),
      };

      fetch(`/api/`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      }).then((resp) => {
        dispatch("add", {
          id: resp.headers.get("location").split("/")[2],
          ...data,
        });
      });
    } else if (is_edit) {
      const data = {
        id: id,
        item: widget.value,
        done_at: done_at,
        updated_at: new Date().toISOString(),
        created_at: new Date().toISOString(),
      };

      fetch(`/api/${id}/`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      }).then(() => {
        dispatch("update", data);
      });
    } else if (is_delete) {
      fetch(`/api/${id}/`, {
        method: "DELETE",
      }).then(() => {
        dispatch("delete", { id: id });
      });
    }
  };

  let cancel_handler = () => {
    if (is_new) {
      dispatch("reset");
    } else if (is_edit) {
      dispatch("update", {
        id: id,
        item: item,
        done_at: done_at,
        updated_at: updated_at,
        created_at: created_at,
      });
    }
  };

  let edit_handler = () => {
    dispatch("active", {
      id: id,
      is_edit: true,
      is_delete: false,
    });
  };

  let delete_handler = () => {
    dispatch("active", {
      id: id,
      is_edit: false,
      is_delete: true,
    });
  };
</script>

<li class={get_widget_class()}>
  {#if !is_edit && !is_delete && !is_new}
    {#if done_at}
      <i class="material-icons circle" on:click={done_handler}>check_circle</i>
    {:else}
      <i class="material-icons circle" on:click={done_handler}
        >check_box_outline_blank</i
      >
    {/if}
  {/if}

  <div class="input-field col s12">
    <input bind:this={widget} id="todo" disabled={!is_editable} value={item} />
  </div>

  <span class="secondary-content">
    {#if is_new}
      <i class="material-icons" on:click={confirm_handler}>check</i>
      <i class="material-icons" on:click={cancel_handler}>clear</i>
    {:else if is_edit}
      <i class="material-icons">edit</i>
      <i class="material-icons" on:click={confirm_handler}>check</i>
      <i class="material-icons" on:click={cancel_handler}>clear</i>
    {:else if is_delete}
      <i class="material-icons">delete</i>
      <i class="material-icons" on:click={confirm_handler}>check</i>
      <i class="material-icons" on:click={cancel_handler}>clear</i>
    {:else}
      <i class="material-icons" on:click={edit_handler}>edit</i>
      <i class="material-icons" on:click={delete_handler}>delete</i>
    {/if}
  </span>
</li>
