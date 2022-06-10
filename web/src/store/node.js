const course_state = {
  result: {nodes:[]},
  loading: false,
}

const result = (state = [], action) => {
  switch (action.type) {
    case "course.result.set":
      return Object.assign({}, action.result);
    case "course.result.clear":
      return Object.assign({}, {});
    default:
      return state;
  }
}

const loading = (state = false, action) => {
  switch (action.type) {
    case "course.loading":
      return true;
    case "course.done":
      return false;
    default:
      return state;
  }
}

const node = (state = course_state, action) => {
  switch (action.type) {
    case "course.result.set":
    case "course.result.clear":
      return Object.assign({}, state, {
        result: result(state.result, action)
      });
    case "course.loading":
    case "course.done":
      return Object.assign({}, state, {
        loading: loading(state.loading, action)
      });
    default:
      return state;
  }
}

export default node;