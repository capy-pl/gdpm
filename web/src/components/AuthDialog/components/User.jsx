import { DialogTitle, DialogContent, DialogContentText, DialogActions, Button, Box, Avatar, Typography } from "@mui/material";
import GoogleIcon from '@mui/icons-material/Google';
import { GoogleLogout } from 'react-google-login';
import { styled } from "@mui/material/styles";
import { useDispatch, useSelector } from "react-redux";
import { googleProfile } from "../../../store/selectors/auth"
import { logout } from "../../../store/actions/auth";

function User() {
  const dispatch = useDispatch();
  const user = useSelector(state => googleProfile(state));
  const avatar = user.picture;
  const name = user.name;
  const email = user.email

  const responseGoogle = (e) => {
    dispatch(logout());
  };

  return (
    <Box>
      <DialogTitle sx={{ display: "flex", alignItems: "center" }}>
        <Avatar alt="Remy Sharp" src={avatar} />
        <Box sx={{ display: "flex", flexDirection: "column" }}>
          <Typography sx={{ margin: "0px 10px" }}>{name}</Typography>
          <Typography sx={{ margin: "0px 10px" }}>{email}</Typography>
        </Box>
      </DialogTitle>
      <DialogActions>
        <GoogleLogout
          clientId="22551978498-3d7pfatc0km7mpm8t6glfuu4ev2jld3a.apps.googleusercontent.com"
          buttonText="Logout"
          render={renderProps => (
            <Button variant="contained" onClick={renderProps.onClick} disabled={renderProps.disabled} startIcon={<GoogleIcon />}>
              Logout
            </Button>
          )}
          onLogoutSuccess={responseGoogle}
        >
        </GoogleLogout>
      </DialogActions>
    </Box>
  )
}

export default User;