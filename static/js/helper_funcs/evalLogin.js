import { isloggedIn } from "./isLoggedIn.js";

export async function evalLogin(fn) {
  let islogged = await isloggedIn();
  if (islogged === "true") {
    await fn();
  } else {
    window.location.replace("/login");
  }
}
