import { Box, Typography, AppBar } from "@mui/material";
import { styled } from '@mui/material/styles';
import { useLocation } from 'react-router-dom';
import Icon from "./components/Icon";
import Account from "./components/Account";
import Nav from "./components/Nav";

const NavbarBox = styled(AppBar)(({ theme }) => ({
  padding: "20px",
  backgroundColor: theme.palette.background.default
}));

function Navbar() {
  const location = useLocation();
  return (
    <NavbarBox position="sticky" color="transparent" elevation={0}>
      <Box sx={{ display: "flex" }}>
        <Icon />
        <Box sx={{ flexGrow: 1 }}>
          {/* <Box sx={{ display: { xs: 'none', sm: 'none', md: 'block' }, height: "100%" }}> */}
          <Box sx={{ display: { }, height: "100%" }}>
          <Typography variant="h5" sx={{ fontWeight: "bold",textAlign: "center" }} color="primary">ETCD節點管理</Typography>
          </Box>
        </Box>
        <Account />
      </Box>
    </NavbarBox>
  )
}

export default Navbar;