import ajax from "../../api"
import axios from "axios"

const set_course_detail = (detail) => {
  return {
    type: "courseDetail.detail.set",
    detail: detail
  }
}

const clear_course_detail = () => {
  return {
    type: "courseDetail.detail.clear"
  }
}

const set_loading = () => {
  return {
    type: "courseDetail.loading"
  }
}

const set_done = () => {
  return {
    type: "courseDetail.done"
  }
}

export const getDetail = (id) => {
  return function (dispatch) {
    dispatch(clear_course_detail());
    dispatch(set_loading());
    return ajax(`/node/${id}/`, "get", { params: { } }).then(res => {
    // console.log('res :', res);
      let services = res.data.Ids.map((value,index)=>{
        return (
          {
            Id:value,
            Command:res.data.Command[index],
            Number:res.data.Number[index]
          }
          )
        })
      console.log('services :', services);
      dispatch(set_course_detail({services,id}));
      dispatch(set_done());
      return res.data;
    })
  }
}

export const createWork = (command,numOfWork) => {
  return function (dispatch, getState) {
    const nodeId = getState().courseDetail.detail.id;
    const url = `/service/`;
    const params = new URLSearchParams()
    params.append('Command', command)
    params.append('InstanceNum', numOfWork)
    const config = {
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded'
      }
    }
    return axios.post(url, params, config).then(res => {
      console.log('res :', res);
      return ajax(`/node/${nodeId}/`, "get", { params: { } })
    }).then(res => {
      dispatch(set_course_detail(res.data));
      dispatch(set_done());
      return res.data;
    })
  }
}
export const updateWork = (serviceId,command,numOfWork) => {
  return function (dispatch, getState) {
    const userId = getState().auth?.profile?.id
    return ajax(`/${serviceId}/`, "post", { data: { command: command, numOfWork: numOfWork } }).then(res => {
      return ajax(`/node/${nodeId}/`, "get", { params: { } })
    }).then(res => {
      dispatch(set_course_detail(res.data));
      dispatch(set_done());
      return res.data;
    })
  }
}
export const deleteWork = (serviceId,command,numOfWork) => {
  return function (dispatch, getState) {
    const userId = getState().auth?.profile?.id
    return ajax(`/${serviceId}/`, "delete", { data: { command: command, numOfWork: numOfWork } }).then(res => {
      return ajax(`/node/${nodeId}/`, "get", { params: { } })
    }).then(res => {
      dispatch(set_course_detail(res.data));
      dispatch(set_done());
      return res.data;
    })
  }
}
