// function onSignIn(googleUser) {
//   var profile = googleUser.getBasicProfile();
//
//   // Prepare the data to send to the backend
//   var data = {
//     id: profile.getId(), // Do not send to your backend! Use an ID token instead.
//     name: profile.getName(),
//     imageUrl: profile.getImageUrl(),
//     email: profile.getEmail(),
//   };
//
//   const url = "/login";
//
//   // Send the POST request using the fetch API
//   fetch(url, {
//     method: "POST",
//     headers: {
//       "Content-Type": "application/json",
//     },
//     body: JSON.stringify(data),
//   })
//     .then((response) => {
//       if (response.ok) {
//         return response.json();
//       } else {
//         throw new Error("Request failed");
//       }
//     })
//     .then((data) => {
//       console.log(data);
//     })
//     .catch((error) => {
//       console.error(error);
//     });
// }

function sendLoginRequest() {
  event.preventDefault();
  const url = "/login";

  let formData = new FormData(document.querySelector("form"));

  // Convert form data to JSON object
  const jsonObject = {};
  formData.forEach((value, key) => {
    jsonObject[key] = value;
  });

  // Convert the JSON object to a JSON string
  const jsonString = JSON.stringify(jsonObject);

  // Send the POST request using the fetch API
  fetch(url, {
    method: "POST",
    body: jsonString,
  })
    .then((res) => res.json())
    .then((out) => {
      if (out.success) {
        console.log("Output: ", out);

        window.location.replace("/");
      } else {
        alert("login error!");
      }
    })
    .catch((err) => console.error(err));
}
