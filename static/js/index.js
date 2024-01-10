import { loadNav, isloggedIn } from "./components/navbar.js";
import { post_cards_component } from "./components/postCard.js";

// global_vars
let gotten_posts = [];

// A function to check if logged in and executes another function
async function evalLogin(fn) {
  let islogged = await isloggedIn();
  if (islogged === "true") {
    fn();
  } else {
    window.location.replace("/login");
  }
}

// App entry
const render_index_page = async () => {
  let nav = await loadNav("/"); // loads navbar
  let body = document.body; // get the html body
  body.insertAdjacentHTML("beforebegin", nav); // attach nav bar to html body

  // Fetch the JSON data from the URL
  fetch("/api/posts")
    .then((response) => response.json())
    .then((data) => {
      // Process the JSON data and create HTML elements
      const jsonContainer =
        document.getElementsByClassName("postcardwrapper")[0];
      // console.log(jsonContainer);
      let i = 0;
      data.posts.forEach((post) => {
        gotten_posts.push(post);

        // parse post
        let postElement = post_cards_component(post, i);
        i++;

        jsonContainer.appendChild(postElement);
      });
    })
    .then(() => {
      /********* Add click listeners ***************/
      const likeButtons = document.querySelectorAll(".likeBtn");
      const dislikeButtons = document.querySelectorAll(".dislikeBtn");
      likeButtons.forEach((btn, index) => {
        btn.addEventListener("click", () =>
          evalLogin(() => LikeEvent(index, btn.id.split("_")[1])),
        );
      });

      dislikeButtons.forEach((btn, index) => {
        console.log(index);
        btn.addEventListener("click", () =>
          evalLogin(() => disLikeEvent(index, btn.id.split("_")[1])),
        );
      });
      /********* END of Adding click listeners ***************/
    })
    .catch((error) => console.error("Error fetching JSON:", error));
};

const loadCats = async () => {
  const catwrapper = document.getElementById("allcats");
  let response = await fetch("/api/categories");
  let data = await response.json();
  data.Categories.forEach((cat) => {
    let child = document.createElement("div");
    child.classList.add("catlisting");
    child.id = `catlisting-${cat.category}`;
    child.addEventListener("click", () => {
      filterBy("/api/category/" + cat.category);
    });
    child.innerText = cat.category;
    catwrapper.append(child);
  });
};

const filterByUser = async () => {
  const createdByUser = document.getElementById("createdByUser");
  createdByUser.addEventListener("click", () => {
    filterBy("/api/created_by_user");
    return;
  });
  const likedByUser = document.getElementById("likedByUser");
  likedByUser.addEventListener("click", () => {
    filterBy("/api/liked_by_user");
    return;
  });
};

const filterBy = async (path) => {
  const jsonContainer = document.getElementsByClassName("postcardwrapper")[0];
  jsonContainer.innerHTML = ``;
  let response = await fetch(path);
  let data = await response.json();
  let i = 0;
  data.posts.forEach((post) => {
    if (post.category === null) {
      return;
    } else {
      let postElement = post_cards_component(post, i);
      i++;
      jsonContainer.appendChild(postElement);
    }
  });

  /********* Add click listeners ***************/
  const likeButtons = document.querySelectorAll(".likeBtn");
  const dislikeButtons = document.querySelectorAll(".dislikeBtn");
  likeButtons.forEach((btn, index) => {
    btn.addEventListener("click", () =>
      evalLogin(() => LikeEvent(index, btn.id.split("_")[1])),
    );
  });

  dislikeButtons.forEach((btn, index) => {
    console.log(index);
    btn.addEventListener("click", () =>
      evalLogin(() => disLikeEvent(index, btn.id.split("_")[1])),
    );
  });
  /********* END of Adding click listeners ***************/
};



// Fetches post data from the specified JSON file path and displays the posts
const dispayPostsCards = async (path) => {
  const jsonContainer = document.getElementsByClassName("postcardwrapper")[0];
  jsonContainer.innerHTML = ``;
  let response = await fetch(path);
  let data = await response.json();
  let i = 0;
  data.posts.forEach((post) => {
    if (post.category === null) {
      return;
    } else {
      let postElement = post_cards_component(post, i);
      i++;
      jsonContainer.appendChild(postElement);
    }
  });
};
/////////////////////////////////////////////////////////////////////

