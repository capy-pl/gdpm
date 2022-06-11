import { Box, Typography, IconButton, Grid, Link, Container, Button, Divider, CircularProgress } from "@mui/material";
import NavBar from "../components/NavBar/Main";
import { useEffect, useState } from "react";
import NodeCard from "../components/NodeCard/Main";
import { styled } from '@mui/material/styles';
import { useSearchParams, useNavigate } from "react-router-dom";
import { useDispatch, useSelector } from "react-redux";
import { coursesPagination, coursesLength } from "../store/selectors/node";
import { getCourse } from "../store/actions/node";
import AddRoundedIcon from '@mui/icons-material/AddRounded';
import AddModal from "../components/NodeModal/AddModal";
import BuildTwoToneIcon from '@mui/icons-material/BuildTwoTone';
import RefreshTwoToneIcon from '@mui/icons-material/RefreshTwoTone';
const HomeBox = styled(Box)(({ theme }) => ({
  height: "100vh",
  display: "flex",
  flexDirection: "column"
}));

const ActionBox = styled(Box)(({ theme }) => ({
  backgroundColor: theme.palette.background.default,
  borderRadius: theme.spacing(2),
  // boxShadow: "-10px -11px 16px 1px rgba(252, 252, 252, 0.7), 9px 14px 24px -10px rgba(0, 0, 0, 0.25)",
  display: "flex",
  flexDirection: "column",
  justifyContent: "space-between",
  padding: theme.spacing(2),
  height: "100%",
  cursor: "pointer"
}));

function Home() {
  const [searchParams, setSearchParams] = useSearchParams();
  const [search, setSearch] = useState(searchParams.get("search"));
  const [page, setPage] = useState(searchParams.get("page") ? searchParams.get("page") : 1);
  const [openAddModal, setOpenAddModal] = useState(false);
  const [refresh, setRefresh] = useState(true);
  const [loading, setLoading] = useState(true);
  const dispatch = useDispatch();
  const courses = useSelector(state => coursesPagination(state, page));
  ;

  // useEffect(() => {
  //   setSearch(searchParams.get("search"));
  //   setPage(searchParams.get("page") ? searchParams.get("page") : 1);
  // }, [searchParams]);

  useEffect(() => {
    const loadNode = async () => {
      setLoading(true);
      await dispatch(getCourse(search));
      setLoading(false);
    }
    loadNode();
  }, [refresh]);

  const coursesCards = courses.map((node, index) => {
    return (
      <Grid item xs={12} key={index}>
        <NodeCard
          {...node}
        />
      </Grid>
    )
  })

  const nodesBlock = () =>{
    if(loading){
        return(
          <Box sx={{ display: "flex", justifyContent: "center" }}>
            <CircularProgress>Loading...</CircularProgress>
          </Box>
        ) 
    }
    if(coursesCards.length == 0 && !loading){
      return(
      <Box sx={{ display: "flex", justifyContent: "center" }}>
        <Typography>查無資料</Typography>
      </Box>
      )
    }else{
      return(
        <Grid container spacing={3} sx={{ paddingBottom: "70px" }}>
         {coursesCards}
        </Grid>
      )
    }
  }

  return (
    <>
    <HomeBox>
      {/* <NavBar /> */}
      <Container maxWidth="sm" sx={{ padding: "15px" }}>
        <Box sx={{ flexGrow: 1 }}>
          <Grid container spacing={3}>
            <Grid item xs={12}>
              <Divider><BuildTwoToneIcon /></Divider>
            </Grid>
          <Grid item xs={12}>
                <ActionBox>
                  <Grid container spacing={3} direction="row" justifyContent="center" alignItems="center">
                    <Grid item xs={6}>
                      <Box textAlign='center'>
                        <Button onClick={()=>setOpenAddModal(true)} variant="outlined" color="success" startIcon={<AddRoundedIcon />}>
                          新增service
                        </Button>
                      </Box>
                    </Grid>
                    <Grid item xs={6}>
                      <Box textAlign='center'>
                        <Button onClick={()=>setRefresh(prev=>!prev)} variant="outlined" sx={{color:"blue"}} startIcon={<RefreshTwoToneIcon />}>
                          刷新nodes
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
            <Grid item xs={12}>
              <Typography sx={{ fontWeight: "bold", fontSize: "2rem", textAlign: "center" }} color="#757575">nodes</Typography> 
            </Grid>
              <Divider></Divider>
            </Grid>
            <Grid item xs={12}>
              {nodesBlock()}
            </Grid>
          </Grid>
        </Box>
      </Container>
    </HomeBox>
     <AddModal open={openAddModal} setOpen={setOpenAddModal} />
    </>
  )

}

export default Home;