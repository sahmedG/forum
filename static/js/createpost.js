import { loadNav } from "./components/navbar.js";

// Loads categories selection buttons
const loadOptions = async () => {
  const catdiv = document.getElementById("catHandler");
  console.log(catdiv);
  let response = await fetch("/api/categories");
  let data = await response.json();
  console.log(data);
  data.Categories.forEach((cat) => {
    catdiv.innerHTML += `
    <div class="checkoption">
              <label class="checkcontainer">
                <input value="${cat.category}" type="checkbox" />
                <svg width="1em" height="1em" viewBox="0 0 64 64">
                  <path
                    class="path"
                    pathLength="575.0541381835938"
                    d="M 0 16 V 56 A 8 8 90 0 0 8 64 H 56 A 8 8 90 0 0 64 56 V 8 A 8 8 90 0 0 56 0 H 8 A 8 8 90 0 0 0 8 V 16 L 32 48 L 64 16 V 8 A 8 8 90 0 0 56 0 H 8 A 8 8 90 0 0 0 8 V 56 A 8 8 90 0 0 8 64 H 56 A 8 8 90 0 0 64 56 V 16"
                  ></path>
                </svg>
              </label>
              <div class="option">${cat.category}</div>
            </div>
    `;
  });
};

const createPost = async () => {
  event.preventDefault();
  try {
    let formData = {
      Post: document.getElementById("Pcontent").value,
      Title: document.getElementById("Ptitle").value,
    };

    // Check for empty title and post content
    if (!formData.Post.trim() || !formData.Title.trim()) {
      let errdiv = document.getElementById("errdiv");
      errdiv.innerText = "Title and Post Content are required fields";
      return;
    }

    // Check the length of the post content
    if (formData.Post.length > 10000) {
      let errdiv = document.getElementById("errdiv");
      errdiv.innerText = "Post Content should be up to 10000 characters long";
      return;
    }

    // Check the length of the title
    if (formData.Title.length > 100) {
      let errdiv = document.getElementById("errdiv");
      errdiv.innerText = "Title should be up to 100 characters long";
      return;
    }

    // get selected cats
    let maincats = [];

    // handling cats
    let checkboxes = document.querySelectorAll(
      'input[type="checkbox"]:checked',
    );

    if (checkboxes.length === 0) {
      maincats.push("General");
    } else {
      checkboxes.forEach((checkbox) => {
        maincats.push(checkbox.value);
      });
    }

    formData["Categories"] = maincats;

    console.log(formData);

    const response = await fetch("/api/create_post", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(formData),
    });

    if (response.ok) {
      console.log("Post created successfully");
      window.location.replace("/");
    } else {
      const errorText = await response.text();
      let errdiv = document.getElementById("errdiv");
      errdiv.innerText = errorText;
      console.error(
        `Failed to create post. Server returned ${response.status} status.`,
      );
    }
  } catch (error) {
    console.error("An error occurred while creating the post:", error);
  }
};

// App enrty point
const loadPage = async () => {
  loadOptions();
  let nav = await loadNav("/");
  let body = document.body;
  body.insertAdjacentHTML("beforebegin", nav);

  /* Add click listeners for button */
  document.getElementById("postBtn").addEventListener("click", createPost);
};

// document.addEventListener("load", loadPage, true);
document.addEventListener("DOMContentLoaded", loadPage);
