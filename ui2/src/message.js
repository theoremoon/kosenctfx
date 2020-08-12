export function message(that, text) {
  that.$eventHub.$emit("message", {
    type: "message",
    text: text
  });
}

export function errorHandle(that, e) {
  let message = "";
  if (!e.response) {
    message = "Network Error";
  } else {
    message = e.response.data.message || "internal server error";
  }

  that.$eventHub.$emit("message", {
    type: "error",
    text: message
  });
}
