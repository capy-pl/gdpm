import { Dialog } from "@mui/material";
import { styled } from "@mui/material/styles";
import { useSelector } from "react-redux";
import { isLoggedIn } from "../../store/selectors/auth"
import Login from "./components/Login";
import User from "./components/User"

const DialogBox = styled(Dialog)(({ theme }) => ({
  backgroundColor: theme.palette.background,
}));

function AuthDialog({ open, handleClose }) {
  const loggedIn = useSelector(state => isLoggedIn(state));

  return (
    <DialogBox onClose={handleClose} open={open}>
      {loggedIn ?
        <User />
        :
        <Login onClose={() => handleClose()} />
      }

    </DialogBox>
  );
}

export default AuthDialog;