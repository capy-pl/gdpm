const error_state = {
  errorStatus: false,
  errorMessage: "testing some error message now"
}

const errorStatus = (state = false, action) => {
  switch (action.type) {
    case "error.errorStatus.error":
      return true;
    case "error.errorStatus.clear":
      return false;
    case "error.errorStatus.get":
    default:
      return state;
  }
}

const errorMessage = (state = "", action) => {
  switch (action.type) {
    case "error.errorMessage.set":
      return action.errorMessage;
    case "error.errorMessage.clear":
      return "";
    case "error.errorMessage.get":
    default:
      return state;
  }
}

const error = (state = error_state, action) => {
  switch (action.type) {
    case "error.errorStatus.error":
    case "error.errorStatus.clear":
    case "error.errorStatus.get":
      return Object.assign({}, state, {
        errorStatus: errorStatus(state.errorStatus, action)
      });
    case "error.errorMessage.set":
    case "error.errorMessage.clear":
    case "error.errorMessage.get":
      return Object.assign({}, state, {
        errorMessage: errorMessage(state.errorMessage, action)
      });
    default:
      return state;
  }
}

export default error;