import React, { useState } from 'react';
import MainButton from "../common/MainButton";
import TextField from '@mui/material/TextField';
import Styles from "../../styles/LoginStyle";

import { Grid } from '@mui/material';
const Login = () => {

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const classes = Styles();

  const login = () => {
    console.log(email)
    console.log(password)
  }

  const onChangeInput = (event) => {
    console.log(event)
    if (event.target.name === 'email') {
      setEmail(event.target.value)
    }
    if (event.target.name === 'password') {
      setPassword(event.target.value)
    }
  }

  return (
    <Grid container sx={classes.container} spacing={2} direction="column" justifyContent="center" alignItems="center">
      <Grid item xs>
        <TextField id="standard-basic" label="Email" name='email' variant="standard" onChange={onChangeInput} />
      </Grid>
      <Grid item xs>
        <TextField id="standard-basic" label="Senha" name='password' variant="standard" type="password" onChange={onChangeInput} />
      </Grid>
      <Grid item xs={12}>
        <MainButton  onClick={login} text={"Login"} />
      </Grid>
    </Grid>
  );
}

export default Login;
