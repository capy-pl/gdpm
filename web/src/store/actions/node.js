import ajax from "../../api"

const set_course_result = (result) => {
  return {
    type: "course.result.set",
    result: result
  }
}

const clear_course_result = () => {
  return {
    type: "course.result.clear"
  }
}

export const getCourse = (searchParam) => {
  return function (dispatch) {
    dispatch({type: "course.result.clear"})
    dispatch({ type: "course.loading" })
    return ajax("/node/", "get", { }).then(res => {
      dispatch({ type: "course.done" });
      let nodes = res.data.Ids.map((value,index)=>{
        return (
          {
            Id:value,
            ServiceNum:res.data.ServiceNum[index],
            Status:res.data.Status[index],
            Times:res.data.Times[index]
          }
        )
      })
      dispatch({
        type: "course.result.set",
        result: {nodes}
      });
      return nodes;
    })
  }
}

export const getCourseDetail = () => { }