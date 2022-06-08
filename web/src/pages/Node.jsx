import { Box, Container, Typography, IconButton, Grid, Button, Paper, Card, CircularProgress } from "@mui/material";
import { styled } from '@mui/material/styles';
import { useEffect, useState } from "react";
import { getDetail } from "../store/actions/nodeDetail";
import { useNavigate, useParams } from "react-router-dom";
import { useDispatch, useSelector } from "react-redux";
import Star from "../assets/star.svg"
import ClearTwoToneIcon from '@mui/icons-material/ClearTwoTone';
import CommandCard from "../components/ServiceCard/Main";
import DeleteIcon from '@mui/icons-material/Delete';
import Divider from '@mui/material/Divider';
import InfoTwoToneIcon from '@mui/icons-material/InfoTwoTone';
import AddRoundedIcon from '@mui/icons-material/AddRounded';
import ModeEditOutlineRoundedIcon from '@mui/icons-material/ModeEditOutlineRounded';
import BuildTwoToneIcon from '@mui/icons-material/BuildTwoTone';
import AddModal from "../components/NodeModal/AddModal";
import UpdateModal from "../components/NodeModal/UpdateModal";
import DeleteModal from "../components/NodeModal/DeleteModal";
import RefreshTwoToneIcon from '@mui/icons-material/RefreshTwoTone';
import AdjustIcon from '@mui/icons-material/Adjust';
import CachedIcon from '@mui/icons-material/Cached';
const ProjextBox = styled(Box)(({ theme }) => ({
  height: "100vh",
  display: "flex",
  flexDirection: "column"
}));

const CourseDetailPannelBox = styled(Card)(({ theme }) => ({
  display: "flex",
  flexDirection: "column",
  boxShadow: "-10px -11px 16px 1px rgba(252, 252, 252, 0.7), 9px 14px 24px -10px rgba(0, 0, 0, 0.25)",
  borderRadius: theme.spacing(2),
  backgroundColor: theme.palette.background.default,
  padding: theme.spacing(2),
}));


