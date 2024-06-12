import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { ThemeProvider, CssBaseline } from '@mui/material';
import { lightTheme, darkTheme } from './styles/Theme';

import sessionStore from './stores/SessionStore';

import Login from './components/login/Login';
import Home from './components/home/Home';
import Page from "./components/Page";
const App = () => {

  const [theme, setTheme] = useState(lightTheme);


  useEffect(() => {
    bind();

    return clear;
  });

  const bind = () => {
    sessionStore.addListener("theme_change", toggleTheme);
  }

  const clear = () => {
    sessionStore.removeListener("theme_change", toggleTheme);
  }

  const toggleTheme = () => {
    setTheme((prevTheme) => (prevTheme.palette.mode === 'light' ? darkTheme : lightTheme));
  };

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Router>
        <Routes>
          <Route path="/" element={<Page />}>
            <Route path="login" element={<Login />} />
            <Route path="home" element={<Home />}/>
            <Route path="license" element={<Home />}/>
          </Route>
        </Routes>
      </Router>
    </ThemeProvider>
  );
}

export default App;
