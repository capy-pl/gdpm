import { Box, Typography, ButtonBase, Link } from "@mui/material";
import { styled } from '@mui/material/styles';
import ChatOutlinedIcon from '@mui/icons-material/ChatOutlined';
import Star from "../../assets/star.svg"
import { useNavigate } from "react-router-dom";
import InfoOutlinedIcon from '@mui/icons-material/InfoOutlined';
import BoltTwoToneIcon from '@mui/icons-material/BoltTwoTone';
import CheckCircleTwoToneIcon from '@mui/icons-material/CheckCircleTwoTone';
import ClearTwoToneIcon from '@mui/icons-material/ClearTwoTone';
import AdjustIcon from '@mui/icons-material/Adjust';

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

function Card({ Id, ServuceNum, Status, Times}) {
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
            {/* <img src={Star} style={{ height: "1rem" }} /> */}
            <AdjustIcon color="action" sx={{color: "green"}}/>
          </>
            // <CheckCircleTwoToneIcon />
            :
            <ClearTwoToneIcon />
          } 
          </Box>
          <Box>
          {/* <Link href={url} target="_blank" variant="body2"> */}
            <Typography variant="h6" sx={{ fontWeight: "bold",textAlign: "center", whiteSpace: "nowrap" }}>在線</Typography>
          {/* </Link> */}
          </Box>
          <Box>
          {/* <Link href={url} target="_blank" variant="body2"> */}
            <Typography variant="h6" sx={{ fontWeight: "bold",textAlign: "center" }}>{Id}</Typography>
          {/* </Link> */}
          </Box>
          <Box>
            {/* <TotalRateAnnounce>{enName}</TotalRateAnnounce> */}
          </Box>
          <Box sx={{ display: "flex", alignItems: "center" }}>
          <Typography variant="h6" sx={{ fontWeight: "bold",textAlign: "center", fontSize: "0.9rem" }}>{`${Times}`}</Typography>
            {/* <InfoOutlinedIcon /> */}
          </Box>
        </Box>
      </Box>
    </CardBox>
  )
}

export default Card;