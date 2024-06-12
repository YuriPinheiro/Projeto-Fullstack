import { useTheme } from '@mui/material/styles';

const MenuBarStyles = () => {
  const theme = useTheme();
    console.log(theme)
  return {
    container: {
     background: theme.palette.colors.onPrimary,
    }
  };
};

export default MenuBarStyles;
