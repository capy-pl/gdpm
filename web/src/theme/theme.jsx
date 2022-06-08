import { LightColor } from "./color";
import { createTheme } from '@mui/material/styles';

const palette= {
  primary: {
    main: LightColor.primary,
    contrastText: LightColor.light,
  },
  secondary: {
    main: LightColor.secondary,
    contrastText: LightColor.light,
  },
  tagRequired: {
    main: LightColor.tagRequired,
    contrastText: LightColor.light,
  },
  tagPartial: {
    main: LightColor.tagPartial,
    contrastText: LightColor.light,
  },
  tagSelect: {
    main: LightColor.tagSelect,
    contrastText: LightColor.light,
  },
  background: {
    default: LightColor.background
  },
  grey: {
    50:"#F0F0F0",
    400: "#B9B9B9"
  },
  yellow:{
    50:"yellow",
    400:""
  }
};

export default createTheme({
  palette: palette,
  breakpoints: {
    values: {
      xs: 0,
      sm: 500,
      md: 800,
      lg: 1100,
      xl: 1500,
    },
  },
  typography: {
    fontFamily: [
      'Noto Sans TC',
      'sans-serif',
    ].join(','),
    fontWeightRegular: 500
  },
});