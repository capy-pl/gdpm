import { createSelector } from "reselect"

export const isLoggedIn = createSelector(
  (state) => state.auth.profile,
  (profile) => profile != null
)

export const googleProfile = createSelector(
  (state) => state.auth.googleProfile,
  (profile) => {
    if(profile) return profile;
    return {
      email: '',
      name: '',
      picture: '',
      locale: 'zh-TW',
    }
  }
)