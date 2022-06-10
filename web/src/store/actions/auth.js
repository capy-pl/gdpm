import jwt_decode from "jwt-decode";
import ajax from "../../api";

const set_profile = (profile) => {
  return {
    type: "auth.profile.set",
    profile: profile
  }
}

const clear_profile = () => {
  return {
    type: "auth.profile.clear",
  }
}

const set_google_profile = (profile) => {
  return {
    type: "auth.googleProfile.set",
    googleProfile: profile
  }
}

const clear_google_profile = () => {
  return {
    type: "auth.googleProfile.clear",
  }
}

export const login = () => {
  return function (dispatch) {
    const authToken = localStorage.getItem("authToken")  
    return ajax("/user/googlesignin", "post", {
      data: {
        token: authToken
      }
    }).then(res => {
      dispatch(set_profile(res.data));
      dispatch(set_google_profile(jwt_decode(authToken)));

      return res.data;
    })
  }
}

export const logout = () => {
  return function (dispatch) {
    dispatch(clear_profile());
    dispatch(clear_google_profile());
    localStorage.removeItem("authToken");
  }
}

export const setProfile = (profile) => {
  return function (dispatch) {
    dispatch(set_profile(profile));
  }
}
