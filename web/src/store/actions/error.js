const error_occur = () => {
  return {
    type: "error.errorStatus.error",
  }
}

const clear_error = () => {
  return {
    type: "error.errorStatus.clear",
  }
}

const set_errorMessage = (errorMessage) => {
  return {
    type: "error.errorMessage.set",
    errorMessage: errorMessage
  }
}

const clear_errorMessage = () => {
  return {
    type: "error.errorMessage.clear",
  }
}

export const setError = (errorMessage) => {
  return function (dispatch) {
    dispatch(error_occur());
    dispatch(set_errorMessage(errorMessage));
  }
}

export const clearError = () => {
  return function (dispatch) {
    dispatch(clear_error());
    dispatch(clear_errorMessage());
  }
}
