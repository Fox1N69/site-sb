export async function prerender(data) {
  return {
    props: {
      postEndpoint: "/posts/set",
    },
  };
}

import { createState, onMount } from "svelte";

const formData = createState({
  Title: "",
  Description: "",
  PostBody: "",
});

async function handleSubmit(event) {
  event.preventDefault();

  try {
    const response = await fetch(
      new Request($postEndpoint, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(formData),
      })
    );

    if (response.ok) {
      console.log("Post successfully sent!");
    } else {
      console.error("Failed to send post.");
    }
  } catch (error) {
    console.error("Error sending post:", error);
  }
}

$: postEndpoint = "";
onMount(() => {
  postEndpoint = $page.props.postEndpoint;
});