const StatusBox = styled(Box)(({ theme }) => ({
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

const ActionBox = styled(Box)(({ theme }) => ({
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


function Node({Id,Status,ServiceNum}) {
  const navigate = useNavigate();
  const params = useParams();
  const nodeId = params.Id;
  const dispatch = useDispatch();
  const courseDetail = useSelector(state => state.courseDetail.detail);
  const open = useSelector(state => state.auth.dialogOpen);
  const [loading, setLoading] = useState(true);
  const [refresh, setRefresh] = useState(true);
  const [openAddModal, setOpenAddModal] = useState(false);
  const [openUpdateModal, setOpenUpdateModal] = useState(false);
  const [openDeleteModal, setOpenDeleteModal] = useState(false);
  
  useEffect(() => {
    setLoading(true);
    dispatch(getDetail(nodeId)).catch(err => {
      console.log(err);
    }).then(() => {
      setLoading(false);
    });
  }, [params,refresh]);


  const ServiceCards = courseDetail.services.map((service)=>{
    return(
      <Grid item xs={12}>
        <CommandCard {...service} />
      </Grid>
    )
  })

  const serviceBlock = () =>{
    if(loading){
        return(
          <Box sx={{ display: "flex", justifyContent: "center" }}>
            <CircularProgress>Loading...</CircularProgress>
          </Box>
        ) 
    }
    if(ServiceCards.length == 0 && !loading){
      return(
      <Box sx={{ display: "flex", justifyContent: "center" }}>
        <Typography>查無資料</Typography>
      </Box>
      )
    }else{
      return(
        <Grid container spacing={3} sx={{ paddingBottom: "70px" }}>
         {ServiceCards}
        </Grid>
      )
    }
  }

  return (
    <>
    <Container maxWidth="md" sx={{ padding: "10px" }}>
      {/* <CourseDetailPannelBox> */}
        <ProjextBox>
            <Box sx={{ flexGrow: 1 }}>
              <Grid container spacing={3}>
                <Grid item xs={12}>
                  <Divider><InfoTwoToneIcon /></Divider>
                </Grid>
                <Grid item xs={12}>
                  <Box sx={{ display: "flex", alignItems: "center", justifyContent: "space-between" }}>
                    <Box sx={{ display: "flex", alignItems: "center", flexWrap: "nowrap" }}>
                      {courseDetail?
                      <>
                        {/* <img src={Star} style={{ height: "1rem" }} /> */}
                        <AdjustIcon color="action" sx={{color: "green"}}/>
                      </>
                        :
                        <ClearTwoToneIcon />
                      } 
                    </Box>
                    <Box>
                        <Typography variant="h6" sx={{ fontWeight: "bold",textAlign: "center", whiteSpace: "nowrap" }}>{courseDetail?"在線":"離線"}</Typography>
                    </Box>
                    <Box>
                      <Typography variant="h6" sx={{ fontWeight: "bold", fontSize: "1.5rem",textAlign: "center" }}>{nodeId}</Typography>
                    </Box>
                    <Box>
                    </Box>
                    <Box sx={{ display: "flex", alignItems: "center" }}>
                      {/* <InfoOutlinedIcon /> */}
                    </Box>
                  </Box>
                </Grid>
                <Grid item xs={12}>
                  <Divider><BuildTwoToneIcon /></Divider>
                </Grid>
                {/* <Grid item xs={12}>
                  <Typography sx={{ fontWeight: "bold", fontSize: "1.5rem", textAlign: "center" }} color="#757575">工作管理</Typography> 
                </Grid> */}
                <Grid item xs={12}>
                <ActionBox>
                  <Grid container spacing={3} direction="row" justifyContent="center" alignItems="center">
                    {/* <Grid item xs={4}>
                      <Box textAlign='center'>
                        <Button onClick={()=>setOpenAddModal(true)} variant="outlined" color="success" startIcon={<AddRoundedIcon />}>
                          新增
                        </Button>
                      </Box>
                    </Grid> */}
                    <Grid item xs={4}>
                      <Box textAlign='center'>
                        <Button onClick={()=>setRefresh(prev=>!prev)} variant="outlined" sx={{color:"blue"}} startIcon={<RefreshTwoToneIcon />}>
                          刷新此NODE的service
                        </Button>
                      </Box>
                    </Grid>
                    {/* <Grid item xs={4}>
                    <Box textAlign='center'>
                        <Button onClick={()=>setOpenUpdateModal(true)} variant="outlined" startIcon={<ModeEditOutlineRoundedIcon />}>
                          更新現有工作數量
                        </Button>
                      </Box>
                    </Grid>
                    <Grid item xs={4}>
                      <Box textAlign='center'>
                        <Button onClick={()=>setOpenDeleteModal(true)} variant="outlined" color="error" startIcon={<DeleteIcon />}>
                        刪除全部
                        </Button>
                      </Box>
                    </Grid> */}
                  </Grid>
                </ActionBox>
                </Grid>
                <Grid item xs={12}>
                  <Typography sx={{ fontWeight: "bold", fontSize: "1.5rem", textAlign: "center" }} color="#757575">運行在此node的SERVICE列表</Typography> 
                </Grid>
                <Grid item xs={12}>
                  <Grid container spacing={3}>
                    <Grid item xs={12}>
                      {serviceBlock()}
                    </Grid>
                  </Grid>
                </Grid>
              </Grid>
            </Box>
        </ProjextBox>
      {/* </CourseDetailPannelBox> */}
     </Container>
     <AddModal open={openAddModal} setOpen={setOpenAddModal} />
     <UpdateModal open={openUpdateModal} setOpen={setOpenUpdateModal} />
     <DeleteModal open={openDeleteModal} setOpen={setOpenDeleteModal} />
     </>
  )

}

export default Node;