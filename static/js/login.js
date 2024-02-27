const followupLogin2 = async (email, pass) => {
  const jsonObject = {
    Email: email,
    Password: pass,
  };

  const jsonString = JSON.stringify(jsonObject);

  fetch("/login", {
    method: "POST",
    body: jsonString,
  })
    .then((response) => {
      if (response.ok) {
        window.location.replace("/");
      } else {
        throw new Error("Request failed");
      }
    })
    .then((data) => {
      console.log(data);
    })
    .catch((error) => {
      console.error(error);
    });
};

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


// function decodeJwtResponse(jwtToken) {
//   // Split the token into its three parts: header, payload, and signature
//   const parts = jwtToken.split('.');

//   // Check if token has payload (it should)
//   if (parts.length === 3) {
//     // Decode the payload (middle part)
//     const decodedPayload = atob(parts[1]);

//     // Parse the JSON string to access the claims
//     const payloadObject = JSON.parse(decodedPayload);

//     // Return the payload object
//     return payloadObject;
//   } else {
//     // Return null or throw an error if token format is invalid
//     return null;
//   }
// }

// async function handleCredentialResponse(response) {
//   const responsePayload = decodeJwtResponse(response.credential);
//   try {
//     var signInData = {
//       email: responsePayload.email.trim(),
//       password: responsePayload.sub.trim(), // Do not send to your backend! Use an ID token instead.
//       EX_ID: responsePayload.sub.trim(),
//     };
//     const loginResponse = await fetch("/login", {
//       method: "POST",
//       body: JSON.stringify(signInData),
//     });
//     const loginResult = await loginResponse.json();
//     if (loginResult.success) {
//       console.log("Output: ", loginResult);
//       window.location.replace("/");
//     } else {
//       var signUpData = {
//         email: responsePayload.email.trim(),
//         uname: responsePayload.name.trim(),
//         password: responsePayload.sub.trim(), // Do not send to your backend! Use an ID token instead.
//         EX_ID: "google",
//       };
//       const signUpResponse = await fetch("/signup", {
//         method: "POST",
//         body: JSON.stringify(signUpData),
//         headers: {
//           "Content-Type": "application/json",
//         },
//       });
//       const sigupResult = await signUpResponse.json();

//       if (!sigupResult.success) {
//         alert("This user was registered under another provider!");
//       } else {
//         await followupLogin2(responsePayload.email.trim(), responsePayload.sub.trim());
//       }
//     }
//   } catch (error) {
//     console.error(error);
//   }
// }

// window.onload = function () {
//   google.accounts.id.initialize({
//     client_id: "653458676586-64o6vf69qvlhnluujbicgesa9fq5kb0f",
//     callback: handleCredentialResponse
//   });
//   google.accounts.id.renderButton(
//     document.getElementById("buttonDiv"),
//     { theme: "outline", size: "large" }  // customization attributes
//   );
//   google.accounts.id.prompt(); // also display the One Tap dialog
// }

async function githubresp2() {
  window.location.replace("/github");
}
async function googleinfo() {
  window.location.replace("/google");
}
