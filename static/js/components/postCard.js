/* Postcard component loader
 * inputs:
 *          post: post object
 *          i: index of post in the html page*/

export const post_cards_component = (post, i) => {
  //parse cats
  let cats = ``;
  if (post.category === null) {
    cats += `<div class="category">null</div>`;
  } else {
    post.category.forEach((cat) => {
      cats += `<div class="category">${cat}</div>`;
    });
  }
  const postElement = document.createElement("div");
  postElement.className = "postcard";
  postElement.innerHTML = `
          <div class="postWrapper">
              <!-- <div class="postImage"></div> -->
              <div class="dataWrapper">
                  <div class="data">
                      <div class="title_category">
                          <a class="title bold_text" href='/post/${
                            post.post_id
                          }'>${post.title}</a>    
                          <div class="categories">
                                ${cats}
                          </div>
                      </div>
                      <div class="user">
                          <div class="userID">by ${post.user_name}</div>
                          <div class="action">
                              <p>Creation Date: ${post.creation_date}</p>
                          </div>
                      </div>
                  </div>
              </div>
              <div class="vote"> 
                    <div id="likeBtn_${post.post_id}" class="likeBtn ${
                      post.isLiked === 1 ? "liked" : ""
                    }" >
                      <img src="${
                        post.isLiked === 1
                          ? "static/assets/icons8-accept-30(1).png"
                          : "static/assets/icons8-accept-30.png"
                      }" alt="Like">
                    </div>
              <!-- Show like counts -->
              <div id="likes_${post.post_id}">
                    ${post.post_likes}
              </div>
                    <div id="dislikeBtn_${post.post_id}" class="dislikeBtn ${
                      post.isLiked === -1 ? "disliked" : ""
                    }" >
                        <img src="${
                          post.isLiked === -1
                            ? "static/assets/icons8-dislike-30(1).png"
                            : "static/assets/icons8-dislike-30.png"
                        }" alt="Dislike">
                    </div>
              <!-- Show like counts -->
              <div id="dislikes_${post.post_id}">
                    ${post.post_dislikes}
              </div>
              </div>
          </div>
          `;

  return postElement;
};
