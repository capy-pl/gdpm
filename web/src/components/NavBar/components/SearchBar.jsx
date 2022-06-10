import { Box } from "@mui/material";
import { styled } from '@mui/material/styles';

const SearchBarBox = styled(Box)(({ theme }) => ({
    flexGrow: 1,
    margin: "0px 30px",
    height: "70%",
    borderRadius: "50px",
    boxShadow: "inset 2px 2px 8px rgba(0, 0, 0, 0.25)",
    display: "flex",
    alignItems: "center"
}));

const SearchBarInput = styled("input")(({ theme }) => ({
    height: "100%",
    width: "calc(100% - 20px)",
    margin: "0px 20px",
    border: "none",
    backgroundColor: "transparent",
    outlineWidth: 0,
    fontWeight: "bold",
    fontFamily: "Noto Sans TC"
}));

function SearchBar({search, setSearch}) {
    return (
        <SearchBarBox>
            <SearchBarInput value={search} onChange={(e) => setSearch(e.target.value)} placeholder="輸入課名、老師或系所" />
        </SearchBarBox>
    )
}

export default SearchBar;