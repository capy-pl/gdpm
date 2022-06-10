import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from "@mui/material/CssBaseline";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import theme from "./theme/theme";
import { useSelector, useDispatch } from "react-redux";
import Home from "./pages/Home";
import Node from "./pages/Node";
import NavBar from "./components/NavBar/Main";
import AuthModal from "./components/AuthDialog/Main";

function App() {
  const dispatch = useDispatch();
  const open = useSelector(state => state.auth.dialogOpen);
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <BrowserRouter>
      <NavBar />
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/node" >
            <Route path=":Id" element={<Node />} />
          </Route>
        </Routes>
        <AuthModal open={open} handleClose={() => dispatch({type: "auth.dialog.close"})} />
      </BrowserRouter>
    </ThemeProvider>
  )
}

export default App
