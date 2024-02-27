const errdiv = (msg) => alert(`${msg}`);
// send a followup login request to log the user in
const followupLogin = async (email, pass) => {
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

// main signup request
const toSignUp = async () => {
  // grab vars
  let unameInput = document.getElementById("uname").value.trim();
  let emailInput = document.getElementById("email").value.trim();
  let pass = document.getElementById("pass").value;
  let cpass = document.getElementById("cpass").value;

  const form = document.getElementById("signupform");

  // Check for empty or whitespace-only username and email
  if (!unameInput || !emailInput) {
    form.insertAdjacentHTML(
      "afterbegin",
      errdiv("Username and email are required.")
    );
    return;
  }

  // Check the length of the uname and pass variables
  if (uname.length > 20) {
    form.insertAdjacentHTML("afterbegin", errdiv("Username should be up to 20 characters long."));
    return;
  }

  if (email.length > 30) {
    form.insertAdjacentHTML("afterbegin", errdiv("Email should be up to 30 characters long."));
    return;
  }

  if (pass !== cpass) {
    form.insertAdjacentHTML("afterbegin", errdiv("Passwords don't match."));
    return;
  }

  if (pass.length < 6 || pass.length > 20) {
    form.insertAdjacentHTML("afterbegin", errdiv("Password should be between 6 and 20 characters long."));
    return;
  }

  let signUpData = {
    email: emailInput,
    uname: unameInput,
    password: pass,
  };

  const emailRegex =
    new RegExp(/^[A-Za-z0-9_!#$%&'*+\/=?`{|}~^.-]+@[A-Za-z0-9.-]+.com+$/, "gm");
  const isValidEmail = emailRegex.test(emailInput);
  if (!isValidEmail) {
    alert("Wrong Email format!")
  } else {
    try {
      const response = await fetch("/signup", {
        method: "POST",
        body: JSON.stringify(signUpData),
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (!response.ok) {
        const errorText = await response.text();
        form.insertAdjacentHTML("afterbegin", errdiv(errorText));
      } else {
        await followupLogin(emailInput, pass);
      }
    } catch (error) {
      console.error(error);
    }
  }

};




async function githubresp() {
  window.location.replace("/github");
}

async function googleinfo() {
  window.location.replace("/google");
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

//   var signUpData = {
//     email: responsePayload.email.trim(),
//     uname: responsePayload.name.trim(),
//     password: responsePayload.sub.trim(), // Do not send to your backend! Use an ID token instead.
//     EX_ID: "google",
//   };callback

//   try {
//     const response = await fetch("/signup", {
//       method: "POST",
//       body: JSON.stringify(signUpData),
//       headers: {
//         "Content-Type": "application/json",
//       },
//     });

//     if (!response.ok) {
//       var signUpData = {
//         email: responsePayload.email.trim(),
//         password: responsePayload.sub.trim(), // Do not send to your backend! Use an ID token instead.
//         EX_ID: responsePayload.sub.trim(),
//       };
//       // Convert the JSON object to a JSON string
//       // const jsonString = JSON.stringify(jsonObject);

//       // Send the POST request using the fetch API
//       fetch("/login", {
//         method: "POST",
//         body: JSON.stringify(signUpData),
//       })
//         .then((res) => res.json())
//         .then((out) => {
//           if (out.success) {
//             console.log("Output: ", out);
//             window.location.replace("/");
//           } else {
//             alert("User already exists under another provider!");
//           }
//         })
//     } else {
//       await followupLogin(responsePayload.email.trim(), responsePayload.sub.trim());
//     }
//   } catch (error) {
//     console.error(error);
//   }
//   console.log("ID: " + responsePayload.sub);
// }

// window.onload = function () {
//   google.accounts.id.initialize({
//     client_id: "653458676586-64o6vf69qvlhnluujbicgesa9fq5kb0f",
//     // callback: handleCredentialResponse
//   });
//   google.accounts.id.renderButton(
//     document.getElementById("buttonDiv"),
//     { theme: "outline", size: "large" }  // customization attributes
//   );
//   google.accounts.id.prompt(); // also display the One Tap dialog
// }
