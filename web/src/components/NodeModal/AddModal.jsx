import * as React from 'react';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import TextField from '@mui/material/TextField';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import { createWork } from "../../store/actions/nodeDetail";
import { useDispatch, useSelector } from "react-redux";
import CircularProgress from '@mui/material/CircularProgress';
export default function FormDialog({open,setOpen}) {
//   const [open, setOpen] = React.useState(false);
  const [loading, setLoading] = React.useState(false);
  const [command, setCommand] = React.useState("python /Users/william/go/src/github.com/gdpm/test.py");
  const [numOfWork, setNumOfWork] = React.useState(1);
  const dispatch = useDispatch();
  const courseDetail = useSelector(state => state.courseDetail.detail);

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const handleSubmit = (e) => {
      e.preventDefault();
      setLoading(true);
    dispatch(createWork(command,numOfWork)).then(() => {
        setOpen(false);
        setLoading(false);
        setOpen(false);
    });
  }

  if(loading) {
    return <CircularProgress>Loading...</CircularProgress>
  }

  return (
    <Box sx={{width:"100%"}}>
        <Dialog open={open} onClose={handleClose}>
            <form onSubmit={handleSubmit}>
                <DialogTitle>新增工作</DialogTitle>
                <DialogContent>
                <DialogContentText>
                    請輸入指令跟工作數
                </DialogContentText>
                    <TextField
                        autoFocus
                        margin="dense"
                        id="name"
                        label="Command"
                        fullWidth
                        required
                        variant="standard"
                        value={command}
                        onChange={(e) => setCommand(e.target.value)}
                        />
                    <TextField
                        id="standard-number"
                        required
                        label="Number"
                        type="number"
                        variant="standard"
                        value={numOfWork}
                        onChange={(e) => setNumOfWork(e.target.value)}
                    />
                </DialogContent>
                <DialogActions>
                <Button onClick={handleClose}>取消</Button>
                <Button type="submit">新增</Button>
                </DialogActions>
            </form>
        </Dialog>
    </Box>
  );
}
