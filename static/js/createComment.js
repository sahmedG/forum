document.addEventListener("DOMContentLoaded", () => {
  // Handle form submission
  const commentForm = document.getElementById("newCommentForm");
  commentForm.addEventListener("submit", async (event) => {
    event.preventDefault(); // Prevent the default form submission behavior

    let islogged = localStorage.getItem("isloggedIn");
    if (islogged === "false") {
      window.location.replace("/login");
      return;
    }

    // Get the comment text from the textarea
    const commentText = document.getElementById("newCommentText").value;

    // Check for empty title and post content
    if (!commentText.trim()) {
      alert("Comment con not be empty");
      return;
    }
    
    // Check the length of the comment text
    if (commentText.length > 500) {
      alert("Comment cannot exceed 500 characters.");
      return;
    }

    // addComment function with the comment text
    try {
      const response = await fetch("../../api/add_comment", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          post_id: postID,
          content: commentText,
        }),
      });
      let done = 0;
      if (response.ok) {
        const newComment = await response.json();

        // Update the UI with the newly added comment
        const commentDiv = document.querySelector(".postcomments .comment");
        const newCommentHTML = await orgComments([newComment]);
        commentDiv.innerHTML += newCommentHTML;
        done = 1;
      } else {
        const errorText = await response.text();
        alert(errorText);
        return;
      }
    } catch (error) {
      console.error("Error adding comment:", error);
    }

    if ((done = 1)) {
      // reloading the page to show the new comment
      location.reload();
      // Clear the textarea after adding the comment
      document.getElementById("newCommentText").value = "";
    }
  });
});

// get the post ID
const postID = parseInt(
  location.href.match(/post\/[0-9]+/)[0].replace("post/", "")
);
