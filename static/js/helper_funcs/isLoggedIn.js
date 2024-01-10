export const isloggedIn = async () => {
  let res = await fetch("/api/islogged");
  let ok = await res.text();
  /* console.log(ok); */
  let isSignedIn;
  if (ok === "1") {
    isSignedIn = "true";
  } else {
    isSignedIn = "false";
  }
  localStorage.setItem("isloggedIn", isSignedIn);
  return isSignedIn;
};
