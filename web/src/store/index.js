import { applyMiddleware, combineReducers } from 'redux';
import { configureStore } from '@reduxjs/toolkit';
import thunkMiddleware from 'redux-thunk';
import { composeWithDevTools } from '@redux-devtools/extension';
import auth from './auth';
import error from './error';
import course from "./node";
import courseDetail from "./nodeDetail";

const App = combineReducers({
  auth,
  error,
  course,
  courseDetail,
})

const store = configureStore({
  reducer: App,
  // Note that this will replace all default middleware
  middleware: [thunkMiddleware],
})

export default store