const LikeEvent = async (index, postID) => {
  let likeBtn = document.querySelectorAll(".likeBtn")[index]; // selects like button
  let dislikeBtn = document.querySelectorAll(".dislikeBtn")[index]; // selects dislike button
  let like_count_area = document.getElementById(`likes_${postID}`); // selects like counts area from dom
  let dislike_count_area = document.getElementById(`dislikes_${postID}`);

  if (gotten_posts[index].isLiked === 1) {
    gotten_posts[index].isLiked = 0;
    likeBtn.classList.remove("liked");
    likeBtn.innerHTML =
      '<img src="static/assets/icons8-accept-30.png" alt="Like">';
    await sendReqPost(postID, 0);
  } else {
    gotten_posts[index].isLiked = 1;
    likeBtn.classList.add("liked");
    likeBtn.innerHTML =
      '<img src="static/assets/icons8-accept-30(1).png" alt="Like">';

    await sendReqPost(postID, 1);
  }

  if (dislikeBtn.classList.contains("disliked")) {
    dislikeBtn.classList.remove("disliked");
    dislikeBtn.innerHTML =
      '<img src="static/assets/icons8-dislike-30.png" alt="Dislike">';
  }
  /* fetch new like and dislike count and update the DOM */
  await get_like_dislike_count(postID).then((likes_dislikes) => {
    like_count_area.innerHTML = likes_dislikes.interactions.like_count;
    dislike_count_area.innerHTML = likes_dislikes.interactions.dislike_count;
  });
};

const disLikeEvent = async (index, postID) => {
  let likeBtn = document.querySelectorAll(".likeBtn")[index]; // selects like button
  let dislikeBtn = document.querySelectorAll(".dislikeBtn")[index]; // selects dislike button
  let like_count_area = document.getElementById(`likes_${postID}`); // selects like counts area from dom
  let dislike_count_area = document.getElementById(`dislikes_${postID}`);

  if (gotten_posts[index].isLiked === -1) {
    gotten_posts[index].isLiked = 0;
    dislikeBtn.classList.remove("disliked");
    dislikeBtn.innerHTML =
      '<img src="static/assets/icons8-dislike-30.png" alt="Dislike">';
    await sendReqPost(postID, 0);
  } else {
    gotten_posts[index].isLiked = -1;
    dislikeBtn.classList.add("disliked");
    dislikeBtn.innerHTML =
      '<img src="static/assets/icons8-dislike-30(1).png" alt="Dislike">';

    await sendReqPost(postID, -1);
  }

  if (likeBtn.classList.contains("liked")) {
    likeBtn.classList.remove("liked");
    likeBtn.innerHTML =
      '<img src="static/assets/icons8-accept-30.png" alt="Like">';
  }
  /* fetch new like and dislike count and update the DOM */
  await get_like_dislike_count(postID).then((likes_dislikes) => {
    like_count_area.innerHTML = likes_dislikes.interactions.like_count;
    dislike_count_area.innerHTML = likes_dislikes.interactions.dislike_count;
  });
};

const get_like_dislike_count = async (postID) => {
  let interactions_obj = {};
  await fetch("/api/postlikes", {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      post_id: postID,
    },
  })
    .then((resp) => {
      return resp.json();
    })
    .then((data) => {
      interactions_obj = data;
    });

  return interactions_obj;
};

const sendReqPost = async (PostID, LikeDislike) => {
  await fetch("/api/likes_post", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      PostID: parseInt(PostID, 10),
      LikeDislike,
    }),
  });
};

/************* App entry point *************/
// Attach the function to the load event

async function initPages() {
  await loadCats();
  await filterByUser();
  // await render_nav_bar();
  //  isloggedIn();
  await render_index_page();
}

// window.addEventListener("load", initPages, true);
window.addEventListener("DOMContentLoaded", initPages);
/************* END *************/
