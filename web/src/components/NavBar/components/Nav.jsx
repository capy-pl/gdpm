import { Box, Typography, IconButton } from "@mui/material";
import { styled } from '@mui/material/styles';
import { useEffect, useState } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import SearchIcon from '@mui/icons-material/Search';
import SemesterSelect from "./SemesterSelect"
import SearchBar from "./SearchBar";

const NavBox = styled(Box)(({ theme }) => ({
  display: "flex",
  alignItems: "center",
  boxShadow: "-10px -11px 16px 1px rgba(252, 252, 252, 0.7), 9px 14px 24px -10px rgba(0, 0, 0, 0.25)",
  borderRadius: "35px",
  margin: "0px 30px",
  padding: "0px 30px",
  height: "100%"
}));


function Nav({ hideTitle }) {
  const [search, setSearch] = useState("");
  const [searchParams, setSearchParams] = useSearchParams();
  const navigate = useNavigate();

  useEffect(() => {
    if (searchParams.get("search")) setSearch(searchParams.get("search"));
  }, [searchParams]);

  const handleSubmit = (e) => {
    e.preventDefault();
    navigate({
      pathname: '/search',
      search: '?search=' + search,
    });
  }

  return (
    <form onSubmit={handleSubmit} style={{ height: "100%" }}>
      {
        hideTitle ?
          <Box sx={{ display: "flex", alignItems: "center", height: "50px" }}>
            <SearchBar search={search} setSearch={setSearch} />
            <IconButton color="primary" type="submit" size="large">
              <SearchIcon fontSize="large" />
            </IconButton>
          </Box>
          :
          <NavBox>
            <Typography variant="h5" sx={{ fontWeight: "bold" }} color="primary">政大課程評價網</Typography>
            <SearchBar search={search} setSearch={setSearch} />
            <IconButton color="primary" type="submit">
              <SearchIcon />
            </IconButton>
          </NavBox>
      }

    </form>
  )
}

export default Nav;