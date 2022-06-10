import { Box, Typography, ButtonBase, Link, Grid, Button, IconButton, Divider } from "@mui/material";
import { styled } from '@mui/material/styles';
import { useEffect, useState } from "react";
import ChatOutlinedIcon from '@mui/icons-material/ChatOutlined';
import Star from "../../assets/star.svg"
import { useNavigate } from "react-router-dom";
import InfoOutlinedIcon from '@mui/icons-material/InfoOutlined';
import BoltTwoToneIcon from '@mui/icons-material/BoltTwoTone';
import CheckCircleTwoToneIcon from '@mui/icons-material/CheckCircleTwoTone';
import ClearTwoToneIcon from '@mui/icons-material/ClearTwoTone';
import UpdateModal from "../NodeModal/UpdateModal";
import DeleteModal from "../NodeModal/DeleteModal";
import ModeEditOutlineRoundedIcon from '@mui/icons-material/ModeEditOutlineRounded';
import DeleteIcon from '@mui/icons-material/Delete';
import CachedIcon from '@mui/icons-material/Cached';
import CircularProgress from '@mui/material/CircularProgress';
const CardBox = styled(Box)(({ theme }) => ({
  backgroundColor: theme.palette.background.default,
  borderRadius: theme.spacing(2),
  boxShadow: "-10px -11px 16px 1px rgba(252, 252, 252, 0.7), 9px 14px 24px -10px rgba(0, 0, 0, 0.25)",
  display: "flex",
  flexDirection: "column",
  justifyContent: "space-between",
  padding: theme.spacing(2),
  height: "100%",
  cursor: "pointer"
}));

const TotalRateAnnounce = styled(Typography)(({ theme }) => ({
  color: theme.palette.grey[400],
  fontSize: "0.9rem"
}));

function CommandCard(prop) {
  const navigate = useNavigate();
  const [openUpdateModal, setOpenUpdateModal] = useState(false);
  const [openDeleteModal, setOpenDeleteModal] = useState(false);
  const navigateToProjext = (id) => {
    navigate("/node/" + id)
  }

  return (
    <>
    <CardBox onClick={() => {}}>
        <Box sx={{ display: "flex", alignItems: "center", justifyContent: "space-between" }}>
          <Grid container columns={24} spacing={1}>
          <Grid item xs={12} md={2} lg={2}>
            <Box sx={{ display: "flex", alignItems: "center", flexWrap: "nowrap" }}>
                {/* <img src={Star} style={{  }} /> */}
                <CircularProgress disableShrink />
                {/* <CachedIcon sx={{ fontWeight: "bold",textAlign: "center", lineHeight:"2.5",fontSize: "2.2rem", color:"green.50" }}/> */}
            </Box>
          </Grid>
          <Grid item xs={12} md={3} lg={3}>
            <Box sx={{ display: "flex", alignItems: "center", flexWrap: "nowrap" }}>
            <Typography variant="h6" sx={{ fontWeight: "bold",textAlign: "center", lineHeight:"2.5",fontSize: "1rem" }}>數量：{prop.Number}</Typography>
            </Box>
          </Grid>
          <Grid item xs={24} md={15} lg={15}>
            {/* <Typography variant="h6" sx={{ fontWeight: "bold",textAlign: "center", lineHeight:"2.5",fontSize: "1rem" }}>數量：{prop.Number}</Typography> */}
            {/* <Divider orientation="vertical" sx={{content: '""'}}>&nbsp;</Divider> */}
            <Typography variant="h6" sx={{ fontWeight: "bold",textAlign: "center", verticalAlign: "middle", fontSize: "1rem", lineHeight:"2.5", fontFamily: "Consolas", whiteSpace: { xs: 'pre-wrap', sm: 'nowrap', md: 'nowrap' } }}>{prop.Command}</Typography>
          {/* </Link> */}
          </Grid>
          <Grid item xs={12} md={2} lg={2}>
          <Box textAlign='center'>
              <IconButton onClick={()=>setOpenUpdateModal(true)} variant="outlined" startIcon={<ModeEditOutlineRoundedIcon />}>
                {/* 更新現有工作數量 */}
                <ModeEditOutlineRoundedIcon />
              </IconButton>
            </Box>
          </Grid>
          <Grid item xs={6} md={2} lg={2}>
            <Box textAlign='center'>
              <IconButton onClick={()=>setOpenDeleteModal(true)} variant="outlined" color="error" startIcon={<DeleteIcon />}>
              {/* 刪除全部 */}
              <DeleteIcon />
              </IconButton>
            </Box>
          </Grid>
          <Box>
          </Box>
          </Grid>
        </Box>
    </CardBox>
    <UpdateModal service={prop} open={openUpdateModal} setOpen={setOpenUpdateModal} />
    <DeleteModal service={prop} open={openDeleteModal} setOpen={setOpenDeleteModal} />
    </>
  )
}

export default CommandCard;