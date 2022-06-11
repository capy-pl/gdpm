import { createSelector } from "reselect"

const PAGE_COUNT = 18

export const coursesPagination = createSelector(
  [
    state => state.course?.result?.nodes,
    (state, page) => page
  ],
  (result, page) => {
    if(!result) return [];
    // return result.filter((course, index) => {
    //   return index >= (page - 1) * PAGE_COUNT && index < page * PAGE_COUNT;
    // })

    return result;
  }
)

export const coursesLength = createSelector(
  (state) => state.course?.result?.courses,
  (result) => {
    if(!result) return 0
    return Math.floor(result.length / PAGE_COUNT) + 1
  }
)