import axios from "axios"

// axios.defaults.baseURL = "";
axios.defaults.baseURL = "https://affb5559ad8b.jp.ngrok.io"
// axios.defaults.baseURL = "http://localhost:8989"
// axios.defaults.headers.post["Content-Type"] = "application/json";
axios.defaults.headers.post["Content-Type"] = "X-www-form-urlencoded";

/*
  Function Usage Sample:

  ajax("/api/user/login", "post", {
    data: {}
    params:{}
  }).then(res => {  
    ...
  })
  
*/

const ajax = (url, method, options) => {
  if (options !== undefined) {
    var { params = {}, data = {} } = options
  } else {
    params = data = null;
  }
  return new Promise((resolve, reject) => {
    axios({ url, method, params, data }).then(res => {
      resolve(res);
    }, res => {
      reject(res);
    })
  })
}

export default ajax