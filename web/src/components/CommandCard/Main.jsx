import { Box, Typography, ButtonBase, Link } from "@mui/material";
import { styled } from '@mui/material/styles';
import ChatOutlinedIcon from '@mui/icons-material/ChatOutlined';
import Star from "../../assets/star.svg"
import { useNavigate } from "react-router-dom";
import InfoOutlinedIcon from '@mui/icons-material/InfoOutlined';
import BoltTwoToneIcon from '@mui/icons-material/BoltTwoTone';
import CheckCircleTwoToneIcon from '@mui/icons-material/CheckCircleTwoTone';
import ClearTwoToneIcon from '@mui/icons-material/ClearTwoTone';
const CardBox = styled(Box)(({ theme }) => ({
  backgroundColor: theme.palette.background.default,
  borderRadius: theme.spacing(2),
  boxShadow: "-10px -11px 16px 1px rgba(252, 252, 252, 0.7), 9px 14px 24px -10px rgba(0, 0, 0, 0.25)",
  display: "flex",
  flexDirection: "column",
  justifyContent: "space-between",
  padding: theme.spacing(3),
  height: "100%",
  cursor: "pointer"
}));

const TotalRateAnnounce = styled(Typography)(({ theme }) => ({
  color: theme.palette.grey[400],
  fontSize: "0.9rem"
}));

function Card({commandStr}) {
  const navigate = useNavigate();

  const navigateToProjext = (id) => {
    navigate("/node/" + id)
  }

  return (
    <CardBox onClick={() => navigateToProjext(Id)}>
      <Box>
        <Box sx={{ display: "flex", alignItems: "center", justifyContent: "space-between" }}>
          <Box sx={{ display: "flex", alignItems: "center", flexWrap: "nowrap" }}>
          {Status===1?
          <>
            <img src={Star} style={{ height: "1rem" }} />
          </>
            // <CheckCircleTwoToneIcon />
            :
            <ClearTwoToneIcon />
          } 
          </Box>
          <Box>
            <Typography variant="h6" sx={{ fontWeight: "bold",textAlign: "center", whiteSpace: "nowrap" }}>在線</Typography>
          </Box>
          <Box>
            <Typography variant="h6" sx={{ fontWeight: "bold",textAlign: "center" }}>{commandStr}</Typography>
          {/* </Link> */}
          </Box>
          <Box>
          </Box>
        </Box>
      </Box>
    </CardBox>
  )
}

export default Card;