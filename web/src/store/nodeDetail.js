const courseDetail_state = {
  detail: {
    services:[]
  },
  loading: false
}

const detail = (state = {}, action) => {
  switch (action.type) {
    case "courseDetail.detail.set":
      return Object.assign({}, {...action.detail});
    case "courseDetail.detail.clear":
      return Object.assign({ services:[]}, {});
    default:
      return state;
  }
}

const loading = (state = false, action) => {
  switch (action.type) {
    case "courseDetail.loading":
      return true;
    case "courseDetail.done":
      return false;
    default:
      return state;
  }
}

const nodeDetail = (state = courseDetail_state, action) => {
  switch (action.type) {
    case "courseDetail.detail.set":
    case "courseDetail.detail.clear":
      return Object.assign({}, state, {
        detail: detail(state.detail, action)
      });
    case "courseDetail.loading":
    case "courseDetail.done":
      return Object.assign({}, state, {
        loading: loading(state.loading, action)
      });
    default:
      return state;
  }
}

export default nodeDetail;