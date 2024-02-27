// function to add post divs to post wrapper
/*
 * inputs:
 *    prop: json object, in this case it is post object
 *    wrapper: is the html tag that holds the rendered infrormations
 *    */

export const orgPostHTML = async (wrapper, prop) => {
  let cats = ``;

  console.log(prop);

  console.log(prop.category);
  prop.category.forEach((cat) => {
    console.log(cat);
    cats += `<div class="category">${cat}</div>`;
  });

  let [comments, gotten_comm] = await orgComments(prop.post_id);

  wrapper.innerHTML += `
    <div class="profilestuff">
      <div class="pfpImage">
        <img src="../../static/assets/reddit.png" alt="reddit lol" class="pimg">
      </div>
      <div class="profileinfo">
        <div class="profileName">${prop.user_name}</div>
        <div class="postinfo">
          <div class="postDate">${prop.creation_date}</div>
          <div class="commentsLink">Comments ${gotten_comm.length}</div>

          <div id="likeBtn_${prop.post_id}" class="likeBtn ${
            prop.isLiked === 1 ? "liked" : ""
          }" >
            <img src="${
              prop.isLiked === 1
                ? "../../static/assets/icons8-accept-30(1).png"
                : "../../static/assets/icons8-accept-30.png"
            }" alt="Like">
          </div>
          <!-- Show like counts -->
          <div id="likes_${prop.post_id}">
                ${prop.post_likes}
          </div>
                <div id="dislikeBtn_${prop.post_id}" class="dislikeBtn ${
                  prop.isLiked === -1 ? "disliked" : ""
                }" >
                    <img src="${
                      prop.isLiked === -1
                        ? "../../static/assets/icons8-dislike-30(1).png"
                        : "../../static/assets/icons8-dislike-30.png"
                    }" alt="Dislike">
                </div>
          <!-- Show like counts -->
          <div id="dislikes_${prop.post_id}">
                ${prop.post_dislikes}
          </div>

        </div>

      </div>
    </div>
    <div class="posttitle">${prop.title}</div>
    <div class="postcats">
      ${cats}
    </div>
    <hr>
    <div class="image-viewing-area" id="imageViewingArea"></div> 
    <div class="postcontent">
      ${prop.text}
    </div>
    <hr style="display: ${(gotten_comm.length < 1) ? 'none' : 'block'};">
    <div class="commentAnnounce" style="display: ${(gotten_comm.length < 1) ? 'none' : 'block'};">
      Comments
    </div>
    <div class="postcomments">
      <div class="comment">
        ${comments}
      </div>
    </div>`;
  // Display image in the image viewing area if the image value is true
  if (prop.image_data) {
    const imageElement = document.createElement("img");
    imageElement.src = `data:image/jpeg;base64,${prop.image_data}`;
    imageViewingArea.appendChild(imageElement);
  }

  return gotten_comm;
};

const orgComments = async (postID) => {
  let comdiv = ``;
  let Response = await fetch("/api/comments", {
    method: "POST",
    body: `
        {
            "post_id" : ${postID}
        }
        `,
  });
  let gotten_comm = [];
  let commentArray = await Response.json();
  let i = 1;
  commentArray.comments.forEach((com) => {
    gotten_comm.push(com);
    comdiv += `<div class="scomment">
                        <div class="nameandlogo">
                            <div class="pfpImage">
                                <img src="../../static/assets/reddit.png" alt="reddit lol" class="pimg">
                            </div>
                            <div class="profileName">${com.user_name}</div>
                            <div class="commentDate">${com.creation_date}</div>
                        </div>
                        <div class="commenttext">
                            ${com.comment}
                        </div>
                        <div class="commentInfo">

                              <div id="likeCommBtn_${
                                com.comment_id
                              }" class="likeCommBtn ${
                                com.isLiked === 1 ? "liked" : ""
                              }" >
                                <img src="${
                                  com.isLiked === 1
                                    ? "../../static/assets/icons8-accept-30(1).png"
                                    : "../../static/assets/icons8-accept-30.png"
                                }" alt="Like">
                              </div>
                              <!-- Show like counts -->
                              <div id="Commlikes_${com.comment_id}">
                                    ${com.comment_likes}
                              </div>
                                    <div id="dislikeCommBtn_${
                                      com.comment_id
                                    }" class="dislikeCommBtn ${
                                      com.isLiked === -1 ? "disliked" : ""
                                    }" >
                                        <img src="${
                                          com.isLiked === -1
                                            ? "../../static/assets/icons8-dislike-30(1).png"
                                            : "../../static/assets/icons8-dislike-30.png"
                                        }" alt="Dislike">
                                    </div>
                              <!-- Show like counts -->
                              <div id="Commdislikes_${com.comment_id}">
                                    ${com.comment_dislikes}
                              </div>
                        </div>
                        </div>
                    `;
    i++;
  });

  return [comdiv, gotten_comm];
};
