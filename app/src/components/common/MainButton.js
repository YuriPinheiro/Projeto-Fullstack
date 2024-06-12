import React from 'react';
import { styled } from '@mui/material/styles';
import Button from '@mui/material/Button';

const MyButton = styled(Button)(({ theme }) => ({
  backgroundColor: theme.palette.primary.main,
  color: theme.palette.common.white,
  '&:hover': {
    backgroundColor: theme.palette.primary.dark,
  },
}));

const StyledButton = (props) => {
  return <MyButton size='sm' onClick={props.onClick}>{props.text}</MyButton>;
}

export default StyledButton;
